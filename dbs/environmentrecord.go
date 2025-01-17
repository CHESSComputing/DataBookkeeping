package dbs

import (
	"database/sql"

	lexicon "github.com/CHESSComputing/golib/lexicon"
)

// EnvironmentRecord represents data input for environment record
type EnvironmentRecord struct {
	Name     string          `json:"name"`
	Version  string          `json:"version"`
	Details  string          `json:"details"`
	Parent   string          `json:"parent_environment",omitempty`
	OsInfo   string          `json:"osinfo",omitempty`
	Packages []PackageRecord `json:"packages",omitempty`
}

// Insert API
func (e *EnvironmentRecord) Insert(tx *sql.Tx) (int64, error) {
	r := Environments{NAME: e.Name, VERSION: e.Version, DETAILS: e.Details}

	// identify parent script id if parent is present
	if e.Parent != "" {
		parent_environment_id, err := GetID(tx, "environments", "environment_id", "name", e.Parent)
		if err == nil {
			r.PARENT_ENVIRONMENT_ID = parent_environment_id
		} else {
			return 0, err
		}
	}

	// insert env record
	eid, err := r.Insert(tx)

	// insert packages if they are provided
	for _, pkg := range e.Packages {
		p := Packages{NAME: pkg.Name, VERSION: pkg.Version}
		pid, err := p.Insert(tx)
		if err != nil {
			msg := "fail to insert package"
			return 0, Error(err, PackagesErrorCode, msg, "dbs.EnvironmentRecord.Insert")
		}
		err = InsertEnvironmentPackage(eid, pid)
		if err != nil {
			msg := "fail to insert environment-package relationship"
			return 0, Error(err, ManyToManyErrorCode, msg, "dbs.EnvironmentRecord.Insert")
		}
	}
	return eid, err
}

// Validate implementation of EnvironmentRecord
func (r *EnvironmentRecord) Validate() error {
	if err := lexicon.CheckPattern("env_name", r.Name); err != nil {
		return Error(err, ValidateErrorCode, "fail env.Name validation", "dbs.EnvironmentRecord.Validate")
	}
	if err := lexicon.CheckPattern("env_version", r.Version); err != nil {
		return Error(err, ValidateErrorCode, "fail env.Version validation", "dbs.EnvironmentRecord.Validate")
	}
	if err := lexicon.CheckPattern("env_details", r.Details); err != nil {
		return Error(err, ValidateErrorCode, "fail env.Details validation", "dbs.EnvironmentRecord.Validate")
	}
	if err := lexicon.CheckPattern("env_parent", r.Parent); err != nil {
		return Error(err, ValidateErrorCode, "fail env.Parent validation", "dbs.EnvironmentRecord.Validate")
	}
	return nil
}
