package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"testing"
)

// basic elements to define a test case
type TestCase struct {
	Description string     `json:"description"` // test case description
	Method      string     `json:"method"`      // http method
	Endpoint    string     `json:"endpoint"`    // url endpoint, optional
	Api         string     `json:"api"`         // server api
	Params      url.Values `json:"params"`      // url parameters, optional
	Input       any        `json:"input"`       // POST and PUT body, optional for GET request
	Output      []any      `json:"output"`      // expected response
	Code        int        `json:"code"`        // expected HTTP response code
}

// run test workflow for a single endpoint
// func runTestWorkflow(t *testing.T, c EndpointTestCase) {
func runTestWorkflow(t *testing.T, v TestCase) {
	initServer()
	t.Logf("Running test case: %+v", v)
	t.Run(v.Description, func(t *testing.T) {

		// create request body
		data, err := json.Marshal(v.Input)
		if err != nil {
			t.Fatal(err.Error())
		}
		reader := bytes.NewReader(data)

		t.Logf("submit method=%s endpoint=%s api=%s data=%s", v.Method, v.Endpoint, v.Api, string(data))
		rr := responseRecorder(t, v.Method, v.Endpoint, v.Api, reader)
		t.Logf("response %+v", rr)
		// check response code
		if rr.Code != v.Code {
			msg := fmt.Sprintf("ERROR: wrong response code, expect=%d received=%d", v.Code, rr.Code)
			t.Fatal(msg)
		}

		// check response
		var d []map[string]any
		if v.Method == "GET" {
			err = json.NewDecoder(rr.Body).Decode(&d)
			if err != nil {
				t.Fatalf("Failed to decode body, %v", err)
			}
		} else if v.Method == "POST" {
			data = rr.Body.Bytes()
			if len(data) != 0 {
				err = json.Unmarshal(data, &d)
				if err != nil {
					msg := fmt.Sprintf("Rsponse data '%s', error %v", string(data), err)
					t.Fatal(msg)
				}
			}
		}

	})
}

// TestIntegration provides integration tests
func TestIntegration(t *testing.T) {
	var testCases []TestCase
	idir := os.Getenv("DBS_INT_TEST_DIR")
	if idir == "" {
		return
	}
	// loop over files in int dir
	files, err := ioutil.ReadDir(idir)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		fname := f.Name()
		if !strings.HasPrefix(fname, "int_") {
			continue
		}
		t.Logf("reding test file %s", fname)
		// load test from test file
		file, err := os.Open(fname)
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()
		data, err := io.ReadAll(file)
		if err != nil {
			t.Fatal(err)
		}
		err = json.Unmarshal(data, &testCases)
		if err != nil {
			t.Fatal(err)
		}
		for _, v := range testCases {
			runTestWorkflow(t, v)
		}
	}
}
