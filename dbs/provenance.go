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

	// Process results
	provenance := DatasetRecord{}
	provenance.Buckets = []string{}
	provenance.InputFiles = []string{}
	provenance.OutputFiles = []string{}
	provenance.Environment.Packages = []PackageRecord{}

	for rows.Next() {
		var parentDID sql.NullString
		var fileType sql.NullString
		var fileName sql.NullString
		var bucketName sql.NullString
		var parentEnv sql.NullString
		var parentScript sql.NullString
		var packageName sql.NullString
		var packageVersion sql.NullString

		// Scan row into variables
		err := rows.Scan(
			&provenance.Did,
			&parentDID,
			&provenance.Processing,
			&provenance.OsInfo.Name,
			&provenance.OsInfo.Kernel,
			&provenance.OsInfo.Version,
			&provenance.Environment.Name,
			&provenance.Environment.Version,
			&provenance.Environment.Details,
			&parentEnv,
			&provenance.Environment.OsInfo,
			&packageName,
			&packageVersion,
			&provenance.Script.Name,
			&provenance.Script.Options,
			&parentScript,
			&provenance.Site,
			&fileName,
			&fileType,
			&bucketName,
		)
		if err != nil {
			return err
		}

		// Handle nullable values
		if parentDID.Valid {
			provenance.Parent = parentDID.String
		}
		if parentEnv.Valid {
			provenance.Environment.Parent = parentEnv.String
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

		// Collect packages
		if provenance.Environment.Packages == nil {
			provenance.Environment.Packages = []PackageRecord{}
		}
		provenance.Environment.Packages = append(provenance.Environment.Packages, PackageRecord{
			Name:    provenance.Environment.Name,
			Version: provenance.Environment.Version,
		})
	}

	// Convert to JSON
	jsonOutput, err := json.MarshalIndent(provenance, "", "  ")
	if err == nil {
		a.Writer.Write(jsonOutput)
	}
	return err
}
