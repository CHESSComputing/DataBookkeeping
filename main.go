package main

// Go implementation of DataBookkeeping (Provenance) service
//
// Copyright (c) 2023 - Valentin Kuznetsov <vkuznet@gmail.com>
//

import (
	"fmt"
	"runtime"
	"time"

	_ "expvar"         // to be used for monitoring, see https://github.com/divan/expvarmon
	_ "net/http/pprof" // profiler, see https://golang.org/pkg/net/http/pprof/

	srvConfig "github.com/CHESSComputing/golib/config"
)

func info() string {
	goVersion := runtime.Version()
	tstamp := time.Now()
	return fmt.Sprintf("git={{VERSION}} go=%s date=%s", goVersion, tstamp)
}

func main() {
	srvConfig.Init()
	Server()
}
