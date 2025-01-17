package dbs

// DBS scripts module
// Copyright (c) 2023 - Valentin Kuznetsov <vkuznet@gmail.com>
//
// nolint: gocyclo

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"

	lexicon "github.com/CHESSComputing/golib/lexicon"
)

// Scripts represents Scripts DBS DB table
type Scripts struct {
	SCRIPT_ID        int64  `json:"script_id"`
	NAME             string `json:"name" validate:"required"`
	OPTIONS          string `json:"options"`
	PARENT_SCRIPT_ID int64  `json:"parent_script_id"`
	CREATE_AT        int64  `json:"create_at"`
	CREATE_BY        string `json:"create_by"`
	MODIFY_AT        int64  `json:"modify_at"`
	MODIFY_BY        string `json:"modify_by"`
}

// Scripts DBS API
//
//gocyclo:ignore
func (a *API) GetScript() error {
	var args []interface{}
	var conds []string
	var err error

	tmpl := make(map[string]any)
	tmpl["Owner"] = DBOWNER
	stm, err := LoadTemplateSQL("select_script", tmpl)
	if err != nil {
		return Error(err, LoadErrorCode, "", "dbs.scripts.Scripts")
	}

	stm = WhereClause(stm, conds)

	// use generic query API to fetch the results from DB
	err = executeAll(a.Writer, a.Separator, stm, args...)
	if err != nil {
		return Error(err, QueryErrorCode, "", "dbs.scripts.Scripts")
	}
	return nil
}

// InsertScript inserts script record into DB
func (a *API) InsertScript() error {
	// the API provides Reader which will be used by Decode function to load the HTTP payload
	// and cast it to Scripts data structure

	// decode input as ScriptRecord
	var rec ScriptRecord

	data, err := io.ReadAll(a.Reader)
	if err != nil {
		log.Println("fail to read data", err)
		return Error(err, ReaderErrorCode, "", "dbs.scripts.InsertScript")
	}
	err = json.Unmarshal(data, &rec)
	if err != nil {
		msg := fmt.Sprintf("fail to decode record")
		log.Println(msg)
		return Error(err, DecodeErrorCode, msg, "dbs.scripts.InsertScript")
	}

	err = rec.Validate()
	if err != nil {
		return Error(err, ValidateErrorCode, "validation error", "dbs.scripts.InsertScript")
	}

	// start transaction
	tx, err := DB.Begin()
	if err != nil {
		return Error(err, TransactionErrorCode, "", "dbs.InsertScript")
	}
	defer tx.Rollback()
	if _, err := rec.Insert(tx); err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return Error(err, CommitErrorCode, "", "dbs.InsertScript")
	}
	return err
}

// UpdateScript inserts script record in DB
func (a *API) UpdateScript() error {
	return nil
}

// DeleteScript deletes script record in DB
func (a *API) DeleteScript() error {
	return nil
}

// Update implementation of Scripts
func (r *Scripts) Update(tx *sql.Tx) error {
	log.Printf("### Update %+v", r)
	return nil
}

// Insert implementation of Scripts
func (r *Scripts) Insert(tx *sql.Tx) (int64, error) {
	var err error
	if r.SCRIPT_ID == 0 {
		scriptID, err := getNextId(tx, "scripts", "script_id")
		if err != nil {
			log.Println("unable to get scriptID", err)
			return 0, Error(err, ScriptsErrorCode, "", "dbs.scripts.Insert")
		}
		r.SCRIPT_ID = scriptID
	}
	// set defaults and validate the record
	r.SetDefaults()
	err = r.Validate()
	if err != nil {
		log.Println("unable to validate record", err)
		return 0, Error(err, ValidateErrorCode, "", "dbs.scripts.Insert")
	}
	// get SQL statement from static area
	stm := getSQL("insert_script")
	if Verbose > 0 {
		log.Printf("Insert Scripts record %+v", r)
	} else if Verbose > 1 {
		log.Printf("Insert Scripts\n%s\n%+v", stm, r)
	}
	_, err = tx.Exec(
		stm,
		r.SCRIPT_ID,
		r.NAME,
		r.OPTIONS,
		r.PARENT_SCRIPT_ID,
		r.CREATE_AT,
		r.CREATE_BY,
		r.MODIFY_AT,
		r.MODIFY_BY)
	if err != nil {
		if Verbose > 0 {
			log.Println("unable to insert scripts, error", err)
		}
		return 0, Error(err, InsertErrorCode, "", "dbs.scripts.Insert")
	}
	return r.SCRIPT_ID, nil
}

// Validate implementation of Scripts
func (r *Scripts) Validate() error {
	if err := RecordValidator.Struct(*r); err != nil {
		return DecodeValidatorError(r, err)
	}
	if matched := lexicon.UnixTimePattern.MatchString(fmt.Sprintf("%d", r.CREATE_AT)); !matched {
		msg := "invalid pattern for creation date"
		return Error(InvalidParamErr, ValidateErrorCode, msg, "dbs.scripts.Validate")
	}
	if matched := lexicon.UnixTimePattern.MatchString(fmt.Sprintf("%d", r.MODIFY_AT)); !matched {
		msg := "invalid pattern for last modification date"
		return Error(InvalidParamErr, ValidateErrorCode, msg, "dbs.scripts.Validate")
	}
	return nil
}

// SetDefaults implements set defaults for Scripts
func (r *Scripts) SetDefaults() {
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

// Decode implementation for Scripts
func (r *Scripts) Decode(reader io.Reader) error {
	// init record with given data record
	data, err := io.ReadAll(reader)
	if err != nil {
		log.Println("fail to read data", err)
		return Error(err, ReaderErrorCode, "", "dbs.scripts.Decode")
	}
	err = json.Unmarshal(data, &r)

	//     decoder := json.NewDecoder(r)
	//     err := decoder.Decode(&rec)
	if err != nil {
		log.Println("fail to decode data", err)
		return Error(err, UnmarshalErrorCode, "", "dbs.scripts.Decode")
	}
	return nil
}
