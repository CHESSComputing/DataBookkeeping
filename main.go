package main

// Go implementation of DataBookkeeping (Provenance) service
//
// Copyright (c) 2023 - Valentin Kuznetsov <vkuznet@gmail.com>
//

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"time"

	_ "expvar"         // to be used for monitoring, see https://github.com/divan/expvarmon
	_ "net/http/pprof" // profiler, see https://golang.org/pkg/net/http/pprof/

	oreConfig "github.com/CHESSComputing/common/config"
)

// srv configuration
var _srvConfig *oreConfig.SrvConfig

func info() string {
	goVersion := runtime.Version()
	tstamp := time.Now()
	return fmt.Sprintf("git={{VERSION}} go=%s date=%s", goVersion, tstamp)
}

func main() {
	var version bool
	flag.BoolVar(&version, "version", false, "Show version")
	var config string
	flag.StringVar(&config, "config", "", "server config JSON file")
	flag.Parse()
	if version {
		fmt.Println("server version:", info())
		return
	}
	oConfig, err := oreConfig.ParseConfig(config)
	if err != nil {
		log.Fatal("ERROR", err)
	}
	_srvConfig = &oConfig
	Server()
}
