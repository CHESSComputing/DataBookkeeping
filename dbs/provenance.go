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
func (a *API) GetProvenance() error {
	if Verbose > 1 {
		log.Printf("datasets params %+v", a.Params)
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

	if val, ok := a.Params["did"]; ok {
		if val != "" {
			conds, args = AddParam("did", "d.did", a.Params, conds, args)
		}
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

	log.Println("SQL", stm, args)

	// Process results
	var provenance DatasetRecord
	//     provenance := DatasetRecord{}
	//     provenance.Buckets = []string{}
	//     provenance.InputFiles = []string{}
	//     provenance.OutputFiles = []string{}

	// keep map of unique packages
	envMap := make(map[int]*EnvironmentRecord)  // Store environments by environment_id
	pkgMap := make(map[int]map[string]struct{}) // Track unique packages per environment

	for rows.Next() {
		var did, processing, osName, osKernel, osVersion string
		var parentDID, fileType, fileName, bucketName sql.NullString
		var site, scriptName, scriptOptions sql.NullString
		var parentEnvName, parentScript, parentScriptName, packageName, packageVersion sql.NullString
		var envName, envVersion, envDetails, envOSName sql.NullString
		var envID int

		// Scan row into variables
		err := rows.Scan(&did, &parentDID, &processing, &osName, &osKernel, &osVersion,
			&envID, &envName, &envVersion, &envDetails, &parentEnvName, &envOSName,
			&packageName, &packageVersion,
			&scriptName, &scriptOptions, &parentScript,
			&site, &fileName, &fileType, &bucketName,
		)
		if err != nil {
			return err
		}
		// Initialize provenance record if first row
		if provenance.Did == "" {
			provenance = DatasetRecord{
				Did:        did,
				Processing: processing,
				OsInfo: OsInfoRecord{
					Name:    osName,
					Kernel:  osKernel,
					Version: osVersion,
				},
				Environments: []EnvironmentRecord{},
				Site:         site.String,
				Script: ScriptRecord{
					Name:    scriptName.String,
					Options: scriptOptions.String,
					Parent:  parentScriptName.String,
				},
				InputFiles:  []string{},
				OutputFiles: []string{},
				Buckets:     []string{},
			}
		}

		// Handle nullable values
		if parentDID.Valid {
			provenance.Parent = parentDID.String
		}
		if parentScript.Valid {
			provenance.Script.Parent = parentScript.String
		}

		// Collect input/output files
		if fileType.Valid && fileName.Valid {
			if fileType.String == "input" {
				provenance.InputFiles = append(provenance.InputFiles, fileName.String)
			} else if fileType.String == "output" {
				provenance.OutputFiles = append(provenance.OutputFiles, fileName.String)
			}
		}

		// Collect buckets
		if bucketName.Valid {
			provenance.Buckets = append(provenance.Buckets, bucketName.String)
		}

		// Handle environments
		if envOSName.Valid {
			osName = envOSName.String
		}
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

	// Convert environments map to slice
	for _, env := range envMap {
		provenance.Environments = append(provenance.Environments, *env)
	}

	// get rid of duplicates in files
	provenance.InputFiles = UniqueList(provenance.InputFiles)
	provenance.OutputFiles = UniqueList(provenance.OutputFiles)

	// Convert to JSON
	jsonOutput, err := json.MarshalIndent(provenance, "", "  ")
	if err == nil {
		a.Writer.Write(jsonOutput)
	}
	return err
}
