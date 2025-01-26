package dbs

// DBS packages module
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

// Packages represents Packages DBS DB table
type Packages struct {
	PACKAGE_ID int64  `json:"package_id"`
	NAME       string `json:"name" validate:"required"`
	VERSION    string `json:"version" validate:"required"`
	CREATE_AT  int64  `json:"create_at"`
	CREATE_BY  string `json:"create_by"`
	MODIFY_AT  int64  `json:"modify_at"`
	MODIFY_BY  string `json:"modify_by"`
}

// Packages DBS API
//
//gocyclo:ignore
func (a *API) GetPackage() error {
	var args []interface{}
	var conds []string
	var err error

	tmpl := make(map[string]any)
	tmpl["Owner"] = DBOWNER
	stm, err := LoadTemplateSQL("select_package", tmpl)
	if err != nil {
		return Error(err, LoadErrorCode, "", "dbs.packages.Packages")
	}

	stm = WhereClause(stm, conds)

	// use generic query API to fetch the results from DB
	err = executeAll(a.Writer, a.Separator, stm, args...)
	if err != nil {
		return Error(err, QueryErrorCode, "", "dbs.packages.Packages")
	}
	return nil
}

// InsertPackage inserts package record into DB
func (a *API) InsertPackage() error {
	// the API provides Reader which will be used by Decode function to load the HTTP payload
	// and cast it to Packages data structure

	// decode input as PackageRecord
	var rec PackageRecord

	data, err := io.ReadAll(a.Reader)
	if err != nil {
		log.Println("fail to read data", err)
		return Error(err, ReaderErrorCode, "", "dbs.packages.InsertPackage")
	}
	err = json.Unmarshal(data, &rec)
	if err != nil {
		msg := fmt.Sprintf("fail to decode record")
		log.Println(msg)
		return Error(err, DecodeErrorCode, msg, "dbs.packages.InsertPackage")
	}

	err = rec.Validate()
	if err != nil {
		return Error(err, ValidateErrorCode, "validation error", "dbs.packages.InsertPackage")
	}

	// start transaction
	tx, err := DB.Begin()
	if err != nil {
		return Error(err, TransactionErrorCode, "", "dbs.InsertPackage")
	}
	defer tx.Rollback()
	if _, err := rec.Insert(tx); err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return Error(err, CommitErrorCode, "", "dbs.InsertPackage")
	}
	return err
}

// UpdatePackage inserts package record in DB
func (a *API) UpdatePackage() error {
	return nil
}

// DeletePackage deletes package record in DB
func (a *API) DeletePackage() error {
	return nil
}

// Delete implementation of Packages
func (r *Packages) Delete(tx *sql.Tx) error {
	return nil
}

// Update implementation of Packages
func (r *Packages) Update(tx *sql.Tx) error {
	return nil
}

// Insert implementation of Packages
func (r *Packages) Insert(tx *sql.Tx) (int64, error) {
	var err error
	if r.PACKAGE_ID == 0 {
		packageID, err := getNextId(tx, "packages", "package_id")
		if err != nil {
			log.Println("unable to get packageID", err)
			return 0, Error(err, ParametersErrorCode, "", "dbs.packages.Insert")
		}
		r.PACKAGE_ID = packageID
	}
	// set defaults and validate the record
	r.SetDefaults()
	err = r.Validate()
	if err != nil {
		log.Println("unable to validate record", err)
		return 0, Error(err, ValidateErrorCode, "", "dbs.packages.Insert")
	}
	// get SQL statement from static area
	stm := getSQL("insert_package")
	if Verbose > 0 {
		log.Printf("Insert Packages record %+v", r)
	} else if Verbose > 1 {
		log.Printf("Insert Packages\n%s\n%+v", stm, r)
	}
	_, err = tx.Exec(
		stm,
		r.PACKAGE_ID,
		r.NAME,
		r.VERSION,
		r.CREATE_AT,
		r.CREATE_BY,
		r.MODIFY_AT,
		r.MODIFY_BY)
	if err != nil {
		if Verbose > 0 {
			log.Println("unable to insert packages, error", err)
		}
		return 0, Error(err, InsertErrorCode, "", "dbs.packages.Insert")
	}
	return r.PACKAGE_ID, nil
}

// Validate implementation of Packages
func (r *Packages) Validate() error {
	if err := RecordValidator.Struct(*r); err != nil {
		return DecodeValidatorError(r, err)
	}
	if matched := lexicon.UnixTimePattern.MatchString(fmt.Sprintf("%d", r.CREATE_AT)); !matched {
		msg := "invalid pattern for creation date"
		return Error(InvalidParamErr, PatternErrorCode, msg, "dbs.packages.Validate")
	}
	if matched := lexicon.UnixTimePattern.MatchString(fmt.Sprintf("%d", r.MODIFY_AT)); !matched {
		msg := "invalid pattern for last modification date"
		return Error(InvalidParamErr, PatternErrorCode, msg, "dbs.packages.Validate")
	}
	return nil
}

// SetDefaults implements set defaults for Packages
func (r *Packages) SetDefaults() {
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
