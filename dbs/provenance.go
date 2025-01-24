package dbs

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/CHESSComputing/golib/utils"
)

//gocyclo:ignore
func (a *API) GetParentDID(did string) (string, error) {
	var args []interface{}
	var conds []string
	tmpl := make(map[string]any)
	tmpl["Owner"] = DBOWNER
	// get SQL statement from static area
	stm, err := LoadTemplateSQL("select_parent_did", tmpl)
	if err != nil {
		return "", Error(err, LoadErrorCode, "fail to load select_parent_did sql template", "dbs.provenance.GetParentDID")
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
	for rows.Next() {
		err := rows.Scan(&did, &parentDID)
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
	stm = fmt.Sprintf("%s ORDER BY d.dataset_id, e.environment_id, pk.package_id", stm)

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
		var bucketName, bucketUUID, bucketMetaData sql.NullString
		var site, scriptName, scriptOptions sql.NullString
		var parentEnvName, parentScript, packageName, packageVersion sql.NullString
		var envName, envVersion, envDetails, envOSName sql.NullString
		var scriptID, scriptOrderIdx sql.NullInt64
		var envID int

		// Scan row into variables
		err := rows.Scan(&did, &processing, &osName, &osKernel, &osVersion,
			&envID, &envName, &envVersion, &envDetails, &parentEnvName, &envOSName,
			&packageName, &packageVersion,
			&scriptID, &scriptName, &scriptOrderIdx, &scriptOptions, &parentScript,
			&site, &bucketName, &bucketUUID, &bucketMetaData,
		)
		if err != nil {
			return err
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
