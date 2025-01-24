package dbs

// DBS sites module
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

// Sites represents Sites DBS DB table
type Sites struct {
	SITE_ID   int64  `json:"site_id"`
	SITE      string `json:"site" validate:"required"`
	CREATE_AT int64  `json:"create_at"`
	CREATE_BY string `json:"create_by"`
	MODIFY_AT int64  `json:"modify_at"`
	MODIFY_BY string `json:"modify_by"`
}

// Sites DBS API
//
//gocyclo:ignore
func (a *API) GetSite() error {
	var args []interface{}
	var conds []string
	var err error

	tmpl := make(map[string]any)
	tmpl["Owner"] = DBOWNER
	stm, err := LoadTemplateSQL("select_site", tmpl)
	if err != nil {
		return Error(err, LoadErrorCode, "", "dbs.sites.Sites")
	}

	stm = WhereClause(stm, conds)

	// use generic query API to fetch the results from DB
	err = executeAll(a.Writer, a.Separator, stm, args...)
	if err != nil {
		return Error(err, QueryErrorCode, "", "dbs.sites.Sites")
	}
	return nil
}

// InsertSite inserts site record into DB
func (a *API) InsertSite() error {
	// the API provides Reader which will be used by Decode function to load the HTTP payload
	// and cast it to Sites data structure
	return insertRecord(&Sites{}, a.Reader)
}

// UpdateSite inserts site record in DB
func (a *API) UpdateSite() error {
	return nil
}

// DeleteSite deletes site record in DB
func (a *API) DeleteSite() error {
	return nil
}

// Update implementation of Sites
func (r *Sites) Update(tx *sql.Tx) error {
	return nil
}

// Insert implementation of Sites
func (r *Sites) Insert(tx *sql.Tx) (int64, error) {
	var err error
	if r.SITE_ID == 0 {
		siteID, err := getNextId(tx, "sites", "site_id")
		if err != nil {
			log.Println("unable to get siteID", err)
			return 0, Error(err, SitesErrorCode, "", "dbs.sites.Insert")
		}
		r.SITE_ID = siteID
	}
	// set defaults and validate the record
	r.SetDefaults()
	err = r.Validate()
	if err != nil {
		log.Println("unable to validate record", err)
		return 0, Error(err, ValidateErrorCode, "", "dbs.sites.Insert")
	}
	// get SQL statement from static area
	stm := getSQL("insert_site")
	if Verbose > 0 {
		log.Printf("Insert Sites record %+v", r)
	} else if Verbose > 1 {
		log.Printf("Insert Sites\n%s\n%+v", stm, r)
	}
	_, err = tx.Exec(
		stm,
		r.SITE_ID,
		r.SITE,
		r.CREATE_AT,
		r.CREATE_BY,
		r.MODIFY_AT,
		r.MODIFY_BY)
	if err != nil {
		if Verbose > 0 {
			log.Println("unable to insert sites, error", err)
		}
		return 0, Error(err, InsertErrorCode, "", "dbs.sites.Insert")
	}
	return r.SITE_ID, nil
}

// Validate implementation of Sites
func (r *Sites) Validate() error {
	if err := RecordValidator.Struct(*r); err != nil {
		return DecodeValidatorError(r, err)
	}
	if matched := lexicon.UnixTimePattern.MatchString(fmt.Sprintf("%d", r.CREATE_AT)); !matched {
		msg := "invalid pattern for creation date"
		return Error(InvalidParamErr, ValidateErrorCode, msg, "dbs.sites.Validate")
	}
	if matched := lexicon.UnixTimePattern.MatchString(fmt.Sprintf("%d", r.MODIFY_AT)); !matched {
		msg := "invalid pattern for last modification date"
		return Error(InvalidParamErr, ValidateErrorCode, msg, "dbs.sites.Validate")
	}
	return nil
}

// SetDefaults implements set defaults for Sites
func (r *Sites) SetDefaults() {
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

// Decode implementation for Sites
func (r *Sites) Decode(reader io.Reader) error {
	// init record with given data record
	data, err := io.ReadAll(reader)
	if err != nil {
		log.Println("fail to read data", err)
		return Error(err, ReaderErrorCode, "", "dbs.sites.Decode")
	}
	err = json.Unmarshal(data, &r)

	//     decoder := json.NewDecoder(r)
	//     err := decoder.Decode(&rec)
	if err != nil {
		log.Println("fail to decode data", err)
		return Error(err, UnmarshalErrorCode, "", "dbs.sites.Decode")
	}
	return nil
}
