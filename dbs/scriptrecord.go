package dbs

import (
	"database/sql"
	"log"
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
