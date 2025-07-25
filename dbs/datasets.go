package dbs

// DBS dataset module
// Copyright (c) 2023 - Valentin Kuznetsov <vkuznet@gmail.com>
//
import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	lexicon "github.com/CHESSComputing/golib/lexicon"
	"github.com/CHESSComputing/golib/utils"
)

// Datasets represents Datasets DBS DB table
type Datasets struct {
	DATASET_ID    int64  `json:"dataset_id"`
	DID           string `json:"did" validate:"required"`
	SITE_ID       int64  `json:"site_id" validate:"required",number`
	PROCESSING_ID int64  `json:"processing_id" validate:"required",number`
	OSINFO_ID     int64  `json:"osinfo_id" validate:"required",number`
	CONFIG_ID     int64  `json:"config_id" validate:"number"`
	CREATE_AT     int64  `json:"create_at" validate:"required,number"`
	CREATE_BY     string `json:"create_by" validate:"required"`
	MODIFY_AT     int64  `json:"modify_at" validate:"required,number"`
	MODIFY_BY     string `json:"modify_by" validate:"required"`
}

// Datasets API
//
//gocyclo:ignore
func (a *API) GetDataset() error {
	if Verbose > 1 {
		log.Printf("datasets params %+v", a.Params)
	}
	var args []interface{}
	var conds []string
	tmpl := make(map[string]any)
	tmpl["Owner"] = DBOWNER

	allowed := []string{"did", "file", "script", "environment", "package", "site", "bucket", "osname", "processing", "config"}
	for k, _ := range a.Params {
		if !utils.InList(k, allowed) {
			msg := fmt.Sprintf("invalid parameter %s", k)
			return errors.New(msg)
		}

	}

	// obtain conditions for provided parameter
	if val, ok := a.Params["did"]; ok && val != "" {
		conds, args = AddParam("did", "D.did", a.Params, conds, args)
	}
	if val, ok := a.Params["file"]; ok && val != "" {
		conds, args = AddParam("file", "F.file", a.Params, conds, args)
		tmpl["Files"] = true
	}
	if val, ok := a.Params["script"]; ok && val != "" {
		conds, args = AddParam("script", "SC.name", a.Params, conds, args)
		tmpl["Scripts"] = true
	}
	if val, ok := a.Params["environment"]; ok && val != "" {
		conds, args = AddParam("environment", "E.name", a.Params, conds, args)
		tmpl["Environments"] = true
	}
	if val, ok := a.Params["package"]; ok && val != "" {
		conds, args = AddParam("package", "EP.name", a.Params, conds, args)
		tmpl["Environments"] = true
	}
	if val, ok := a.Params["site"]; ok && val != "" {
		conds, args = AddParam("site", "S.site", a.Params, conds, args)
		tmpl["Sites"] = true
	}
	if val, ok := a.Params["bucket"]; ok && val != "" {
		conds, args = AddParam("bucket", "B.bucket", a.Params, conds, args)
		tmpl["Buckets"] = true
	}
	if val, ok := a.Params["processing"]; ok && val != "" {
		conds, args = AddParam("processing", "PR.processing", a.Params, conds, args)
		tmpl["Processing"] = true
	}
	if val, ok := a.Params["osname"]; ok && val != "" {
		conds, args = AddParam("osname", "O.name", a.Params, conds, args)
		tmpl["Osinfo"] = true
	}

	// get SQL statement from static area
	stm, err := LoadTemplateSQL("select_dataset", tmpl)
	if err != nil {
		return Error(err, LoadErrorCode, "", "dbs.datasets.Datasets")
	}
	stm = WhereClause(stm, conds)

	// use generic query API to fetch the results from DB
	err = executeAll(a.Writer, a.Separator, stm, args...)
	if err != nil {
		return Error(err, QueryErrorCode, "", "dbs.datasets.Datasets")
	}
	return nil
}

// InsertDataset inserts dataset into database
func (a *API) InsertDataset() error {
	// the API provides Reader which will be used by Decode function to load the HTTP payload
	// and cast it to Datasets data structure

	// read given input
	data, err := io.ReadAll(a.Reader)
	if err != nil {
		log.Println("fail to read data", err)
		return Error(err, ReaderErrorCode, "", "dbs.datasets.InsertDataset")
	}
	rec := DatasetRecord{}
	if a.ContentType == "application/json" {
		err = json.Unmarshal(data, &rec)
	} else {
		log.Println("Parser dataset record using default application/json mtime")
		err = json.Unmarshal(data, &rec)
	}
	if err != nil {
		log.Println("reading", a.ContentType)
		log.Println("reading data", string(data))
		log.Println("fail to decode data", err)
		return Error(err, UnmarshalErrorCode, "", "dbs.datasets.InsertDataset")
	}
	err = rec.Validate()
	if err != nil {
		return Error(err, ValidateErrorCode, "validation error", "dbs.datasets.InsertDataset")
	}

	// parse incoming DatasetRequest and insert relationships, e.g.
	// site, bucket, parent, processing, files
	record := Datasets{
		DID:       rec.Did,
		CREATE_BY: a.CreateBy,
		MODIFY_BY: a.CreateBy,
	}
	record.SetDefaults()
	err = record.Validate()
	if err != nil {
		return Error(err, ValidateErrorCode, "validation error", "dbs.datasets.InsertDataset")
	}
	err = insertParts(&rec, &record)
	if err != nil {
		return Error(err, DatasetErrorCode, "fail to insert parts of dataset", "dbs.insertRecord")
	}
	return nil
}

// helper function to insert parts of the dataset relationships
func insertParts(rec *DatasetRecord, record *Datasets) error {
	// start transaction
	tx, err := DB.Begin()
	if err != nil {
		return Error(err, TransactionErrorCode, "", "dbs.insertRecord")
	}
	defer tx.Rollback()
	var siteId, processingId, datasetId, osId, scriptId, fileId, bucketId, configId int64
	var envIds, scriptIds []int64

	// insert site info
	if rec.Site != "" {
		siteId, err = GetID(tx, "sites", "site_id", "site", rec.Site)
		if err != nil || siteId == 0 {
			site := Sites{SITE: rec.Site}
			siteId, err = site.Insert(tx)
			if err != nil {
				return err
			}
		}
	}
	record.SITE_ID = siteId

	// insert os info
	if rec.OsInfo.Name != "" {
		osId, err = GetID(tx, "osinfo", "os_id", "name", rec.OsInfo.Name)
		if err != nil || osId == 0 {
			osId, err = rec.OsInfo.Insert(tx)
			if err != nil {
				return err
			}
		}
		record.OSINFO_ID = osId
	} else {
		// osinfo must be present
		return errors.New("no osinfo is provided")
	}

	// insert environment info
	for _, env := range rec.Environments {
		if env.Name != "" {
			environmentId, err := GetID(tx, "environments", "environment_id", "name", env.Name)
			if err != nil || environmentId == 0 {
				environmentId, err = env.Insert(tx)
				if err != nil {
					return err
				}
			}
			if environmentId == 0 {
				msg := fmt.Sprintf("unable to obtain environment id for %s", env.Name)
				return errors.New(msg)
			}
			envIds = append(envIds, environmentId)
		}
	}

	// insert script info
	for _, script := range rec.Scripts {
		if script.Name != "" {
			scriptId, err = GetID(tx, "scripts", "script_id", "name", script.Name)
			if err != nil || scriptId == 0 {
				scriptId, err = script.Insert(tx)
				if err != nil && !strings.Contains(err.Error(), "UNIQUE") {
					return err
				}
			}
			scriptIds = append(scriptIds, scriptId)
		} else {
			return errors.New("no script info is provide")
		}
	}

	// insert processing info
	if rec.Processing == "" {
		return errors.New("procesing part of provenance records is empty")
	}
	processingId, err = GetID(tx, "processing", "processing_id", "processing", rec.Processing)
	if err != nil || processingId == 0 {
		processing := Processing{
			PROCESSING: rec.Processing,
		}
		processingId, err = processing.Insert(tx)
		if err != nil {
			return err
		}
	}
	record.PROCESSING_ID = processingId

	// insert dataset info
	datasetId, err = GetID(tx, "datasets", "dataset_id", "did", rec.Did)
	if err != nil {
		record.SITE_ID = siteId
		record.PROCESSING_ID = processingId
		datasetId, err = record.Insert(tx)
		if err != nil {
			return err
		}
	}
	record.DATASET_ID = datasetId

	// insert dataset-environments relationships
	for _, envId := range envIds {
		err = InsertManyToMany(tx, "insert_dataset_environment", datasetId, envId)
		if err != nil && !strings.Contains(err.Error(), "UNIQUE") {
			return err
		}
	}
	// insert dataset-scripts relationships
	for _, sid := range scriptIds {
		err = InsertManyToMany(tx, "insert_dataset_script", datasetId, sid)
		if err != nil && !strings.Contains(err.Error(), "UNIQUE") {
			return err
		}
	}
	// insert dataset-configs relationships
	config := Config{CONFIG: rec.Config}
	configId, err = config.Insert(tx)
	err = InsertManyToMany(tx, "insert_dataset_config", datasetId, configId)
	if err != nil && !strings.Contains(err.Error(), "UNIQUE") {
		return err
	}
	record.CONFIG_ID = configId

	// insert parent info
	if rec.Parent != "" {
		// look-up parent_id for given did name (rec.Parent)
		parentId, err := GetID(tx, "datasets", "dataset_id", "did", rec.Parent)
		if err != nil {
			return err
		}
		if parentId == 0 {
			msg := fmt.Sprintf("no did found for did %s", rec.Parent)
			return errors.New(msg)
		}
		// if not found we need to insert relation between parent id and dataset id
		parent := Parents{PARENT_ID: parentId, DATASET_ID: datasetId}
		parentId, err = parent.Insert(tx)
		if err != nil {
			return err
		}
	}

	// insert all buckets
	for _, b := range rec.Buckets {
		bucketId, err = GetID(tx, "buckets", "bucket_id", "bucket", b.Name)
		if err != nil || bucketId == 0 {
			bucket := Buckets{
				BUCKET:     b.Name,
				UUID:       b.UUID,
				META_DATA:  b.MetaData,
				DATASET_ID: datasetId,
			}
			if _, err = bucket.Insert(tx); err != nil {
				log.Printf("Bucket %+v already exist", bucket)
			}
		}
	}

	// insert all input files
	for _, f := range rec.InputFiles {
		fileId, err = GetID(tx, "files", "file_id", "file", f.Name)
		if err != nil {
			fileId, err = f.Insert(tx)
			if err != nil {
				log.Printf("Fail to insert %+v, error %v ", f, err)
			}
		}
		if fileId == 0 {
			msg := fmt.Sprintf("unable to find input file %+v", f)
			return errors.New(msg)
		}
		err = InsertManyToMany(tx, "insert_dataset_file", datasetId, fileId, "input")
		if err != nil && !strings.Contains(err.Error(), "UNIQUE") {
			return err
		}
	}

	// insert all output files
	for _, f := range rec.OutputFiles {
		fileId, err = GetID(tx, "files", "file_id", "file", f.Name)
		if err != nil {
			fileId, err = f.Insert(tx)
			if err != nil {
				log.Printf("Fail to insert %+v, error %v ", f, err)
			}
		}
		if fileId == 0 {
			msg := fmt.Sprintf("unable to find output file %+v", f)
			return errors.New(msg)
		}
		err = InsertManyToMany(tx, "insert_dataset_file", datasetId, fileId, "output")
		if err != nil && !strings.Contains(err.Error(), "UNIQUE") {
			return err
		}
	}

	// commit all transactions
	err = tx.Commit()
	return err
}

func (a *API) UpdateDataset() error {
	// extract payload from API and initialize dataset attributes
	data, err := io.ReadAll(a.Reader)
	if err != nil {
		return err
	}
	rec := &Datasets{}
	return DBOperation("update", rec, data, "dbs.UpdateDatset")
}
func (a *API) DeleteDataset() error {
	// extract payload from API and initialize dataset attributes
	data, err := io.ReadAll(a.Reader)
	if err != nil {
		return err
	}

	rec := &Datasets{}
	return DBOperation("delete", rec, data, "dbs.DeleteDatset")
}

// Delete implementation of Datasets
func (r *Datasets) Delete(tx *sql.Tx) error {
	return nil
}

// Update implementation of Datasets
func (r *Datasets) Update(tx *sql.Tx) error {
	var err error
	if r.DATASET_ID == 0 {
		return Error(err, UpdateErrorCode, "Dataset should have valid id", "dbs.datasets.Update")
	}
	// set defaults and validate the record
	r.SetDefaults()
	err = r.Validate()
	if err != nil {
		log.Println("unable to validate record", err)
		return Error(err, ValidateErrorCode, "", "dbs.datasets.Update")
	}
	// get SQL statement from static area
	stm := getSQL("update_dataset")
	if Verbose > 0 {
		log.Printf("Update Datasets\n%s\n%+v", stm, r)
	}
	// make final SQL statement to insert dataset record
	_, err = tx.Exec(
		stm,
		r.SITE_ID,
		r.PROCESSING_ID,
		r.MODIFY_AT,
		r.MODIFY_BY,
		r.DATASET_ID,
	)
	if err != nil {
		if Verbose > 0 {
			log.Printf("unable to update dataset table %+v", err)
		}
		return Error(err, InsertErrorCode, "", "dbs.datasets.Update")
	}
	return nil
}

// Insert implementation of Datasets
func (r *Datasets) Insert(tx *sql.Tx) (int64, error) {
	var tid int64
	var err error
	if r.DATASET_ID == 0 {
		if DBOWNER == "sqlite" || DBOWNER == "mysql" {
			tid, err = LastInsertID(tx, "datasets", "dataset_id")
			r.DATASET_ID = tid + 1
		} else {
			tid, err = IncrementSequence(tx, "SEQ_DS")
			r.DATASET_ID = tid
		}
	}
	if err != nil {
		return 0, Error(err, LastInsertErrorCode, "", "dbs.datasets.Insert")
	}
	// set defaults and validate the record
	r.SetDefaults()
	err = r.Validate()
	if err != nil {
		log.Println("unable to validate record", err)
		return 0, Error(err, ValidateErrorCode, "", "dbs.datasets.Insert")
	}
	// get SQL statement from static area
	stm := getSQL("insert_dataset")
	if Verbose > 0 {
		log.Printf("Insert Datasets\n%s\n%+v", stm, r)
	}
	// make final SQL statement to insert dataset record
	_, err = tx.Exec(
		stm,
		r.DATASET_ID,
		r.DID,
		r.SITE_ID,
		r.PROCESSING_ID,
		r.OSINFO_ID,
		r.CONFIG_ID,
		r.CREATE_AT,
		r.CREATE_BY,
		r.MODIFY_AT,
		r.MODIFY_BY)
	if err != nil {
		if Verbose > 0 {
			log.Printf("unable to insert Datasets %+v", err)
		}
		return 0, Error(err, InsertErrorCode, "", "dbs.datasets.Insert")
	}
	return r.DATASET_ID, nil
}

// Validate implementation of Datasets
//
//gocyclo:ignore
func (r *Datasets) Validate() error {
	if err := lexicon.CheckPattern("did", r.DID); err != nil {
		return Error(err, ValidateErrorCode, "", "dbs.datasets.Validate")
	}
	if matched := lexicon.UnixTimePattern.MatchString(fmt.Sprintf("%d", r.CREATE_AT)); !matched {
		msg := "invalid pattern for creation date"
		return Error(InvalidParamErr, ValidateErrorCode, msg, "dbs.datasets.Validate")
	}
	if r.CREATE_AT == 0 {
		msg := "missing create_at"
		return Error(InvalidParamErr, ValidateErrorCode, msg, "dbs.datasets.Validate")
	}
	if r.CREATE_BY == "" {
		msg := "missing create_by"
		return Error(InvalidParamErr, ValidateErrorCode, msg, "dbs.datasets.Validate")
	}
	if r.MODIFY_AT == 0 {
		msg := "missing modify_at"
		return Error(InvalidParamErr, ValidateErrorCode, msg, "dbs.datasets.Validate")
	}
	if r.MODIFY_BY == "" {
		msg := "missing modify_by"
		return Error(InvalidParamErr, ValidateErrorCode, msg, "dbs.datasets.Validate")
	}
	return nil
}

// SetDefaults implements set defaults for Datasets
func (r *Datasets) SetDefaults() {
	if r.CREATE_BY == "" {
		r.CREATE_BY = "Server"
	}
	if r.CREATE_AT == 0 {
		r.CREATE_AT = Date()
	}
	if r.MODIFY_BY == "" {
		r.MODIFY_BY = "Server"
	}
	if r.MODIFY_AT == 0 {
		r.MODIFY_AT = Date()
	}
}
