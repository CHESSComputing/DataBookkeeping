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
//
// Get profile's output
// visit http://localhost:<port>/debug/pprof
// or generate png plots
// go tool pprof -png http://localhost:<port>/debug/pprof/heap > /tmp/heap.png
// go tool pprof -png http://localhost:<port>/debug/pprof/profile > /tmp/profile.png

import (
	"log"

	"github.com/CHESSComputing/DataBookkeeping/dbs"
	srvConfig "github.com/CHESSComputing/golib/config"
	"github.com/CHESSComputing/golib/lexicon"
	server "github.com/CHESSComputing/golib/server"
	sqldb "github.com/CHESSComputing/golib/sqldb"
	"github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"

	// GO profiler
	_ "net/http/pprof"
)

// Verbose controls verbosity level
var Verbose int

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

		server.Route{Method: "GET", Path: "/parents", Handler: ParentHandler, Authorized: false},

		server.Route{Method: "GET", Path: "/children", Handler: ChildHandler, Authorized: false},
		server.Route{Method: "GET", Path: "/child", Handler: ChildHandler, Authorized: false},

		// authorized routes

		// dataset routes
		server.Route{Method: "POST", Path: "/dataset", Handler: DatasetHandler, Authorized: true, Scope: "write"},
		server.Route{Method: "PUT", Path: "/dataset/*name", Handler: DatasetHandler, Authorized: true, Scope: "write"},
		server.Route{Method: "DELETE", Path: "/dataset/*name", Handler: DatasetHandler, Authorized: true, Scope: "delete"},

		// file routes
		server.Route{Method: "POST", Path: "/file", Handler: FileHandler, Authorized: true, Scope: "write"},
		server.Route{Method: "PUT", Path: "/file/*name", Handler: FileHandler, Authorized: true, Scope: "write"},
		server.Route{Method: "DELETE", Path: "/file/*name", Handler: FileHandler, Authorized: true, Scope: "delete"},

		// parent routes
		server.Route{Method: "POST", Path: "/parent", Handler: ParentHandler, Authorized: true, Scope: "write"},
		server.Route{Method: "PUT", Path: "/parent/*name", Handler: ParentHandler, Authorized: true, Scope: "write"},
		server.Route{Method: "DELETE", Path: "/parent/*name", Handler: ParentHandler, Authorized: true, Scope: "write"},
	}
	r := server.Router(routes, nil, "static", srvConfig.Config.DataBookkeeping.WebServer)
	return r
}

// Server defines our HTTP server
func Server() {
	// be verbose
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// initialize record validator
	dbs.RecordValidator = validator.New()
	Verbose = srvConfig.Config.DataBookkeeping.WebServer.Verbose
	dbs.Verbose = Verbose
	dbs.StaticDir = srvConfig.Config.DataBookkeeping.WebServer.StaticDir

	// set database connection once
	log.Println("parse Config.DBFile:", srvConfig.Config.DataBookkeeping.DBFile)
	dbtype, dburi, dbowner := sqldb.ParseDBFile(srvConfig.Config.DataBookkeeping.DBFile)
	log.Println("InitDB: type=%s owner=%s", dbtype, dbowner)

	// setup DBS
	db, dberr := sqldb.InitDB(dbtype, dburi)
	if dberr != nil {
		log.Fatal(dberr)
	}
	dbs.DB = db
	dbs.DBTYPE = dbtype
	dbsql := dbs.LoadSQL(dbowner)
	dbs.DBSQL = dbsql
	dbs.DBOWNER = dbowner
	defer dbs.DB.Close()

	// load Lexicon patterns
	lexPatterns, err := lexicon.LoadPatterns(srvConfig.Config.DataBookkeeping.LexiconFile)
	if err != nil {
		log.Fatal(err)
	}
	lexicon.LexiconPatterns = lexPatterns

	// setup web router and start the service
	r := setupRouter()
	webServer := srvConfig.Config.DataBookkeeping.WebServer
	server.StartServer(r, webServer)
}
