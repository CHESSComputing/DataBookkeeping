package dbs

import (
	"database/sql"
	"log"

	lexicon "github.com/CHESSComputing/golib/lexicon"
)

// ScriptRecord represents script input data record
type ScriptRecord struct {
	Name    string `json:"name"`
	Options string `json:"options"`
	Parent  string `json:"parent_script",omitempty`
}

// Insert API
func (e *ScriptRecord) Insert(tx *sql.Tx) (int64, error) {
	r := Scripts{NAME: e.Name, OPTIONS: e.Options}
	if r.SCRIPT_ID == 0 {
		id, err := getNextId(tx, "scripts", "script_id")
		if err != nil {
			log.Println("unable to get script id", err)
			return 0, Error(err, ParametersErrorCode, "", "dbs.script.Insert")
		}
		r.SCRIPT_ID = id
	}
	// identify parent script id if parent is present
	if e.Parent != "" {
		parent_script_id, err := GetID(tx, "scripts", "script_id", "name", e.Parent)
		if err == nil {
			r.PARENT_SCRIPT_ID = parent_script_id
		} else {
			return 0, err
		}
	}
	err := r.Insert(tx)
	return r.SCRIPT_ID, err
}

// Validate implementation of ScriptRecord
func (r *ScriptRecord) Validate() error {
	if err := lexicon.CheckPattern("script_name", r.Name); err != nil {
		return Error(err, PatternErrorCode, "fail script.Name validation", "dbs.datasets.DatasetRecord.Validate")
	}
	if err := lexicon.CheckPattern("script_options", r.Options); err != nil {
		return Error(err, PatternErrorCode, "fail script.Options validation", "dbs.datasets.DatasetRecord.Validate")
	}
	if err := lexicon.CheckPattern("script_parent", r.Parent); err != nil {
		return Error(err, PatternErrorCode, "fail script.Parent validation", "dbs.datasets.DatasetRecord.Validate")
	}
	return nil
}
