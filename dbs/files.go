package dbs

// DBS files module
// Copyright (c) 2023 - Valentin Kuznetsov <vkuznet@gmail.com>
//
// nolint: gocyclo

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	lexicon "github.com/CHESSComputing/golib/lexicon"
)

// Files represents Files DBS DB table
type Files struct {
	FILE_ID       int64  `json:"file_id"`
	FILE          string `json:"file" validate:"required"`
	IS_FILE_VALID int64  `json:"is_file_valid" validate:"number"`
	CREATE_AT     int64  `json:"create_at" validate:"required,number,gt=0"`
	CREATE_BY     string `json:"create_by" validate:"required"`
	MODIFY_AT     int64  `json:"modify_at" validate:"required,number,gt=0"`
	MODIFY_BY     string `json:"modify_by" validate:"required"`
}

// FileRecord represents input file record from HTTP request
type FileRecord struct {
	File string `json:"file" validate:"required"`
	Did  string `json:"did" validate:"required"`
}

// Files DBS API
//
//gocyclo:ignore
func (a *API) GetFile() error {
	var args []interface{}
	var conds []string
	var err error

	if len(a.Params) == 0 {
		msg := "Files API with empty parameter map"
		return Error(InvalidParamErr, ParametersErrorCode, msg, "dbs.files.Files")
	}
	if val, ok := a.Params["file"]; ok {
		if val != "" {
			conds, args = AddParam("file", "F.file", a.Params, conds, args)
		}
	}
	if val, ok := a.Params["did"]; ok {
		if val != "" {
			conds, args = AddParam("did", "D.did", a.Params, conds, args)
		}
	}
	if Verbose > 0 {
		log.Println("### /file params", a.Params, conds, args)
	}

	tmpl := make(map[string]any)
	tmpl["Owner"] = DBOWNER
	stm, err := LoadTemplateSQL("select_file", tmpl)
	if err != nil {
		return Error(err, LoadErrorCode, "", "dbs.files.Files")
	}

	stm = WhereClause(stm, conds)
	log.Println("QUERY:\n", stm, conds)

	// use generic query API to fetch the results from DB
	err = executeAll(a.Writer, a.Separator, stm, args...)
	if err != nil {
		return Error(err, QueryErrorCode, "", "dbs.files.Files")
	}
	return nil
}

func (a *API) InsertFile() error {
	// the API provides Reader which will be used by Decode function to load the HTTP payload
	// and cast it to Files data structure
	return insertRecord(&Files{}, a.Reader)
}
func (a *API) UpdateFile() error {
	return nil
}
func (a *API) DeleteFile() error {
	return nil
}

// Update implementation of Files
func (r *Files) Update(tx *sql.Tx) error {
	log.Printf("### Update %+v", r)
	return nil
}

// Insert implementation of Files
func (r *Files) Insert(tx *sql.Tx) (int64, error) {
	var err error
	if r.FILE_ID == 0 {
		fileID, err := getNextId(tx, "files", "file_id")
		if err != nil {
			log.Println("unable to get fileID", err)
			return 0, Error(err, FilesErrorCode, "", "dbs.files.Insert")
		}
		r.FILE_ID = fileID
	}
	// set defaults and validate the record
	r.SetDefaults()
	err = r.Validate()
	if err != nil {
		log.Println("unable to validate record", err)
		return 0, Error(err, ValidateErrorCode, "", "dbs.files.Insert")
	}
	// get SQL statement from static area
	stm := getSQL("insert_file")
	if Verbose > 0 {
		log.Printf("Insert Files file_id=%d lfn=%s", r.FILE_ID, r.FILE)
	} else if Verbose > 1 {
		log.Printf("Insert Files\n%s\n%+v", stm, r)
	}
	_, err = tx.Exec(
		stm,
		r.FILE_ID,
		r.FILE,
		r.IS_FILE_VALID,
		r.CREATE_AT,
		r.CREATE_BY,
		r.MODIFY_AT,
		r.MODIFY_BY)
	if err != nil {
		if Verbose > 0 {
			log.Println("unable to insert files, error", err)
		}
		return 0, Error(err, InsertErrorCode, "", "dbs.files.Insert")
	}
	return r.FILE_ID, nil
}

// Validate implementation of Files
func (r *Files) Validate() error {
	if err := RecordValidator.Struct(*r); err != nil {
		return DecodeValidatorError(r, err)
	}
	if err := lexicon.CheckPattern("file", r.FILE); err != nil {
		return Error(err, ValidateErrorCode, "", "dbs.files.Validate")
	}
	if matched := lexicon.UnixTimePattern.MatchString(fmt.Sprintf("%d", r.CREATE_AT)); !matched {
		msg := "invalid pattern for creation date"
		return Error(InvalidParamErr, ValidateErrorCode, msg, "dbs.files.Validate")
	}
	if matched := lexicon.UnixTimePattern.MatchString(fmt.Sprintf("%d", r.MODIFY_AT)); !matched {
		msg := "invalid pattern for last modification date"
		return Error(InvalidParamErr, ValidateErrorCode, msg, "dbs.files.Validate")
	}
	return nil
}

// SetDefaults implements set defaults for Files
func (r *Files) SetDefaults() {
	if r.CREATE_AT == 0 {
		r.CREATE_AT = Date()
	}
	if r.MODIFY_AT == 0 {
		r.MODIFY_AT = Date()
	}
}

// Decode implementation for Files
func (r *Files) Decode(reader io.Reader) error {
	// init record with given data record
	data, err := io.ReadAll(reader)
	if err != nil {
		log.Println("fail to read data", err)
		return Error(err, ReaderErrorCode, "", "dbs.files.Decode")
	}
	err = json.Unmarshal(data, &r)

	// check if is_file_valid was present in request, if not set it to 1
	if !strings.Contains(string(data), "is_file_valid") {
		r.IS_FILE_VALID = 1
	}

	//     decoder := json.NewDecoder(r)
	//     err := decoder.Decode(&rec)
	if err != nil {
		log.Println("fail to decode data", err)
		return Error(err, UnmarshalErrorCode, "", "dbs.files.Decode")
	}
	return nil
}
