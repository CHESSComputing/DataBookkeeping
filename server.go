package main

// DBS server
// Copyright (c) 2023 - Valentin Kuznetsov <vkuznet@gmail.com>
//
//
// Some links:  http://www.alexedwards.net/blog/golang-response-snippets
//              http://blog.golang.org/json-and-go
// Go patterns: http://www.golangpatterns.info/home
// Templates:   http://gohugo.io/templates/go-templates/
//              http://golang.org/pkg/html/template/
// Go examples: https://gobyexample.com/
// for Go database API: http://go-database-sql.org/overview.html
// Oracle drivers:
//   _ "gopkg.in/rana/ora.v4"
//   _ "github.com/mattn/go-oci8"
// MySQL driver:
//   _ "github.com/go-sql-driver/mysql"
// SQLite driver:
//  _ "github.com/mattn/go-sqlite3"
//
// Get profile's output
// visit http://localhost:<port>/debug/pprof
// or generate png plots
// go tool pprof -png http://localhost:<port>/debug/pprof/heap > /tmp/heap.png
// go tool pprof -png http://localhost:<port>/debug/pprof/profile > /tmp/profile.png

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/CHESSComputing/DataBookkeeping/dbs"
	"github.com/CHESSComputing/DataBookkeeping/utils"
	srvConfig "github.com/CHESSComputing/golib/config"
	server "github.com/CHESSComputing/golib/server"
	"github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"

	// GO profiler
	_ "net/http/pprof"

	// imports for supported DB drivers
	// go-sqlite driver
	_ "github.com/mattn/go-sqlite3"
)

// helper function to setup our router
func setupRouter() *gin.Engine {
	routes := []server.Route{
		// routes without authorization
		server.Route{Method: "GET", Path: "/datasets", Handler: DatasetHandler, Authorized: false},
		server.Route{Method: "GET", Path: "/dataset", Handler: DatasetHandler, Authorized: false},
		server.Route{Method: "GET", Path: "/dataset/*name", Handler: DatasetHandler, Authorized: false},
		server.Route{Method: "GET", Path: "/files", Handler: FileHandler, Authorized: false},
		server.Route{Method: "GET", Path: "/file", Handler: FileHandler, Authorized: false},
		server.Route{Method: "GET", Path: "/file/*name", Handler: FileHandler, Authorized: false},
		// authorized routes
		server.Route{Method: "POST", Path: "/dataset", Handler: DatasetHandler, Authorized: true, Scope: "write"},
		server.Route{Method: "POST", Path: "/file", Handler: FileHandler, Authorized: true, Scope: "write"},
		server.Route{Method: "PUT", Path: "/dataset/*name", Handler: DatasetHandler, Authorized: true, Scope: "write"},
		server.Route{Method: "PUT", Path: "/file/*name", Handler: FileHandler, Authorized: true, Scope: "write"},
		server.Route{Method: "DELETE", Path: "/dataset/*name", Handler: DatasetHandler, Authorized: true, Scope: "delete"},
		server.Route{Method: "DELETE", Path: "/file/*name", Handler: FileHandler, Authorized: true, Scope: "delete"},
	}
	r := server.Router(routes, nil, "static", srvConfig.Config.DataBookkeeping.WebServer)
	return r
}

// helper function to initialize DB access
func dbInit(dbtype, dburi string) (*sql.DB, error) {
	db, dberr := sql.Open(dbtype, dburi)
	if dberr != nil {
		log.Printf("unable to open %s, error %v", dbtype, dburi)
		return nil, dberr
	}
	dberr = db.Ping()
	if dberr != nil {
		log.Println("DB ping error", dberr)
		return nil, dberr
	}
	db.SetMaxOpenConns(srvConfig.Config.DataBookkeeping.MaxDBConnections)
	db.SetMaxIdleConns(srvConfig.Config.DataBookkeeping.MaxIdleConnections)
	// Disables connection pool for sqlite3. This enables some concurrency with sqlite3 databases
	// See https://stackoverflow.com/questions/57683132/turning-off-connection-pool-for-go-http-client
	// and https://sqlite.org/wal.html
	// This only will apply to sqlite3 databases
	if dbtype == "sqlite3" {
		db.Exec("PRAGMA journal_mode=WAL;")
	}
	return db, nil
}

// Server defines our HTTP server
func Server() {
	// be verbose
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// initialize record validator
	dbs.RecordValidator = validator.New()

	// set database connection once
	log.Println("parse Config.DBFile:", srvConfig.Config.DataBookkeeping.DBFile)
	dbtype, dburi, dbowner := dbs.ParseDBFile(srvConfig.Config.DataBookkeeping.DBFile)
	// for oci driver we know it is oracle backend
	if strings.HasPrefix(dbtype, "oci") {
		utils.ORACLE = true
	}
	log.Println("DBOWNER", dbowner)
	// set static dir
	utils.STATICDIR = srvConfig.Config.DataBookkeeping.WebServer.StaticDir
	utils.VERBOSE = srvConfig.Config.DataBookkeeping.WebServer.Verbose

	// setup DBS
	db, dberr := dbInit(dbtype, dburi)
	if dberr != nil {
		log.Fatal(dberr)
	}
	dbs.DB = db
	dbs.DBTYPE = dbtype
	dbsql := dbs.LoadSQL(dbowner)
	dbs.DBSQL = dbsql
	dbs.DBOWNER = dbowner
	defer dbs.DB.Close()

	r := setupRouter()
	sport := fmt.Sprintf(":%d", srvConfig.Config.DataBookkeeping.WebServer.Port)
	log.Printf("Start HTTP server %s", sport)
	r.Run(sport)
}
