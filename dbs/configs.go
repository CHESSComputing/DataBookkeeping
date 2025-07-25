package dbs

// DBS config module
// Copyright (c) 2023 - Valentin Kuznetsov <vkuznet@gmail.com>
//
// nolint: gocyclo

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"reflect"
	"strconv"
)

// Config represents Config DBS DB table
type Config struct {
	CONFIG_ID int64  `json:"config_id"`
	CONTENT   any    `json:"content"`
	CREATE_AT int64  `json:"create_at"`
	CREATE_BY string `json:"create_by"`
	MODIFY_AT int64  `json:"modify_at"`
	MODIFY_BY string `json:"modify_by"`
}

// Config DBS API
//
//gocyclo:ignore
func (a *API) GetConfig() error {
	var args []interface{}
	var conds []string
	var err error

	tmpl := make(map[string]any)
	tmpl["Owner"] = DBOWNER
	stm, err := LoadTemplateSQL("select_config", tmpl)
	if err != nil {
		return Error(err, LoadErrorCode, "", "dbs.config.Config")
	}
	log.Printf("### a.Params=%+v", a.Params)
	if val, ok := a.Params["did"]; ok && val != "" {
		conds, args = AddParam("did", "D.did", a.Params, conds, args)
	}

	stm = WhereClause(stm, conds)
	log.Printf("#### stm=%s conds=%v args=%v", stm, conds, args)

	// use generic query API to fetch the results from DB
	err = executeAll(a.Writer, a.Separator, stm, args...)
	if err != nil {
		return Error(err, QueryErrorCode, "", "dbs.config.Config")
	}
	return nil
}

// InsertConfig inserts config record into DB
func (a *API) InsertConfig() error {
	// the API provides Reader which will be used by Decode function to load the HTTP payload
	// and cast it to Config data structure
	return insertRecord(&Config{}, a.Reader)
}

// UpdateConfig inserts config record in DB
func (a *API) UpdateConfig() error {
	// extract payload from API and initialize config attributes
	data, err := io.ReadAll(a.Reader)
	if err != nil {
		return err
	}
	rec := &Config{}
	return DBOperation("update", rec, data, "dbs.UpdateConfig")
}

// DeleteConfig deletes config record in DB
func (a *API) DeleteConfig() error {
	// extract payload from API and initialize config attributes
	data, err := io.ReadAll(a.Reader)
	if err != nil {
		return err
	}
	rec := &Config{}
	return DBOperation("delete", rec, data, "dbs.DeleteConfig")
}

// Delete implementation of Config
func (r *Config) Delete(tx *sql.Tx) error {
	return nil
}

// Update implementation of Config
func (r *Config) Update(tx *sql.Tx) error {
	return nil
}

// Insert implementation of Config
func (r *Config) Insert(tx *sql.Tx) (int64, error) {
	var err error
	if r.CONFIG_ID == 0 {
		configID, err := getNextId(tx, "configs", "config_id")
		if err != nil {
			log.Println("unable to get configID", err)
			return 0, Error(err, ConfigsErrorCode, "", "dbs.config.Insert")
		}
		r.CONFIG_ID = configID
	}
	// we need to convert content any type to a string for database insertion
	if val, err := marshalForSQL(r.CONTENT); err == nil {
		r.CONTENT = val
	} else {
		r.CONTENT = fmt.Sprintf("%v", r.CONTENT)
	}

	// set defaults and validate the record
	r.SetDefaults()
	err = r.Validate()
	if err != nil {
		log.Println("unable to validate record", err)
		return 0, Error(err, ValidateErrorCode, "", "dbs.config.Insert")
	}
	// get SQL statement from static area
	stm := getSQL("insert_config")
	if Verbose > 0 {
		log.Printf("Insert Config record %+v", r)
	} else if Verbose > 1 {
		log.Printf("Insert Config\n%s\n%+v", stm, r)
	}
	_, err = tx.Exec(
		stm,
		r.CONFIG_ID,
		r.CONTENT,
		r.CREATE_AT,
		r.CREATE_BY,
		r.MODIFY_AT,
		r.MODIFY_BY)
	if err != nil {
		if Verbose > 0 {
			log.Println("unable to insert config, error", err)
		}
		return 0, Error(err, InsertErrorCode, "", "dbs.config.Insert")
	}
	log.Printf("##### insertConfig r=%+v", r)
	return r.CONFIG_ID, nil
}

// Validate implementation of Config
func (r *Config) Validate() error {
	return nil
}

// SetDefaults implements set defaults for Config
func (r *Config) SetDefaults() {
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

// helper function to marshal given value to SQL string
func marshalForSQL(val any) (string, error) {
	if val == nil {
		return "NULL", nil
	}

	// Handle string directly
	switch v := val.(type) {
	case string:
		return strconv.Quote(v), nil // wrap in quotes
	case fmt.Stringer:
		return strconv.Quote(v.String()), nil
	case int, int64, float64, bool:
		return fmt.Sprint(v), nil // raw numbers/booleans
	}

	// For maps, slices, structs → JSON
	rv := reflect.ValueOf(val)
	kind := rv.Kind()
	if kind == reflect.Map || kind == reflect.Slice || kind == reflect.Struct {
		j, err := json.Marshal(val)
		if err != nil {
			return "", err
		}
		return strconv.Quote(string(j)), nil // wrap JSON string in SQL quotes
	}

	// Fallback: convert to string
	return strconv.Quote(fmt.Sprint(val)), nil
}
