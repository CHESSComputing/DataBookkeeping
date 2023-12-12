package dbs

// nolint: gocyclo

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/CHESSComputing/DataBookkeeping/utils"
)

// Files represents Files DBS DB table
type Files struct {
	FILE_ID                int64  `json:"file_id"`
	LOGICAL_FILE_NAME      string `json:"logical_file_name" validate:"required"`
	IS_FILE_VALID          int64  `json:"is_file_valid" validate:"number"`
	DATASET_ID             int64  `json:"dataset_id" validate:"number,gt=0"`
	META_ID                string `json:"meta_id" validate:"required"`
	CREATION_DATE          int64  `json:"creation_date" validate:"required,number,gt=0"`
	CREATE_BY              string `json:"create_by" validate:"required"`
	LAST_MODIFICATION_DATE int64  `json:"last_modification_date" validate:"required,number,gt=0"`
	LAST_MODIFIED_BY       string `json:"last_modified_by" validate:"required"`
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
	if val, ok := a.Params["logical_file_name"]; ok {
		if val != "" {
			conds, args = AddParam("logical_file_name", "F.FILE", a.Params, conds, args)
		}
	}
	if val, ok := a.Params["dataset"]; ok {
		if val != "" {
			conds, args = AddParam("dataset", "D.DATASET", a.Params, conds, args)
		}
	}
	if utils.VERBOSE > 0 {
		log.Println("### /file params", a.Params, conds, args)
	}

	tmpl := make(Record)
	tmpl["Owner"] = DBOWNER
	stm, err := LoadTemplateSQL("select_file", tmpl)
	if err != nil {
		return Error(err, LoadErrorCode, "", "dbs.files.Files")
	}

	stm = WhereClause(stm, conds)

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

// Insert implementation of Files
func (r *Files) Insert(tx *sql.Tx) error {
	var err error
	if r.FILE_ID == 0 {
		fileID, err := getNextId(tx, "FILES", "FILE_ID")
		if err != nil {
			log.Println("unable to get fileID", err)
			return Error(err, ParametersErrorCode, "", "dbs.files.Insert")
		}
		r.FILE_ID = fileID
	}
	// set defaults and validate the record
	r.SetDefaults()
	err = r.Validate()
	if err != nil {
		log.Println("unable to validate record", err)
		return Error(err, ValidateErrorCode, "", "dbs.files.Insert")
	}
	// get SQL statement from static area
	stm := getSQL("insert_file")
	if utils.VERBOSE > 0 {
		log.Printf("Insert Files file_id=%d lfn=%s", r.FILE_ID, r.LOGICAL_FILE_NAME)
	} else if utils.VERBOSE > 1 {
		log.Printf("Insert Files\n%s\n%+v", stm, r)
	}
	_, err = tx.Exec(
		stm,
		r.FILE_ID,
		r.LOGICAL_FILE_NAME,
		r.IS_FILE_VALID,
		r.DATASET_ID,
		r.META_ID,
		r.CREATION_DATE,
		r.CREATE_BY,
		r.LAST_MODIFICATION_DATE,
		r.LAST_MODIFIED_BY)
	if err != nil {
		if utils.VERBOSE > 0 {
			log.Println("unable to insert files, error", err)
		}
		return Error(err, InsertErrorCode, "", "dbs.files.Insert")
	}
	return nil
}

// Validate implementation of Files
func (r *Files) Validate() error {
	if err := RecordValidator.Struct(*r); err != nil {
		return DecodeValidatorError(r, err)
	}
	if err := CheckPattern("logical_file_name", r.LOGICAL_FILE_NAME); err != nil {
		return Error(err, PatternErrorCode, "", "dbs.files.Validate")
	}
	if matched := unixTimePattern.MatchString(fmt.Sprintf("%d", r.CREATION_DATE)); !matched {
		msg := "invalid pattern for creation date"
		return Error(InvalidParamErr, PatternErrorCode, msg, "dbs.files.Validate")
	}
	if matched := unixTimePattern.MatchString(fmt.Sprintf("%d", r.LAST_MODIFICATION_DATE)); !matched {
		msg := "invalid pattern for last modification date"
		return Error(InvalidParamErr, PatternErrorCode, msg, "dbs.files.Validate")
	}
	return nil
}

// SetDefaults implements set defaults for Files
func (r *Files) SetDefaults() {
	if r.CREATION_DATE == 0 {
		r.CREATION_DATE = Date()
	}
	if r.LAST_MODIFICATION_DATE == 0 {
		r.LAST_MODIFICATION_DATE = Date()
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
