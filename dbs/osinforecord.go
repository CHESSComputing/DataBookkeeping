package dbs

import (
	"database/sql"

	lexicon "github.com/CHESSComputing/golib/lexicon"
)

// OsInfoRecord represent input os info record
type OsInfoRecord struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Kernel  string `json:"kernel"`
}

// Insert API
func (o *OsInfoRecord) Insert(tx *sql.Tx) (int64, error) {
	r := OsInfo{NAME: o.Name, VERSION: o.Version, KERNEL: o.Kernel}
	oid, err := r.Insert(tx)
	return oid, err
}

// Validate implementation of OsInfoRecord
func (r *OsInfoRecord) Validate() error {
	if err := lexicon.CheckPattern("osinfo_name", r.Name); err != nil {
		return Error(err, PatternErrorCode, "fail osinfo.name validation", "dbs.OsInfoRecord.Validate")
	}
	if err := lexicon.CheckPattern("osinfo_version", r.Version); err != nil {
		return Error(err, PatternErrorCode, "fail osinfo.version validation", "dbs.OsInfoRecord.Validate")
	}
	if err := lexicon.CheckPattern("osinfo_kernel", r.Kernel); err != nil {
		return Error(err, PatternErrorCode, "fail osinfo.name validation", "dbs.OsInfoRecord.Validate")
	}
	return nil
}
