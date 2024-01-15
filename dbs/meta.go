package dbs

// DBS metadata module
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

// Metadata represents Metadata DBS DB table
type Metadata struct {
	META_ID   int64  `json:"meta_id"`
	DID       string `json:"did" validate:"required"`
	CREATE_AT int64  `json:"create_at"`
	CREATE_BY string `json:"create_by"`
	MODIFY_AT int64  `json:"modify_at"`
	MODIFY_BY string `json:"modify_by"`
}

// metadata DBS API
//
//gocyclo:ignore
func (a *API) GetMeta() error {
	var args []interface{}
	var conds []string
	var err error

	tmpl := make(map[string]any)
	tmpl["Owner"] = DBOWNER
	stm, err := LoadTemplateSQL("select_meta", tmpl)
	if err != nil {
		return Error(err, LoadErrorCode, "", "dbs.metadata.Metadata")
	}

	stm = WhereClause(stm, conds)

	// use generic query API to fetch the results from DB
	err = executeAll(a.Writer, a.Separator, stm, args...)
	if err != nil {
		return Error(err, QueryErrorCode, "", "dbs.meta.GetMeta")
	}
	return nil
}

// InsertMetadata inserts metadata record into DB
func (a *API) InsertMetadata() error {
	// the API provides Reader which will be used by Decode function to load the HTTP payload
	// and cast it to Metadata data structure
	return insertRecord(&Metadata{}, a.Reader)
}

// UpdateMetadata inserts metadata record in DB
func (a *API) UpdateMetadata() error {
	return nil
}

// DeleteMetadata deletes metadata record in DB
func (a *API) DeleteMetadata() error {
	return nil
}

// Insert implementation of Metadata
func (r *Metadata) Insert(tx *sql.Tx) error {
	var err error
	if r.META_ID == 0 {
		metaID, err := getNextId(tx, "metadata", "meta_id")
		if err != nil {
			log.Println("unable to get metaID", err)
			return Error(err, ParametersErrorCode, "", "dbs.metadata.Insert")
		}
		r.META_ID = metaID
	}
	// set defaults and validate the record
	r.SetDefaults()
	err = r.Validate()
	if err != nil {
		log.Println("unable to validate record", err)
		return Error(err, ValidateErrorCode, "", "dbs.metadata.Insert")
	}
	// get SQL statement from static area
	stm := getSQL("insert_meta")
	if Verbose > 0 {
		log.Printf("Insert Metadata record %+v", r)
	} else if Verbose > 1 {
		log.Printf("Insert Metadata\n%s\n%+v", stm, r)
	}
	_, err = tx.Exec(
		stm,
		r.META_ID,
		r.DID,
		r.CREATE_AT,
		r.CREATE_BY,
		r.MODIFY_AT,
		r.MODIFY_BY)
	if err != nil {
		if Verbose > 0 {
			log.Println("unable to insert metadata, error", err)
		}
		return Error(err, InsertErrorCode, "", "dbs.metadata.Insert")
	}
	return nil
}

// Validate implementation of Metadata
func (r *Metadata) Validate() error {
	if err := RecordValidator.Struct(*r); err != nil {
		return DecodeValidatorError(r, err)
	}
	if matched := unixTimePattern.MatchString(fmt.Sprintf("%d", r.CREATE_AT)); !matched {
		msg := "invalid pattern for creation date"
		return Error(InvalidParamErr, PatternErrorCode, msg, "dbs.metadata.Validate")
	}
	if matched := unixTimePattern.MatchString(fmt.Sprintf("%d", r.MODIFY_AT)); !matched {
		msg := "invalid pattern for last modification date"
		return Error(InvalidParamErr, PatternErrorCode, msg, "dbs.metadata.Validate")
	}
	return nil
}

// SetDefaults implements set defaults for Metadata
func (r *Metadata) SetDefaults() {
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

// Decode implementation for Metadata
func (r *Metadata) Decode(reader io.Reader) error {
	// init record with given data record
	data, err := io.ReadAll(reader)
	if err != nil {
		log.Println("fail to read data", err)
		return Error(err, ReaderErrorCode, "", "dbs.metadata.Decode")
	}
	err = json.Unmarshal(data, &r)

	//     decoder := json.NewDecoder(r)
	//     err := decoder.Decode(&rec)
	if err != nil {
		log.Println("fail to decode data", err)
		return Error(err, UnmarshalErrorCode, "", "dbs.metadata.Decode")
	}
	return nil
}
