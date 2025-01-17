package dbs

// DBS osinfo module
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

// OsInfo represents OsInfo DBS DB table
type OsInfo struct {
	OS_ID     int64  `json:"os_id"`
	NAME      string `json:"name" validate:"required"`
	VERSION   string `json:"version" validate:"required"`
	KERNEL    string `json:"kernel" validate:"required"`
	CREATE_AT int64  `json:"create_at"`
	CREATE_BY string `json:"create_by"`
	MODIFY_AT int64  `json:"modify_at"`
	MODIFY_BY string `json:"modify_by"`
}

// OsInfo DBS API
//
//gocyclo:ignore
func (a *API) GetOsInfo() error {
	var args []interface{}
	var conds []string
	var err error

	tmpl := make(map[string]any)
	tmpl["Owner"] = DBOWNER
	stm, err := LoadTemplateSQL("select_osinfo", tmpl)
	if err != nil {
		return Error(err, LoadErrorCode, "", "dbs.osinfo.OsInfo")
	}

	stm = WhereClause(stm, conds)

	// use generic query API to fetch the results from DB
	err = executeAll(a.Writer, a.Separator, stm, args...)
	if err != nil {
		return Error(err, QueryErrorCode, "", "dbs.osinfo.OsInfo")
	}
	return nil
}

// InsertOsInfo inserts osinfo record into DB
func (a *API) InsertOsInfo() error {
	// the API provides Reader which will be used by Decode function to load the HTTP payload
	// and cast it to OsInfo data structure
	return insertRecord(&OsInfo{}, a.Reader)
}

// UpdateOsInfo inserts osinfo record in DB
func (a *API) UpdateOsInfo() error {
	return nil
}

// DeleteOsInfo deletes osinfo record in DB
func (a *API) DeleteOsInfo() error {
	return nil
}

// Update implementation of OsInfo
func (r *OsInfo) Update(tx *sql.Tx) error {
	log.Printf("### Update %+v", r)
	return nil
}

// Insert implementation of OsInfo
func (r *OsInfo) Insert(tx *sql.Tx) (int64, error) {
	var err error
	if r.OS_ID == 0 {
		osinfoID, err := getNextId(tx, "osinfo", "os_id")
		if err != nil {
			log.Println("unable to get osinfoID", err)
			return 0, Error(err, OsInfoErrorCode, "", "dbs.osinfo.Insert")
		}
		r.OS_ID = osinfoID
	}
	// set defaults and validate the record
	r.SetDefaults()
	err = r.Validate()
	if err != nil {
		log.Println("unable to validate record", err)
		return 0, Error(err, ValidateErrorCode, "", "dbs.osinfo.Insert")
	}
	// get SQL statement from static area
	stm := getSQL("insert_osinfo")
	if Verbose > 0 {
		log.Printf("Insert OsInfo record %+v", r)
	} else if Verbose > 1 {
		log.Printf("Insert OsInfo\n%s\n%+v", stm, r)
	}
	_, err = tx.Exec(
		stm,
		r.OS_ID,
		r.NAME,
		r.VERSION,
		r.KERNEL,
		r.CREATE_AT,
		r.CREATE_BY,
		r.MODIFY_AT,
		r.MODIFY_BY)
	if err != nil {
		if Verbose > 0 {
			log.Println("unable to insert osinfo, error", err)
		}
		return 0, Error(err, InsertErrorCode, "", "dbs.osinfo.Insert")
	}
	return r.OS_ID, nil
}

// Validate implementation of OsInfo
func (r *OsInfo) Validate() error {
	if err := RecordValidator.Struct(*r); err != nil {
		return DecodeValidatorError(r, err)
	}
	if matched := lexicon.UnixTimePattern.MatchString(fmt.Sprintf("%d", r.CREATE_AT)); !matched {
		msg := "invalid pattern for creation date"
		return Error(InvalidParamErr, ValidateErrorCode, msg, "dbs.osinfo.Validate")
	}
	if matched := lexicon.UnixTimePattern.MatchString(fmt.Sprintf("%d", r.MODIFY_AT)); !matched {
		msg := "invalid pattern for last modification date"
		return Error(InvalidParamErr, ValidateErrorCode, msg, "dbs.osinfo.Validate")
	}
	return nil
}

// SetDefaults implements set defaults for OsInfo
func (r *OsInfo) SetDefaults() {
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

// Decode implementation for OsInfo
func (r *OsInfo) Decode(reader io.Reader) error {
	// init record with given data record
	data, err := io.ReadAll(reader)
	if err != nil {
		log.Println("fail to read data", err)
		return Error(err, ReaderErrorCode, "", "dbs.osinfo.Decode")
	}
	err = json.Unmarshal(data, &r)

	//     decoder := json.NewDecoder(r)
	//     err := decoder.Decode(&rec)
	if err != nil {
		log.Println("fail to decode data", err)
		return Error(err, UnmarshalErrorCode, "", "dbs.osinfo.Decode")
	}
	return nil
}
