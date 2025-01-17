package dbs

import (
	"database/sql"
	"log"

	lexicon "github.com/CHESSComputing/golib/lexicon"
)

// PackageRecord represents data input for package record
type PackageRecord struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// Insert API
func (e *PackageRecord) Insert(tx *sql.Tx) (int64, error) {
	r := Packages{NAME: e.Name, VERSION: e.Version}
	if r.PACKAGE_ID == 0 {
		id, err := getNextId(tx, "packages", "package_id")
		if err != nil {
			log.Println("unable to get package id", err)
			return 0, Error(err, ParametersErrorCode, "", "dbs.package.Insert")
		}
		r.PACKAGE_ID = id
	}
	err := r.Insert(tx)
	return r.PACKAGE_ID, err
}

// Validate implementation of PackageRecord
func (r *PackageRecord) Validate() error {
	if err := lexicon.CheckPattern("env_name", r.Name); err != nil {
		return Error(err, PatternErrorCode, "fail env.Name validation", "dbs.datasets.DatasetRecord.Validate")
	}
	if err := lexicon.CheckPattern("env_version", r.Version); err != nil {
		return Error(err, PatternErrorCode, "fail env.Version validation", "dbs.datasets.DatasetRecord.Validate")
	}
	return nil
}
