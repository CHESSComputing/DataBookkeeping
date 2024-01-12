package dbs

// DBS buckets module
// Copyright (c) 2023 - Valentin Kuznetsov <vkuznet@gmail.com>
//
// nolint: gocyclo

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/CHESSComputing/DataBookkeeping/utils"
)

// Buckets represents Buckets DBS DB table
type Buckets struct {
	BUCKET_ID  int64  `json:"bucket_id"`
	BUCKET     string `json:"bucket" validate:"required"`
	META_ID    string `json:"meta_id" validate:"required"`
	DATASET_ID int64  `json:"dataset_id" validate:"required"`
	CREATE_AT  int64  `json:"create_at"`
	CREATE_BY  string `json:"create_by"`
	MODIFY_AT  int64  `json:"modify_at"`
	MODIFY_BY  string `json:"modify_by"`
}

// Buckets DBS API
//
//gocyclo:ignore
func (a *API) GetBucket() error {
	var args []interface{}
	var conds []string
	var err error

	tmpl := make(map[string]any)
	tmpl["Owner"] = DBOWNER
	stm, err := LoadTemplateSQL("select_bucket", tmpl)
	if err != nil {
		return Error(err, LoadErrorCode, "", "dbs.buckets.Buckets")
	}

	stm = WhereClause(stm, conds)

	// use generic query API to fetch the results from DB
	err = executeAll(a.Writer, a.Separator, stm, args...)
	if err != nil {
		return Error(err, QueryErrorCode, "", "dbs.buckets.Buckets")
	}
	return nil
}

// InsertBucket inserts bucket record into DB
func (a *API) InsertBucket() error {
	// the API provides Reader which will be used by Decode function to load the HTTP payload
	// and cast it to Buckets data structure
	return insertRecord(&Buckets{}, a.Reader)
}

// UpdateBucket inserts bucket record in DB
func (a *API) UpdateBucket() error {
	return nil
}

// DeleteBucket deletes bucket record in DB
func (a *API) DeleteBucket() error {
	return nil
}

// Insert implementation of Buckets
func (r *Buckets) Insert(tx *sql.Tx) error {
	var err error
	if r.BUCKET_ID == 0 {
		bucketID, err := getNextId(tx, "buckets", "bucket_id")
		if err != nil {
			log.Println("unable to get bucketID", err)
			return Error(err, ParametersErrorCode, "", "dbs.buckets.Insert")
		}
		r.BUCKET_ID = bucketID
	}
	// set defaults and validate the record
	r.SetDefaults()
	err = r.Validate()
	if err != nil {
		log.Println("unable to validate record", err)
		return Error(err, ValidateErrorCode, "", "dbs.buckets.Insert")
	}
	// get SQL statement from static area
	stm := getSQL("insert_bucket")
	if utils.VERBOSE > 0 {
		log.Printf("Insert Buckets record %+v", r)
	} else if utils.VERBOSE > 1 {
		log.Printf("Insert Buckets\n%s\n%+v", stm, r)
	}
	_, err = tx.Exec(
		stm,
		r.BUCKET_ID,
		r.BUCKET,
		r.META_ID,
		r.DATASET_ID,
		r.CREATE_AT,
		r.CREATE_BY,
		r.MODIFY_AT,
		r.MODIFY_BY)
	if err != nil {
		if utils.VERBOSE > 0 {
			log.Println("unable to insert buckets, error", err)
		}
		return Error(err, InsertErrorCode, "", "dbs.buckets.Insert")
	}
	return nil
}

// Validate implementation of Buckets
func (r *Buckets) Validate() error {
	if err := RecordValidator.Struct(*r); err != nil {
		return DecodeValidatorError(r, err)
	}
	if matched := unixTimePattern.MatchString(fmt.Sprintf("%d", r.CREATE_AT)); !matched {
		msg := "invalid pattern for creation date"
		return Error(InvalidParamErr, PatternErrorCode, msg, "dbs.buckets.Validate")
	}
	if matched := unixTimePattern.MatchString(fmt.Sprintf("%d", r.MODIFY_AT)); !matched {
		msg := "invalid pattern for last modification date"
		return Error(InvalidParamErr, PatternErrorCode, msg, "dbs.buckets.Validate")
	}
	return nil
}

// SetDefaults implements set defaults for Buckets
func (r *Buckets) SetDefaults() {
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

// Decode implementation for Buckets
func (r *Buckets) Decode(reader io.Reader) error {
	// init record with given data record
	data, err := io.ReadAll(reader)
	if err != nil {
		log.Println("fail to read data", err)
		return Error(err, ReaderErrorCode, "", "dbs.buckets.Decode")
	}
	err = json.Unmarshal(data, &r)

	//     decoder := json.NewDecoder(r)
	//     err := decoder.Decode(&rec)
	if err != nil {
		log.Println("fail to decode data", err)
		return Error(err, UnmarshalErrorCode, "", "dbs.buckets.Decode")
	}
	return nil
}
