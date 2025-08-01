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

	lexicon "github.com/CHESSComputing/golib/lexicon"
)

// Parents represents Parents DBS DB table
type Parents struct {
	PARENT_ID  int64  `json:"parent_id"`
	DATASET_ID int64  `json:"dataset_id"`
	CREATE_AT  int64  `json:"create_at"`
	CREATE_BY  string `json:"create_by"`
	MODIFY_AT  int64  `json:"modify_at"`
	MODIFY_BY  string `json:"modify_by"`
}

// ParentRecord represents input parent record from HTTP request
type ParentRecord struct {
	Parent string `json:"parent" validate:"required"`
	Did    string `json:"did" validate:"required"`
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
	//stm, err := LoadTemplateSQL("select_parent_did", tmpl)
	stm, err := LoadTemplateSQL("select_parent", tmpl)
	if err != nil {
		return Error(err, LoadErrorCode, "", "dbs.parents.GetParent")
	}
	if val, ok := a.Params["did"]; ok {
		if val != "" {
			conds, args = AddParam("did", "D.did", a.Params, conds, args)
		}
	} else {
		return Error(err, QueryErrorCode, "no did is provided", "dbs.parents.GetParent")
	}

	stm = WhereClause(stm, conds)

	// use generic query API to fetch the results from DB
	err = executeAll(a.Writer, a.Separator, stm, args...)
	if err != nil {
		return Error(err, QueryErrorCode, "", "dbs.parents.GetParent")
	}
	return nil
}

// InsertParent inserts parent record into DB
func (a *API) InsertParent() error {
	// read given input
	data, err := io.ReadAll(a.Reader)
	if err != nil {
		log.Println("fail to read data", err)
		return Error(err, ReaderErrorCode, "", "dbs.parents.InsertParent")
	}
	rec := ParentRecord{}
	if a.ContentType == "application/json" {
		err = json.Unmarshal(data, &rec)
	} else {
		log.Println("Parser dataset record using default application/json mtime")
		err = json.Unmarshal(data, &rec)
	}
	if err != nil {
		log.Println("reading", a.ContentType)
		log.Println("reading data", string(data))
		log.Println("fail to decode data", err)
		return Error(err, UnmarshalErrorCode, "", "dbs.parents.InsertParent")
	}
	// find parent id and current did
	tx, err := DB.Begin()
	if err != nil {
		return Error(err, TransactionErrorCode, "", "dbs.insertRecord")
	}
	defer tx.Rollback()
	datasetId, err := GetID(tx, "datasets", "dataset_id", "did", rec.Did)
	if err != nil {
		return Error(err, LoadErrorCode, "", "dbs.parents.InsertRecord")
	}
	parentId, err := GetID(tx, "datasets", "dataset_id", "did", rec.Parent)
	if err != nil {
		return Error(err, LoadErrorCode, "", "dbs.parents.InsertRecord")
	}
	// the API provides Reader which will be used by Decode function to load the HTTP payload
	// and cast it to Parents data structure
	record := Parents{
		PARENT_ID:  parentId,
		DATASET_ID: datasetId,
	}
	_, err = record.Insert(tx)
	if err != nil {
		return err
	}
	err = tx.Commit()
	return err
}

// UpdateParent inserts parent record in DB
func (a *API) UpdateParent() error {
	return nil
}

// DeleteParent deletes parent record in DB
func (a *API) DeleteParent() error {
	return nil
}

// Delete implementation of Parents
func (r *Parents) Delete(tx *sql.Tx) error {
	return nil
}

// Update implementation of Parents
func (r *Parents) Update(tx *sql.Tx) error {
	return nil
}

// Insert implementation of Parents
func (r *Parents) Insert(tx *sql.Tx) (int64, error) {
	var err error
	if r.PARENT_ID == 0 {
		parentID, err := getNextId(tx, "parents", "parent_id")
		if err != nil {
			log.Println("unable to get parentID", err)
			return 0, Error(err, ParentsErrorCode, "", "dbs.parents.Insert")
		}
		r.PARENT_ID = parentID
	}
	// set defaults and validate the record
	r.SetDefaults()
	err = r.Validate()
	if err != nil {
		log.Println("unable to validate record", err)
		return 0, Error(err, ValidateErrorCode, "", "dbs.parents.Insert")
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
		r.DATASET_ID,
		r.CREATE_AT,
		r.CREATE_BY,
		r.MODIFY_AT,
		r.MODIFY_BY)
	if err != nil {
		if Verbose > 0 {
			log.Println("unable to insert parents, error", err)
		}
		return 0, Error(err, InsertErrorCode, "", "dbs.parents.Insert")
	}
	return r.PARENT_ID, nil
}

// Validate implementation of Parents
func (r *Parents) Validate() error {
	if err := RecordValidator.Struct(*r); err != nil {
		return DecodeValidatorError(r, err)
	}
	if matched := lexicon.UnixTimePattern.MatchString(fmt.Sprintf("%d", r.CREATE_AT)); !matched {
		msg := "invalid pattern for creation date"
		return Error(InvalidParamErr, ValidateErrorCode, msg, "dbs.parents.Validate")
	}
	if matched := lexicon.UnixTimePattern.MatchString(fmt.Sprintf("%d", r.MODIFY_AT)); !matched {
		msg := "invalid pattern for last modification date"
		return Error(InvalidParamErr, ValidateErrorCode, msg, "dbs.parents.Validate")
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
