package dbs

import (
	"database/sql"
	"fmt"
)

// ConfigRecord represent input config record
type ConfigRecord struct {
	Content any `json:"content"`
}

// IsEmpty checks if given record is empty
func (c *ConfigRecord) IsEmpty() bool {
	return c.Content == nil || fmt.Sprintf("%v", c.Content) == ""
}

// Insert API
func (o *ConfigRecord) Insert(tx *sql.Tx) (int64, error) {
	r := Config{CONTENT: o.Content}
	cid, err := r.Insert(tx)
	return cid, err
}

// Validate implementation of ConfigRecord
func (r *ConfigRecord) Validate() error {
	return nil
}
