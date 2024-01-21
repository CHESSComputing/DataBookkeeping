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

	"github.com/CHESSComputing/golib/utils"
)

// Datasets represents Datasets DBS DB table
type Datasets struct {
	DATASET_ID    int64  `json:"dataset_id"`
	DATASET       string `json:"dataset" validate:"required"`
	META_ID       int64  `json:"meta_id" validate:"required"`
	SITE_ID       int64  `json:"site_id" validate:"required"`
	PROCESSING_ID int64  `json:"processing_id" validate:"required"`
	PARENT_ID     int64  `json:"parent_id" validate:"required"`
	CREATE_AT     int64  `json:"create_at" validate:"required,number"`
	CREATE_BY     string `json:"create_by" validate:"required"`
	MODIFY_AT     int64  `json:"modify_at" validate:"required,number"`
	MODIFY_BY     string `json:"modify_by" validate:"required"`
}

// DatasetRecord represents input dataset record from HTTP request
type DatasetRecord struct {
	Dataset    string   `json:"dataset" validate:"required"`
	Buckets    []string `json:"buckets" validate:"required"`
	Site       string   `json:"site" validate:"required"`
	Processing string   `json:"processing" validate:"required"`
	Parent     string   `json:"parent_dataset" validate:"required"`
	MetaId     string   `json:"meta_id" validate:"required"`
	Files      []string `json:"files" validate:"required"`
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

	allowed := []string{"dataset", "did"}
	for k, _ := range a.Params {
		if !utils.InList(k, allowed) {
			msg := fmt.Sprintf("invalid parameter %s", k)
			return errors.New(msg)
		}

	}

	if val, ok := a.Params["dataset"]; ok {
		if val != "" {
			conds, args = AddParam("dataset", "D.dataset", a.Params, conds, args)
		}
	}
	if val, ok := a.Params["did"]; ok {
		if val != "" {
			conds, args = AddParam("did", "M.did", a.Params, conds, args)
		}
	}
	if Verbose > 0 {
		log.Println("### /dataset params", a.Params, conds, args)
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
	log.Printf("### input DatasetRecord %+v", rec)

	// parse incoming DatasetRequest and insert relationships, e.g.
	// site, bucket, parent, processing, files
	record := Datasets{
		DATASET:   rec.Dataset,
		CREATE_BY: a.CreateBy,
		MODIFY_BY: a.CreateBy,
	}
	err = insertParts(&rec, &record)
	if err != nil {
		return Error(err, CommitErrorCode, "", "dbs.insertRecord")
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
	var siteId, processingId, parentId, datasetId, metaId int64

	// insert site info
	siteId, err = GetID(tx, "sites", "site_id", "site", rec.Site)
	if err != nil {
		site := Sites{SITE: rec.Site}
		if err = site.Insert(tx); err != nil {
			return err
		}
		siteId, err = GetID(tx, "sites", "site_id", "site", rec.Site)
		if err != nil {
			return err
		}
	}
	record.SITE_ID = siteId

	// insert meta info
	metaId, err = GetID(tx, "metadata", "meta_id", "did", rec.MetaId)
	if err != nil {
		meta := Metadata{DID: rec.MetaId}
		if err = meta.Insert(tx); err != nil {
			return err
		}
		metaId, err = GetID(tx, "metadata", "meta_id", "did", rec.MetaId)
		if err != nil {
			return err
		}
	}
	record.META_ID = metaId

	// insert processing info
	processingId, err = GetID(tx, "processing", "processing_id", "processing", rec.Processing)
	if err != nil {
		processing := Processing{PROCESSING: rec.Processing}
		if err = processing.Insert(tx); err != nil {
			return err
		}
		processingId, err = GetID(tx, "processing", "processing_id", "processing", rec.Processing)
		if err != nil {
			return err
		}
	}
	record.PROCESSING_ID = processingId

	// insert parent info
	parentId, err = GetID(tx, "parents", "parent_id", "parent", rec.Parent)
	if err != nil {
		if rec.Parent != "" {
			parent := Parents{PARENT: rec.Parent}
			if err = parent.Insert(tx); err != nil {
				return err
			}
			parentId, err = GetID(tx, "parents", "parent_id", "parent", rec.Parent)
			if err != nil {
				return err
			}
		}
	}
	record.PARENT_ID = parentId

	// insert dataset info
	datasetId, err = GetID(tx, "datasets", "dataset_id", "dataset", rec.Dataset)
	if err != nil {
		record.SITE_ID = siteId
		record.PARENT_ID = parentId
		record.PROCESSING_ID = processingId
		if err = record.Insert(tx); err != nil {
			return err
		}
		datasetId, err = GetID(tx, "datasets", "dataset_id", "dataset", rec.Dataset)
		if err != nil {
			return err
		}
	}

	// insert all buckets
	for _, b := range rec.Buckets {
		bucket := Buckets{
			BUCKET:     b,
			DATASET_ID: datasetId,
			META_ID:    metaId,
		}
		if err = bucket.Insert(tx); err != nil {
			log.Printf("Bucket %+v already exist", bucket)
		}
	}

	// insert all files
	for _, f := range rec.Files {
		file := Files{
			FILE:          f,
			IS_FILE_VALID: 1, // by default all files are valid
			DATASET_ID:    datasetId,
			META_ID:       metaId,
			CREATE_BY:     record.CREATE_BY,
			MODIFY_BY:     record.CREATE_BY,
		}
		if err = file.Insert(tx); err != nil {
			log.Printf("File %+v already exist", file)
		}
	}

	// commit all transactions
	err = tx.Commit()
	return err
}

func (a *API) UpdateDataset() error {
	return nil
}
func (a *API) DeleteDataset() error {
	return nil
}

// Insert implementation of Datasets
func (r *Datasets) Insert(tx *sql.Tx) error {
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
		return Error(err, LastInsertErrorCode, "", "dbs.datasets.Insert")
	}
	// set defaults and validate the record
	r.SetDefaults()
	err = r.Validate()
	if err != nil {
		log.Println("unable to validate record", err)
		return Error(err, ValidateErrorCode, "", "dbs.datasets.Insert")
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
		r.DATASET,
		r.META_ID,
		r.SITE_ID,
		r.PROCESSING_ID,
		r.PARENT_ID,
		r.CREATE_AT,
		r.CREATE_BY,
		r.MODIFY_AT,
		r.MODIFY_BY)
	if err != nil {
		if Verbose > 0 {
			log.Printf("unable to insert Datasets %+v", err)
		}
		return Error(err, InsertErrorCode, "", "dbs.datasets.Insert")
	}
	return nil
}

// Validate implementation of Datasets
//
//gocyclo:ignore
func (r *Datasets) Validate() error {
	if err := CheckPattern("dataset", r.DATASET); err != nil {
		return Error(err, PatternErrorCode, "", "dbs.datasets.Validate")
	}
	if matched := unixTimePattern.MatchString(fmt.Sprintf("%d", r.CREATE_AT)); !matched {
		msg := "invalid pattern for creation date"
		return Error(InvalidParamErr, PatternErrorCode, msg, "dbs.datasets.Validate")
	}
	if r.CREATE_AT == 0 {
		msg := "missing create_at"
		return Error(InvalidParamErr, ParametersErrorCode, msg, "dbs.datasets.Validate")
	}
	if r.CREATE_BY == "" {
		msg := "missing create_by"
		return Error(InvalidParamErr, ParametersErrorCode, msg, "dbs.datasets.Validate")
	}
	if r.MODIFY_AT == 0 {
		msg := "missing modify_at"
		return Error(InvalidParamErr, ParametersErrorCode, msg, "dbs.datasets.Validate")
	}
	if r.MODIFY_BY == "" {
		msg := "missing modify_by"
		return Error(InvalidParamErr, ParametersErrorCode, msg, "dbs.datasets.Validate")
	}
	return nil
}

// SetDefaults implements set defaults for Datasets
func (r *Datasets) SetDefaults() {
	if r.CREATE_AT == 0 {
		r.CREATE_AT = Date()
	}
	if r.MODIFY_AT == 0 {
		r.MODIFY_AT = Date()
	}
}

// Decode implementation for Datasets
func (r *Datasets) Decode(reader io.Reader) error {
	if reader == nil {
		return nil
	}
	// init record with given data record
	data, err := io.ReadAll(reader)
	if err != nil {
		log.Println("fail to read data", err)
		return Error(err, ReaderErrorCode, "", "dbs.datasets.Decode")
	}
	err = json.Unmarshal(data, &r)

	//     decoder := json.NewDecoder(r)
	//     err := decoder.Decode(&rec)
	if err != nil {
		log.Println("fail to decode data", err)
		return Error(err, UnmarshalErrorCode, "", "dbs.datasets.Decode")
	}
	return nil
}
