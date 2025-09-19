package dbs

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/CHESSComputing/golib/utils"
)

// ProvenanceRecord represent structure of provenance record
type ProvenanceRecord struct {
	Did          string              `json:"did"`
	Site         string              `json:"site"`
	Processing   string              `json:"processing"`
	ParentDid    string              `json:"parent_did"`
	InputFiles   []FileRecord        `json:"input_files"`
	OutputFiles  []FileRecord        `json:"output_files"`
	Environments []EnvironmentRecord `json:"environments"`
	Scripts      []ScriptRecord      `json:"scripts"`
	OsInfo       OsInfoRecord        `json:"osinfo"`
	Buckets      []BucketRecord      `json:"buckets"`
	Config       ConfigRecord        `json:"config"`
}

// IsEmpty checks if given record is empty
func (p *ProvenanceRecord) IsEmpty() bool {
	// we allow Did to be non-empty, everything else must be empty
	if strings.TrimSpace(p.Site) != "" || strings.TrimSpace(p.Processing) != "" || strings.TrimSpace(p.ParentDid) != "" {
		return false
	}
	if len(p.InputFiles) > 0 {
		return false
	}
	if len(p.OutputFiles) > 0 {
		return false
	}
	for _, env := range p.Environments {
		if !env.IsEmpty() {
			return false
		}
	}
	for _, s := range p.Scripts {
		if !s.IsEmpty() {
			return false
		}
	}
	if !p.OsInfo.IsEmpty() {
		return false
	}
	for _, b := range p.Buckets {
		if !b.IsEmpty() {
			return false
		}
	}
	if !p.Config.IsEmpty() {
		return false
	}
	return true
}

//gocyclo:ignore
func (a *API) GetParentDID(did string) (string, error) {
	var args []interface{}
	var conds []string
	tmpl := make(map[string]any)
	tmpl["Owner"] = DBOWNER
	// get SQL statement from static area
	//stm, err := LoadTemplateSQL("select_parent_did", tmpl)
	stm, err := LoadTemplateSQL("select_parent", tmpl)
	if err != nil {
		return "", Error(err, LoadErrorCode, "fail to load select_parent_did sql template", "dbs.provenance.GetParentDID")
	}
	if did != "" {
		args = append(args, did)
		conds = append(conds, "d.did = ?")
	}
	stm = WhereClause(stm, conds)

	tx, err := DB.Begin()
	if err != nil {
		return "", Error(err, TransactionErrorCode, "fail to get DB transaction", "dbs.provenance.GetProvenance")
	}
	defer tx.Rollback()
	rows, err := tx.Query(stm, args...)

	if err != nil {
		return "", err
	}
	defer rows.Close()

	log.Println("QUERY:\n", stm, args)

	var parentDID sql.NullString
	var cat, cby, mat, mby any
	for rows.Next() {
		err := rows.Scan(&did, &parentDID, &cat, &cby, &mat, &mby)
		if err != nil {
			return "", err
		}
		break
	}
	if parentDID.Valid {
		return parentDID.String, nil
	}
	msg := fmt.Sprintf("parent for did %s is not found", did)
	return "", errors.New(msg)
}

//gocyclo:ignore
func (a *API) GetProvenance() error {
	if Verbose > 1 {
		log.Printf("provenance params %+v", a.Params)
	}
	var args []interface{}
	var conds []string
	tmpl := make(map[string]any)
	tmpl["Owner"] = DBOWNER

	allowed := []string{"did"}
	for k, _ := range a.Params {
		if !utils.InList(k, allowed) {
			msg := fmt.Sprintf("invalid parameter %s", k)
			return errors.New(msg)
		}

	}

	var dataset_did string
	if val, ok := a.Params["did"]; ok {
		if val != "" {
			conds, args = AddParam("did", "d.did", a.Params, conds, args)
			dataset_did = fmt.Sprintf("%v", val)
			dataset_did = strings.Replace(dataset_did, "[", "", -1)
			dataset_did = strings.Replace(dataset_did, "]", "", -1)
		}
	} else {
		msg := fmt.Sprintf("/provenance API requires did input, got %+v\n", a.Params)
		return errors.New(msg)
	}

	// get SQL statement from static area
	stm, err := LoadTemplateSQL("select_provenance", tmpl)
	if err != nil {
		return Error(err, LoadErrorCode, "fail to load select_provenance sql template", "dbs.datasets.Datasets")
	}
	stm = WhereClause(stm, conds)
	stm = fmt.Sprintf("%s ORDER BY D.dataset_id, E.environment_id, PK.package_id", stm)

	tx, err := DB.Begin()
	if err != nil {
		return Error(err, TransactionErrorCode, "fail to get DB transaction", "dbs.provenance.GetProvenance")
	}
	defer tx.Rollback()
	rows, err := tx.Query(stm, args...)

	if err != nil {
		return err
	}
	defer rows.Close()

	log.Println("QUERY:\n", stm, args)

	// Process results
	var provenance DatasetRecord

	// keep map of unique packages
	envMap := make(map[int]*EnvironmentRecord)  // Store environments by environment_id
	pkgMap := make(map[int]map[string]struct{}) // Track unique packages per environment
	scriptMap := make(map[int64]*ScriptRecord)  // Store scripts by script_id

	// find parent did
	parentDID, err := a.GetParentDID(dataset_did)
	if err != nil {
		log.Println("WARNING:", err)
	}

	// main query
	for rows.Next() {
		var did, processing, osName, osKernel, osVersion string
		var didSql, processingSql, osNameSql, osKernelSql, osVersionSql sql.NullString
		var bucketName, bucketUUID, bucketMetaData sql.NullString
		var site, scriptName, scriptOptions sql.NullString
		var parentEnvName, parentScript, packageName, packageVersion sql.NullString
		var envName, envVersion, envDetails, envOSName, config sql.NullString
		var scriptID, scriptOrderIdx sql.NullInt64
		var envIDSql sql.NullInt32
		var envID int

		// Scan row into variables
		err := rows.Scan(&didSql, &processingSql, &osNameSql, &osKernelSql, &osVersionSql,
			&envIDSql, &envName, &envVersion, &envDetails, &parentEnvName, &envOSName,
			&packageName, &packageVersion,
			&scriptID, &scriptName, &scriptOrderIdx, &scriptOptions, &parentScript,
			&site, &config, &bucketName, &bucketUUID, &bucketMetaData,
		)
		if err != nil {
			log.Println("ERROR: unable to scan rows", err)
			return err
		}
		if didSql.Valid {
			did = didSql.String
		}
		if envIDSql.Valid {
			envID = int(envIDSql.Int32)
		}
		if processingSql.Valid {
			processing = processingSql.String
		}
		if osNameSql.Valid {
			osName = osNameSql.String
		}
		if osKernelSql.Valid {
			osKernel = osKernelSql.String
		}
		if osVersionSql.Valid {
			osVersion = osVersionSql.String
		}
		// Initialize provenance record if first row
		if provenance.Did == "" {
			provenance = DatasetRecord{
				Did:        did,
				Parent:     parentDID,
				Processing: processing,
				OsInfo: OsInfoRecord{
					Name:    osName,
					Kernel:  osKernel,
					Version: osVersion,
				},
				Environments: []EnvironmentRecord{},
				Site:         site.String,
				Scripts:      []ScriptRecord{},
				Buckets:      []BucketRecord{},
			}
		}

		// config
		if config.Valid {
			provenance.Config = config.String
		}

		// Collect buckets
		b := BucketRecord{}
		if bucketName.Valid {
			b.Name = bucketName.String
		}
		if bucketUUID.Valid {
			b.UUID = bucketUUID.String
		}
		if bucketMetaData.Valid {
			b.MetaData = bucketMetaData.String
		}
		provenance.Buckets = append(provenance.Buckets, b)

		if envOSName.Valid {
			osName = envOSName.String
		}
		// Handle scripts
		var sid int64
		if scriptID.Valid {
			sid = scriptID.Int64
			if _, exists := scriptMap[sid]; !exists {
				scriptMap[sid] = &ScriptRecord{
					Name:     scriptName.String,
					OrderIdx: scriptOrderIdx.Int64,
					Options:  scriptOptions.String,
					Parent:   parentScript.String,
				}
			}
		}

		// Handle environments
		if _, exists := envMap[envID]; !exists {
			envMap[envID] = &EnvironmentRecord{
				Name:     envName.String,
				Version:  envVersion.String,
				Details:  envDetails.String,
				Parent:   parentEnvName.String,
				OSName:   osName,
				Packages: []PackageRecord{},
			}
			pkgMap[envID] = make(map[string]struct{}) // Track unique packages
		}

		// Check if the package is already in the set before adding
		if packageName.Valid && packageVersion.Valid {
			pkgKey := packageName.String + "|" + packageVersion.String
			if _, exists := pkgMap[envID][pkgKey]; !exists {
				envMap[envID].Packages = append(envMap[envID].Packages, PackageRecord{
					Name:    packageName.String,
					Version: packageVersion.String,
				})
				pkgMap[envID][pkgKey] = struct{}{}
			}
		}

	}

	// Convert environments map to list of environments in provenance record
	for _, env := range envMap {
		provenance.Environments = append(provenance.Environments, *env)
	}
	// Convert scripts map to list of scripts in provenance record
	smap := make(map[string]struct{})
	for _, script := range scriptMap {
		if _, exists := smap[script.Name]; !exists {
			provenance.Scripts = append(provenance.Scripts, *script)
			smap[script.Name] = struct{}{}
		}
	}

	// get rid of duplicates
	provenance.Buckets = UniqueBucketRecords(provenance.Buckets)

	// Convert to JSON
	var out []DatasetRecord
	out = append(out, provenance)
	jsonOutput, err := json.MarshalIndent(out, "", "  ")
	if err == nil {
		a.Writer.Write(jsonOutput)
	}
	return err
}

// UniqueBucketRecords removes duplicates from a slice and returns a new slice with unique elements.
func UniqueBucketRecords(bucketRecords []BucketRecord) []BucketRecord {
	seen := make(map[string]bool)
	var result []BucketRecord

	for _, r := range bucketRecords {
		if !seen[r.Name] {
			seen[r.Name] = true
			result = append(result, r)
		}
	}

	return result
}

// helper function to convert user files into FileRecord list
func userFiles(val any) []FileRecord {
	var files []FileRecord
	switch input := val.(type) {
	case string:
		for _, f := range strings.Split(input, " ") {
			files = append(files, FileRecord{Name: strings.Trim(f, " ")})
		}
	case []string:
		for _, f := range input {
			files = append(files, FileRecord{Name: f})
		}
	case []any:
		for _, f := range input {
			files = append(files, FileRecord{Name: fmt.Sprintf("%v", f)})
		}
	}
	return files
}

// InsertProvenance inserts provenance record into DB
func (a *API) InsertProvenance() error {
	// extract payload from API
	data, err := io.ReadAll(a.Reader)
	if err != nil {
		return err
	}
	var userRecord map[string]any
	err = json.Unmarshal(data, &userRecord)
	if err != nil {
		return err
	}

	// parameters for provenance record
	var inputFiles, outputFiles []FileRecord
	var user, did, parentDid, application, site string
	var config any

	// extract all possible values from input user record
	if val, ok := userRecord["input_files"]; ok {
		inputFiles = userFiles(val)
	}
	if val, ok := userRecord["output_files"]; ok {
		outputFiles = userFiles(val)
	}
	if val, ok := userRecord["user"]; ok {
		user = val.(string)
	} else {
		return errors.New("no user value found in user record")
	}
	if val, ok := userRecord["parent_did"]; ok {
		parentDid = val.(string)
		tstamp := time.Now().Format("20060102_150405")
		did = fmt.Sprintf("%s/%s:%s", val, user, tstamp)
	}
	if val, ok := userRecord["did"]; ok {
		did = val.(string)
	}
	if val, ok := userRecord["application"]; ok {
		application = val.(string)
	} else {
		application = "N/A"
	}
	if val, ok := userRecord["site"]; ok {
		site = val.(string)
	} else {
		site = "Cornell"
	}
	if val, ok := userRecord["config"]; ok {
		config = val
	}

	var buckets []BucketRecord
	buckets = append(buckets, BucketRecord{Name: "UserBucket"})

	var environments []EnvironmentRecord
	environments = append(environments, EnvironmentRecord{Name: "UserEnvironment", Version: "N/A", Details: "N/A"})

	var scripts []ScriptRecord
	scripts = append(scripts, ScriptRecord{Name: "UserScript"})

	osinfo := OsInfoRecord{Name: "UserInfo", Version: "N/A", Kernel: "N/A"}

	rec := DatasetRecord{
		Did:          did,
		Site:         site,
		Processing:   application,
		Parent:       parentDid,
		InputFiles:   inputFiles,
		OutputFiles:  outputFiles,
		Environments: environments,
		Scripts:      scripts,
		OsInfo:       osinfo,
		Buckets:      buckets,
		Config:       config,
	}
	if _, ok := userRecord["user_metadata"]; ok {
		fillProvenanceFromUserMetadata(&rec, userRecord)
	}
	record := Datasets{
		DID:       did,
		CREATE_BY: user,
		MODIFY_BY: user,
	}
	record.SetDefaults()
	err = record.Validate()
	if err != nil {
		return Error(err, ValidateErrorCode, "validation error", "dbs.provenance.InsertProvenance")
	}
	err = insertParts(&rec, &record)
	if err != nil {
		return Error(err, DatasetErrorCode, "fail to insert parts of dataset", "dbs.provenance.insertRecord")
	}
	return nil
}

func fillProvenanceFromUserMetadata(rec *DatasetRecord, userRecord map[string]any) {
	if Verbose > 1 {
		log.Println("DEBUG: fill out provenance info from user metadata", userRecord)
	}
	// check if userRecord contains user_record part
	urec, ok := userRecord["user_metadata"]
	if !ok {
		return
	}
	if Verbose > 1 {
		log.Printf("DEBUG: urec=%v type=%T", urec, urec)
	}
	var mrec map[string]any
	switch data := urec.(type) {
	case map[string]any:
		mrec = data
	}
	if len(mrec) == 0 {
		return
	}
	if Verbose > 1 {
		log.Println("DEBUG: mrec", mrec)
	}
	// try to extract from user_record metadata
	prov, ok := mrec["metadata"]
	if !ok {
		return
	}
	if Verbose > 1 {
		log.Printf("DEBUG: prov=%+v type=%T", prov, prov)
	}
	// try to map metadata into provenance record
	var prec ProvenanceRecord
	switch data := prov.(type) {
	case map[string]any:
		bytes, err := json.Marshal(data)
		if err == nil {
			if err := json.Unmarshal(bytes, &prec); err != nil {
				log.Printf("WARNING: unable to map data %v to provenance record, error=%v", data, err)
				return
			}
		}
	case string:
		if err := json.Unmarshal([]byte(data), &prec); err != nil {
			log.Printf("WARNING: unable to map data %v to provenance record, error=%v", data, err)
			return
		}
	}
	if &prec != nil {
		if Verbose > 1 {
			log.Printf("DEBUG: succefully extrated provenance record %+v", prec)
			log.Printf("DEBUG: will fill out provenance record %+v", rec)
		}
		rec.Processing = prec.Processing
		rec.OsInfo = prec.OsInfo
		rec.Environments = prec.Environments
		rec.Scripts = prec.Scripts
		if len(prec.InputFiles) > 0 {
			for _, r := range prec.InputFiles {
				rec.InputFiles = append(rec.InputFiles, r)
			}
		}
		if len(prec.OutputFiles) > 0 {
			for _, r := range prec.OutputFiles {
				rec.OutputFiles = append(rec.OutputFiles, r)
			}
		}
		if Verbose > 1 {
			log.Printf("DEBUG: will fill out provenance record %+v", rec)
		}
	}
}
