package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/CHESSComputing/DataBookkeeping/dbs"
	srvConfig "github.com/CHESSComputing/golib/config"
	server "github.com/CHESSComputing/golib/server"
	"github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"
	_ "github.com/mattn/go-sqlite3"
)

// helper function to initialize DB for tests
func initDB(dryRun bool, dburi string) *sql.DB {
	srvConfig.Init()
	log.SetFlags(0)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// current directory is a <pwd>/test
	_, err := os.Getwd()
	if err != nil {
		log.Fatal("unable to get current working dir")
	}
	dbs.StaticDir = "static"
	dbtype := "sqlite3"
	dbowner := "sqlite"

	db, err := sql.Open(dbtype, dburi)
	if err != nil {
		log.Fatal("unable to open db file", err)
	}
	dbs.DB = db
	dbs.DBTYPE = dbtype
	dbsql := dbs.LoadSQL(dbowner)
	dbs.DBSQL = dbsql
	dbs.DBOWNER = dbowner
	dbs.Verbose = 1
	if dryRun {
		dbs.DRYRUN = true
	}
	// init validator
	dbs.RecordValidator = validator.New()
	dbs.FileLumiChunkSize = 1000

	return db
}

var router *gin.Engine
var db *sql.DB

func initServer() {
	if db == nil {
		// initialize DB for testing
		dburi := os.Getenv("DBS_DB_FILE")
		if dburi == "" {
			log.Fatal("DBS_DB_FILE not defined")
		}
		db = initDB(false, dburi)
	}
	if router == nil {
		routes := []server.Route{
			server.Route{Method: "GET", Path: "/dataset", Handler: DatasetHandler, Authorized: false},
			server.Route{Method: "POST", Path: "/dataset", Handler: DatasetHandler, Authorized: false},
		}
		router = server.Router(routes, nil, "static", srvConfig.Config.DataBookkeeping.WebServer)
	}
}

// helper function to print any struct in formatted way
func logStruct(t *testing.T, msg string, data any) {
	body, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		//         t.Fatal(err)
		t.Logf("%s\n%+v\n", msg, data)
		return
	}
	t.Logf("%s\n%s\n", msg, string(body))
}

// helper function to create http test response recorder
// for given HTTP Method, endPoint, reader and DBS web handler
func responseRecorder(t *testing.T, method, endPoint, api string, reader io.Reader) *httptest.ResponseRecorder {
	// setup HTTP request
	req, err := http.NewRequest(method, api, reader)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Accept", "application/json")
	if method == "POST" {
		req.Header.Set("Content-Type", "application/json")
	}

	// create response recorder
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	logStruct(t, "HTTP request", req)
	logStruct(t, "HTTP response", rr)
	return rr
}

// helper function to create http test response recorder
// for given HTTP Method, endPoint, reader and DBS web handler
func respRecorder(t *testing.T, method, endPoint, api string, reader io.Reader) (*httptest.ResponseRecorder, error) {
	rr := responseRecorder(t, method, endPoint, api, reader)
	if status := rr.Code; status != http.StatusOK {
		data, e := io.ReadAll(rr.Body)
		var emsg string
		if e != nil {
			emsg = fmt.Sprintf("unable to read reasponse body, error: %v", e)
		} else {
			emsg = fmt.Sprintf("handler returned status code: %v message: %s",
				status, string(data))
		}
		msg := fmt.Sprintf("HTTP status %v, error %s", status, emsg)
		return nil, errors.New(msg)
	}
	return rr, nil
}

// TestDBSWriter provides a test to DBS writer functionality
func TestDBSWriter(t *testing.T) {
	initServer()
	var err error

	endPoint := "/dataset"
	t.Log("insert datasets")
	insertData(t, db, "POST", endPoint, "data/datasets.json", "dataset")
	t.Log("re-insert primary dataset")
	insertData(t, db, "POST", endPoint, "data/datasets.json", "dataset")

	t.Logf("finish DBS writer test")
	err = db.Close()
	if err != nil {
		t.Error(err.Error())
	}
}

// insertData provides a test to insert DBS data
func insertData(t *testing.T, db *sql.DB, method, endPoint, dataFile, attr string) {
	// setup HTTP request
	var data []byte
	var err error
	var rr *httptest.ResponseRecorder
	data, err = os.ReadFile(dataFile)
	if err != nil {
		t.Logf("ERROR: unable to read %s error %v", dataFile, err.Error())
		t.Fatal(err.Error())
	}
	var rec map[string]any
	err = json.Unmarshal(data, &rec)
	if err != nil {
		t.Logf("unable to unmarshal received data into map[string]any, error %v, try next []dbs.Record", err)
		// let's try to load list of records
		var rrr []map[string]any
		err = json.Unmarshal(data, &rrr)
		if err != nil {
			t.Fatalf("ERROR: unable to unmarshal received data '%s', error %v", string(data), err)
		}
		t.Log("succeed to load record as []map[string]any")
	}
	reader := bytes.NewReader(data)

	// test writer DBS API
	postApi := endPoint
	rr, err = respRecorder(t, method, endPoint, postApi, reader)
	if err != nil {
		t.Logf("ERROR: unable to make %s HTTP request with endPoint=%s, error %v", method, endPoint, err)
		t.Fatal(err)
	}

	t.Logf("writer endPoint %s send data:\n%v", endPoint, string(data))

	// if no attribute is provided we'll skip GET API test
	if attr == "" {
		t.Log("skip get API call since no attr is provided")
		return
	}
	// test reader GET DBS API
	val, ok := rec[attr]
	if !ok {
		t.Fatalf("ERROR: unable to extract %s from loaded record", attr)
	}
	var value string
	switch v := val.(type) {
	case string:
		value = url.QueryEscape(v)
	default:
		value = fmt.Sprintf("%v", v)
	}
	getApi := fmt.Sprintf("%s?%s=%s", endPoint, attr, value)
	rr, err = respRecorder(t, "GET", endPoint, getApi, reader)
	if err != nil {
		t.Logf("ERROR: unable to place GET HTTP request with api=%s, error %v", getApi, err)
		t.Fatal(err)
	}

	// unmarshal received records
	var records []map[string]any
	data = rr.Body.Bytes()
	err = json.Unmarshal(data, &records)
	if err != nil {
		t.Fatalf("ERROR: unable to unmarshal received data '%s', error %v", string(data), err)
	}
	t.Logf("reader api %s received data:\n%v", getApi, string(data))
}
