package dbs

import (
	"database/sql"
	"log"
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
	if r.OS_ID == 0 {
		id, err := getNextId(tx, "osinfo", "os_id")
		if err != nil {
			log.Println("unable to get osinfo id", err)
			return 0, Error(err, ParametersErrorCode, "", "dbs.osinfo.Insert")
		}
		r.OS_ID = id
	}
	err := r.Insert(tx)
	return r.OS_ID, err
}
