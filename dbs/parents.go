package dbs

// DBS parents module
// Copyright (c) 2023 - Valentin Kuznetsov <vkuznet@gmail.com>
//
// nolint: gocyclo

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
)

// Parents represents Parents DBS DB table
type Parents struct {
	PARENT_ID int64  `json:"parent_id"`
	PARENT    string `json:"parent" validate:"required"`
	CREATE_AT int64  `json:"create_at"`
	CREATE_BY string `json:"create_by"`
	MODIFY_AT int64  `json:"modify_at"`
	MODIFY_BY string `json:"modify_by"`
}

// Parents DBS API
//
//gocyclo:ignore
func (a *API) GetParent() error {
	var args []interface{}
	var conds []string
	var err error

	tmpl := make(map[string]any)
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
		parentID, err := getNextId(tx, "parents", "parent_id")
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
	if Verbose > 0 {
		log.Printf("Insert Parents record %+v", r)
	} else if Verbose > 1 {
		log.Printf("Insert Parents\n%s\n%+v", stm, r)
	}
	_, err = tx.Exec(
		stm,
		r.PARENT_ID,
		r.PARENT,
		r.CREATE_AT,
		r.CREATE_BY,
		r.MODIFY_AT,
		r.MODIFY_BY)
	if err != nil {
		if Verbose > 0 {
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
	if matched := unixTimePattern.MatchString(fmt.Sprintf("%d", r.CREATE_AT)); !matched {
		msg := "invalid pattern for creation date"
		return Error(InvalidParamErr, PatternErrorCode, msg, "dbs.parents.Validate")
	}
	if matched := unixTimePattern.MatchString(fmt.Sprintf("%d", r.MODIFY_AT)); !matched {
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
	if r.CREATE_AT == 0 {
		r.CREATE_AT = Date()
	}
	if r.MODIFY_BY == "" {
		r.MODIFY_BY = "Server"
	}
	if r.MODIFY_AT == 0 {
		r.MODIFY_AT = Date()
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
