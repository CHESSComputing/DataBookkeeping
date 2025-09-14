package dbs

import (
	"database/sql"

	lexicon "github.com/CHESSComputing/golib/lexicon"
)

// FileRecord represents script input data record
type FileRecord struct {
	Name     string `json:"name"`
	Checksum string `json:"checksum,omitempty"`
	Size     int64  `json:"size,omitempty"`
	IsValid  int64  `json:"isvalid,omitempty"`
}

// IsEmpty checks if given record is empty
func (r *FileRecord) IsEmpty() bool {
	return r.Name == "" && r.Checksum == "" && r.Size == 0 && r.IsValid == 0
}

// IsEmpty checks if given record is empty
// Insert API
func (r *FileRecord) Insert(tx *sql.Tx) (int64, error) {
	f := Files{
		FILE:          r.Name,
		CHECKSUM:      r.Checksum,
		SIZE:          r.Size,
		IS_FILE_VALID: 1, // when insert new file into always set is file valid to 1
		CREATE_BY:     "Server",
		CREATE_AT:     Date(),
		MODIFY_BY:     "Server",
		MODIFY_AT:     Date(),
	}
	fid, err := f.Insert(tx)
	return fid, err
}

// Validate implementation of FileRecord
func (r *FileRecord) Validate() error {
	if err := lexicon.CheckPattern("file", r.Name); err != nil {
		return Error(err, ValidateErrorCode, "fail file validation", "dbs.FileRecord.Validate")
	}
	return nil
}
