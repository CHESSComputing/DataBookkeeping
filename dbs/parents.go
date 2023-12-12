package dbs

// nolint: gocyclo

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/CHESSComputing/DataBookkeeping/utils"
)

// Parents represents Parents DBS DB table
type Parents struct {
	PARENT_ID              int64  `json:"parent_id"`
	PARENT                 string `json:"parent" validate:"required"`
	CREATION_DATE          int64  `json:"creation_date"`
	CREATE_BY              string `json:"create_by"`
	LAST_MODIFICATION_DATE int64  `json:"last_modification_date"`
	LAST_MODIFIED_BY       string `json:"last_modified_by"`
}

// Parents DBS API
//
//gocyclo:ignore
func (a *API) GetParent() error {
	var args []interface{}
	var conds []string
	var err error

	tmpl := make(Record)
	tmpl["Owner"] = DBOWNER
	stm, err := LoadTemplateSQL("select_parent", tmpl)
	if err != nil {
		return Error(err, LoadErrorCode, "", "dbs.parents.Parents")
	}

	stm = WhereClause(stm, conds)

	// use generic query API to fetch the results from DB
	err = executeAll(a.Writer, a.Separator, stm, args...)
	if err != nil {
		return Error(err, QueryErrorCode, "", "dbs.parents.Parents")
	}
	return nil
}

// InsertParent inserts parent record into DB
func (a *API) InsertParent() error {
	// the API provides Reader which will be used by Decode function to load the HTTP payload
	// and cast it to Parents data structure
	return insertRecord(&Parents{}, a.Reader)
}

// UpdateParent inserts parent record in DB
func (a *API) UpdateParent() error {
	return nil
}

// DeleteParent deletes parent record in DB
func (a *API) DeleteParent() error {
	return nil
}

// Insert implementation of Parents
func (r *Parents) Insert(tx *sql.Tx) error {
	var err error
	if r.PARENT_ID == 0 {
		parentID, err := getNextId(tx, "PARENTS", "PARENT_ID")
		if err != nil {
			log.Println("unable to get parentID", err)
			return Error(err, ParametersErrorCode, "", "dbs.parents.Insert")
		}
		r.PARENT_ID = parentID
	}
	// set defaults and validate the record
	r.SetDefaults()
	err = r.Validate()
	if err != nil {
		log.Println("unable to validate record", err)
		return Error(err, ValidateErrorCode, "", "dbs.parents.Insert")
	}
	// get SQL statement from static area
	stm := getSQL("insert_parent")
	if utils.VERBOSE > 0 {
		log.Printf("Insert Parents record %+v", r)
	} else if utils.VERBOSE > 1 {
		log.Printf("Insert Parents\n%s\n%+v", stm, r)
	}
	_, err = tx.Exec(
		stm,
		r.PARENT_ID,
		r.PARENT,
		r.CREATION_DATE,
		r.CREATE_BY,
		r.LAST_MODIFICATION_DATE,
		r.LAST_MODIFIED_BY)
	if err != nil {
		if utils.VERBOSE > 0 {
			log.Println("unable to insert parents, error", err)
		}
		return Error(err, InsertErrorCode, "", "dbs.parents.Insert")
	}
	return nil
}

// Validate implementation of Parents
func (r *Parents) Validate() error {
	if err := RecordValidator.Struct(*r); err != nil {
		return DecodeValidatorError(r, err)
	}
	if matched := unixTimePattern.MatchString(fmt.Sprintf("%d", r.CREATION_DATE)); !matched {
		msg := "invalid pattern for creation date"
		return Error(InvalidParamErr, PatternErrorCode, msg, "dbs.parents.Validate")
	}
	if matched := unixTimePattern.MatchString(fmt.Sprintf("%d", r.LAST_MODIFICATION_DATE)); !matched {
		msg := "invalid pattern for last modification date"
		return Error(InvalidParamErr, PatternErrorCode, msg, "dbs.parents.Validate")
	}
	return nil
}

// SetDefaults implements set defaults for Parents
func (r *Parents) SetDefaults() {
	if r.CREATE_BY == "" {
		r.CREATE_BY = "Server"
	}
	if r.CREATION_DATE == 0 {
		r.CREATION_DATE = Date()
	}
	if r.LAST_MODIFIED_BY == "" {
		r.LAST_MODIFIED_BY = "Server"
	}
	if r.LAST_MODIFICATION_DATE == 0 {
		r.LAST_MODIFICATION_DATE = Date()
	}
}

// Decode implementation for Parents
func (r *Parents) Decode(reader io.Reader) error {
	// init record with given data record
	data, err := io.ReadAll(reader)
	if err != nil {
		log.Println("fail to read data", err)
		return Error(err, ReaderErrorCode, "", "dbs.parents.Decode")
	}
	err = json.Unmarshal(data, &r)

	//     decoder := json.NewDecoder(r)
	//     err := decoder.Decode(&rec)
	if err != nil {
		log.Println("fail to decode data", err)
		return Error(err, UnmarshalErrorCode, "", "dbs.parents.Decode")
	}
	return nil
}
