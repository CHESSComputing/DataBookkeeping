package dbs

import (
	"database/sql"
	"fmt"
)

// InsertManyToMany provides function to insert many-to-many relationship from given
// template name and set of parameters
func InsertManyToMany(tx *sql.Tx, tmplName string, args ...interface{}) error {
	stm := getSQL(tmplName)
	_, err := tx.Exec(stm, args...)
	if err != nil {
		msg := fmt.Sprintf("fail to insert %s template", tmplName)
		return Error(err, InsertErrorCode, msg, "dbs.manytomany.InsertManyToMany")
	}
	return nil
}
