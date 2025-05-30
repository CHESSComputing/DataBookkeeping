package dbs

// DBS environments module
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

// Environments represents Environments DBS DB table
type Environments struct {
	ENVIRONMENT_ID        int64  `json:"environment_id"`
	NAME                  string `json:"name" validate:"required"`
	VERSION               string `json:"version" validate:"required"`
	DETAILS               string `json:"details" validate:"required"`
	OS_ID                 int64  `json:"os_id"`
	PARENT_ENVIRONMENT_ID int64  `json:"parent_environment_id"`
	CREATE_AT             int64  `json:"create_at"`
	CREATE_BY             string `json:"create_by"`
	MODIFY_AT             int64  `json:"modify_at"`
	MODIFY_BY             string `json:"modify_by"`
}

// Environments DBS API
//
//gocyclo:ignore
func (a *API) GetEnvironment() error {
	var args []interface{}
	var conds []string
	var err error

	tmpl := make(map[string]any)
	tmpl["Owner"] = DBOWNER
	stm, err := LoadTemplateSQL("select_environment", tmpl)
	if err != nil {
		return Error(err, LoadErrorCode, "", "dbs.environments.Environments")
	}
	if val, ok := a.Params["did"]; ok && val != "" {
		conds, args = AddParam("did", "D.did", a.Params, conds, args)
	}

	stm = WhereClause(stm, conds)

	// use generic query API to fetch the results from DB
	err = executeAll(a.Writer, a.Separator, stm, args...)
	if err != nil {
		return Error(err, QueryErrorCode, "", "dbs.environments.Environments")
	}
	return nil
}

// InsertEnvironment inserts environment record into DB
func (a *API) InsertEnvironment() error {
	// the API provides Reader which will be used by Decode function to load the HTTP payload
	// and cast it to Environments data structure

	// decode input as EnvironmentRecord
	var rec EnvironmentRecord

	data, err := io.ReadAll(a.Reader)
	if err != nil {
		log.Println("fail to read data", err)
		return Error(err, ReaderErrorCode, "", "dbs.environments.InsertEnvironment")
	}
	err = json.Unmarshal(data, &rec)
	if err != nil {
		msg := fmt.Sprintf("fail to decode record")
		log.Println(msg)
		return Error(err, DecodeErrorCode, msg, "dbs.environments.InsertEnvironment")
	}

	err = rec.Validate()
	if err != nil {
		return Error(err, ValidateErrorCode, "validation error", "dbs.environments.InsertEnvironment")
	}

	// start transaction
	tx, err := DB.Begin()
	if err != nil {
		return Error(err, TransactionErrorCode, "", "dbs.InsertEnvironment")
	}
	defer tx.Rollback()
	if _, err := rec.Insert(tx); err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return Error(err, CommitErrorCode, "", "dbs.InsertEnvironment")
	}
	return err
}

// UpdateEnvironment inserts environment record in DB
func (a *API) UpdateEnvironment() error {
	// extract payload from API and initialize environment attributes
	data, err := io.ReadAll(a.Reader)
	if err != nil {
		return err
	}
	rec := &Environments{}
	return DBOperation("update", rec, data, "dbs.UpdateEnvironment")
}

// DeleteEnvironment deletes environment record in DB
func (a *API) DeleteEnvironment() error {
	// extract payload from API and initialize environment attributes
	data, err := io.ReadAll(a.Reader)
	if err != nil {
		return err
	}
	rec := &Environments{}
	return DBOperation("delete", rec, data, "dbs.DeleteEnvironment")
}

// Delete implementation of Environments
func (r *Environments) Delete(tx *sql.Tx) error {
	return nil
}

// Update implementation of Environments
func (r *Environments) Update(tx *sql.Tx) error {
	return nil
}

// Insert implementation of Environments
func (r *Environments) Insert(tx *sql.Tx) (int64, error) {
	var err error
	if r.ENVIRONMENT_ID == 0 {
		environmentID, err := getNextId(tx, "environments", "environment_id")
		if err != nil {
			log.Println("unable to get environmentID", err)
			return 0, Error(err, EnvironmentsErrorCode, "", "dbs.environments.Insert")
		}
		r.ENVIRONMENT_ID = environmentID
	}
	// set defaults and validate the record
	r.SetDefaults()
	err = r.Validate()
	if err != nil {
		log.Println("unable to validate record", err)
		return 0, Error(err, ValidateErrorCode, "", "dbs.environments.Insert")
	}
	// get SQL statement from static area
	stm := getSQL("insert_environment")
	if Verbose > 0 {
		log.Printf("Insert Environments record %+v", r)
	} else if Verbose > 1 {
		log.Printf("Insert Environments\n%s\n%+v", stm, r)
	}
	_, err = tx.Exec(
		stm,
		r.ENVIRONMENT_ID,
		r.NAME,
		r.VERSION,
		r.DETAILS,
		r.OS_ID,
		r.PARENT_ENVIRONMENT_ID,
		r.CREATE_AT,
		r.CREATE_BY,
		r.MODIFY_AT,
		r.MODIFY_BY)
	if err != nil {
		if Verbose > 0 {
			log.Println("unable to insert environments, error", err)
		}
		return 0, Error(err, InsertErrorCode, "", "dbs.environments.Insert")
	}
	return r.ENVIRONMENT_ID, nil
}

// Validate implementation of Environments
func (r *Environments) Validate() error {
	if err := RecordValidator.Struct(*r); err != nil {
		return DecodeValidatorError(r, err)
	}
	if matched := lexicon.UnixTimePattern.MatchString(fmt.Sprintf("%d", r.CREATE_AT)); !matched {
		msg := "invalid pattern for creation date"
		return Error(InvalidParamErr, PatternErrorCode, msg, "dbs.environments.Validate")
	}
	if matched := lexicon.UnixTimePattern.MatchString(fmt.Sprintf("%d", r.MODIFY_AT)); !matched {
		msg := "invalid pattern for last modification date"
		return Error(InvalidParamErr, PatternErrorCode, msg, "dbs.environments.Validate")
	}
	return nil
}

// SetDefaults implements set defaults for Environments
func (r *Environments) SetDefaults() {
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
