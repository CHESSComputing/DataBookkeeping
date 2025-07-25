package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/CHESSComputing/DataBookkeeping/dbs"
	srvConfig "github.com/CHESSComputing/golib/config"
	"github.com/CHESSComputing/golib/lexicon"
	server "github.com/CHESSComputing/golib/server"
	"github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"
	_ "github.com/mattn/go-sqlite3"
)

// helper function to initialize DBS for tests
func initDBS(dryRun bool, dburi string) *sql.DB {
	config := os.Getenv("FOXDEN_CONFIG")
	if cobj, err := srvConfig.ParseConfig(config); err == nil {
		srvConfig.Config = &cobj
	}

	// set server log file to stdout
	srvConfig.Config.DataBookkeeping.LogFile = ""
	// set log flags
	log.SetFlags(0)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// current directory is a <pwd>/test
	_, err := os.Getwd()
	if err != nil {
		log.Fatal("unable to get current working dir")
	}
	// load Lexicon patterns
	var lexPatterns map[string]lexicon.LexiconPattern
	if srvConfig.Config.DataBookkeeping.LexiconFile != "" {
		lexPatterns, err = lexicon.LoadPatterns(srvConfig.Config.DataBookkeeping.LexiconFile)
	} else {
		// use for unit test in github
		lexPatterns, err = lexicon.LoadPatterns("data/dbs_lexicon.json")
	}
	if err != nil {
		log.Fatal(err)
	}
	lexicon.LexiconPatterns = lexPatterns

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
	dbs.Verbose = 0
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
		db = initDBS(false, dburi)
	}
	if router == nil {
		routes := []server.Route{
			// GET APIs for integration tests
			server.Route{Method: "GET", Path: "/datasets", Handler: DatasetHandler, Authorized: false},
			server.Route{Method: "GET", Path: "/files", Handler: FileHandler, Authorized: false},
			server.Route{Method: "GET", Path: "/scripts", Handler: ScriptHandler, Authorized: false},
			server.Route{Method: "GET", Path: "/configs", Handler: ConfigHandler, Authorized: false},
			server.Route{Method: "GET", Path: "/packages", Handler: PackageHandler, Authorized: false},
			server.Route{Method: "GET", Path: "/environments", Handler: EnvironmentHandler, Authorized: false},
			server.Route{Method: "GET", Path: "/provenance", Handler: ProvenanceHandler, Authorized: false},

			// POST APIs for integration tests
			server.Route{Method: "POST", Path: "/provenance", Handler: ProvenanceHandler, Authorized: false},
			server.Route{Method: "POST", Path: "/dataset", Handler: DatasetHandler, Authorized: false},
			server.Route{Method: "POST", Path: "/file", Handler: FileHandler, Authorized: false},
			server.Route{Method: "POST", Path: "/script", Handler: ScriptHandler, Authorized: false},
			server.Route{Method: "POST", Path: "/config", Handler: ConfigHandler, Authorized: false},
			server.Route{Method: "POST", Path: "/environment", Handler: EnvironmentHandler, Authorized: false},
		}
		router = server.Router(routes, nil, "static", srvConfig.Config.DataBookkeeping.WebServer)
	}
}

// helper function to print any struct in formatted way
func logStruct(t *testing.T, msg string, data any) {
	body, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		t.Logf("%s\n%+v\n", msg, data)
		return
	}
	t.Logf("%s\n%s\n", msg, string(body))
}

// helper function to create http test response recorder
// for given HTTP Method, endPoint, reader and DBS web handler
func responseRecorder(t *testing.T, v TestCase) *httptest.ResponseRecorder {
	// read data from the inpit
	data, err := json.Marshal(v.Input)
	if err != nil {
		t.Fatal(err.Error())
	}
	reader := bytes.NewReader(data)

	if v.Verbose > 0 {
		t.Logf("submit method=%s endpoint=%s url=%s input=%v output=%v code=%v fail=%v data=%s", v.Method, v.Endpoint, v.Url, v.Input, v.Output, v.Code, v.Fail, string(data))
	}
	// setup HTTP request
	req, err := http.NewRequest(v.Method, v.Url, reader)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Accept", "application/json")
	if v.Method == "POST" {
		req.Header.Set("Content-Type", "application/json")
	}

	// create response recorder
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if v.Verbose > 1 {
		logStruct(t, "HTTP request", req)
		logStruct(t, "HTTP response", rr)
	}
	return rr
}
