package dbs

import (
	"net/http"
	"time"
)

// Timeout represents DBS timeout used by HttpClient
var Timeout int

// HttpClient is HTTP client for urlfetch server
func HttpClient(tout int) *http.Client {
	timeout := time.Duration(tout) * time.Second
	if tout > 0 {
		return &http.Client{Timeout: timeout}
	}
	return &http.Client{}
}
