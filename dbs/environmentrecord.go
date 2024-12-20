package dbs

import (
	"database/sql"
	"log"
)

// EnvironmentRecord represents data input for environment record
type EnvironmentRecord struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Details string `json:"details"`
	Parent  string `json:"parent_environment",omitempty`
}

// Insert API
func (e *EnvironmentRecord) Insert(tx *sql.Tx) (int64, error) {
	r := Environments{NAME: e.Name, VERSION: e.Version, DETAILS: e.Details}
	if r.ENVIRONMENT_ID == 0 {
		id, err := getNextId(tx, "environments", "environment_id")
		if err != nil {
			log.Println("unable to get environment id", err)
			return 0, Error(err, ParametersErrorCode, "", "dbs.environment.Insert")
		}
		r.ENVIRONMENT_ID = id
	}
	// identify parent script id if parent is present
	if e.Parent != "" {
		parent_environment_id, err := GetID(tx, "environments", "environment_id", "name", e.Parent)
		if err == nil {
			r.PARENT_ENVIRONMENT_ID = parent_environment_id
		} else {
			return 0, err
		}
	}
	err := r.Insert(tx)
	return r.ENVIRONMENT_ID, err
}
