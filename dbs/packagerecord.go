package dbs

import (
	"database/sql"

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
	pid, err := r.Insert(tx)
	msg := "unable to insert packages record"
	return pid, Error(err, PackagesErrorCode, msg, "dbs.PackageRecord.Insert")
}

// Validate implementation of PackageRecord
func (r *PackageRecord) Validate() error {
	if err := lexicon.CheckPattern("env_name", r.Name); err != nil {
		return Error(err, ValidateErrorCode, "fail env.Name validation", "dbs.datasets.DatasetRecord.Validate")
	}
	if err := lexicon.CheckPattern("env_version", r.Version); err != nil {
		return Error(err, ValidateErrorCode, "fail env.Version validation", "dbs.datasets.DatasetRecord.Validate")
	}
	return nil
}
