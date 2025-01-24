package dbs

import (
	"database/sql"

	lexicon "github.com/CHESSComputing/golib/lexicon"
)

// ScriptRecord represents script input data record
type ScriptRecord struct {
	Name     string `json:"name"`
	Options  string `json:"options"`
	Parent   string `json:"parent_script"`
	OrderIdx int64  `json:"order_idx"`
}

// Insert API
func (e *ScriptRecord) Insert(tx *sql.Tx) (int64, error) {
	r := Scripts{NAME: e.Name, OPTIONS: e.Options, ORDER_IDX: e.OrderIdx}
	// identify parent script id if parent is present
	if e.Parent != "" {
		parent_script_id, err := GetID(tx, "scripts", "script_id", "name", e.Parent)
		if err == nil {
			r.PARENT_SCRIPT_ID = parent_script_id
		} else {
			return 0, err
		}
	}
	sid, err := r.Insert(tx)
	return sid, err
}

// Validate implementation of ScriptRecord
func (r *ScriptRecord) Validate() error {
	if err := lexicon.CheckPattern("script_name", r.Name); err != nil {
		return Error(err, ValidateErrorCode, "fail script.Name validation", "dbs.ScriptRecord.Validate")
	}
	if err := lexicon.CheckPattern("script_options", r.Options); err != nil {
		return Error(err, ValidateErrorCode, "fail script.Options validation", "dbs.ScriptRecord.Validate")
	}
	if err := lexicon.CheckPattern("script_parent", r.Parent); err != nil {
		return Error(err, ValidateErrorCode, "fail script.Parent validation", "dbs.ScriptRecord.Validate")
	}
	return nil
}
