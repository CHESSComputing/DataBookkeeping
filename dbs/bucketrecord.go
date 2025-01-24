package dbs

import (
	"database/sql"

	lexicon "github.com/CHESSComputing/golib/lexicon"
)

// BucketRecord represents bucket input data record
type BucketRecord struct {
	Name     string `json:"name"`
	UUID     string `json:"uuid"`
	MetaData string `json:"meta_data"`
}

// Insert API
func (e *BucketRecord) Insert(tx *sql.Tx) (int64, error) {
	r := Buckets{BUCKET: e.Name, UUID: e.UUID, META_DATA: e.MetaData}
	bid, err := r.Insert(tx)
	return bid, err
}

// Validate implementation of BucketRecord
func (r *BucketRecord) Validate() error {
	if err := lexicon.CheckPattern("bucket_name", r.Name); err != nil {
		return Error(err, ValidateErrorCode, "fail bucket.Name validation", "dbs.BucketRecord.Validate")
	}
	return nil
}
