package dbs

// DBS processing module
// Copyright (c) 2023 - Valentin Kuznetsov <vkuznet@gmail.com>
//
// nolint: gocyclo

import (
	"database/sql"
	"fmt"
	"log"

	lexicon "github.com/CHESSComputing/golib/lexicon"
)

// Processing represents Processing DBS DB table
type Processing struct {
	PROCESSING_ID int64  `json:"processing_id"`
	PROCESSING    string `json:"processing" validate:"required"`
	CREATE_AT     int64  `json:"create_at"`
	CREATE_BY     string `json:"create_by"`
	MODIFY_AT     int64  `json:"modify_at"`
	MODIFY_BY     string `json:"modify_by"`
}

// Processing DBS API
//
//gocyclo:ignore
func (a *API) GetProcessing() error {
	var args []interface{}
	var conds []string
	var err error

	tmpl := make(map[string]any)
	tmpl["Owner"] = DBOWNER
	stm, err := LoadTemplateSQL("select_processing", tmpl)
	if err != nil {
		return Error(err, LoadErrorCode, "", "dbs.processing.Processing")
	}

	stm = WhereClause(stm, conds)

	// use generic query API to fetch the results from DB
	err = executeAll(a.Writer, a.Separator, stm, args...)
	if err != nil {
		return Error(err, QueryErrorCode, "", "dbs.processing.Processing")
	}
	return nil
}

// InsertProcessing inserts processing record into DB
func (a *API) InsertProcessing() error {
	// the API provides Reader which will be used by Decode function to load the HTTP payload
	// and cast it to Processing data structure
	return insertRecord(&Processing{}, a.Reader)
}

// UpdateProcessing inserts processing record in DB
func (a *API) UpdateProcessing() error {
	return nil
}

// DeleteProcessing deletes processing record in DB
func (a *API) DeleteProcessing() error {
	return nil
}

// Delete implementation of Processing
func (r *Processing) Delete(tx *sql.Tx) error {
	return nil
}

// Update implementation of Processing
func (r *Processing) Update(tx *sql.Tx) error {
	return nil
}

// Insert implementation of Processing
func (r *Processing) Insert(tx *sql.Tx) (int64, error) {
	var err error
	if r.PROCESSING_ID == 0 {
		processingID, err := getNextId(tx, "processing", "processing_id")
		if err != nil {
			log.Println("unable to get processingID", err)
			return 0, Error(err, ProcessingErrorCode, "", "dbs.processing.Insert")
		}
		r.PROCESSING_ID = processingID
	}
	// set defaults and validate the record
	r.SetDefaults()
	err = r.Validate()
	if err != nil {
		log.Println("unable to validate record", err)
		return 0, Error(err, ValidateErrorCode, "", "dbs.processing.Insert")
	}
	// get SQL statement from static area
	stm := getSQL("insert_processing")
	if Verbose > 0 {
		log.Printf("Insert Processing record %+v", r)
	} else if Verbose > 1 {
		log.Printf("Insert Processing\n%s\n%+v", stm, r)
	}
	_, err = tx.Exec(
		stm,
		r.PROCESSING_ID,
		r.PROCESSING,
		r.CREATE_AT,
		r.CREATE_BY,
		r.MODIFY_AT,
		r.MODIFY_BY)
	if err != nil {
		if Verbose > 0 {
			log.Println("unable to insert processing, error", err)
		}
		return 0, Error(err, InsertErrorCode, "", "dbs.processing.Insert")
	}
	return r.PROCESSING_ID, nil
}

// Validate implementation of Processing
func (r *Processing) Validate() error {
	if err := RecordValidator.Struct(*r); err != nil {
		return DecodeValidatorError(r, err)
	}
	if matched := lexicon.UnixTimePattern.MatchString(fmt.Sprintf("%d", r.CREATE_AT)); !matched {
		msg := "invalid pattern for creation date"
		return Error(InvalidParamErr, ValidateErrorCode, msg, "dbs.processing.Validate")
	}
	if matched := lexicon.UnixTimePattern.MatchString(fmt.Sprintf("%d", r.MODIFY_AT)); !matched {
		msg := "invalid pattern for last modification date"
		return Error(InvalidParamErr, ValidateErrorCode, msg, "dbs.processing.Validate")
	}
	return nil
}

// SetDefaults implements set defaults for Processing
func (r *Processing) SetDefaults() {
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
