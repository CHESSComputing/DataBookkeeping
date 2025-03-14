package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// basic elements to define a test case
type TestCase struct {
	Description  string   `json:"description"`   // test case description
	Method       string   `json:"method"`        // http method
	Endpoint     string   `json:"endpoint"`      // url endpoint, optional
	Url          string   `json:"url"`           // url parameters, optional
	Input        any      `json:"input"`         // POST and PUT body, optional for GET request
	Output       []string `json:"output"`        // expected response patterns
	Code         int      `json:"code"`          // expected HTTP response code
	Verbose      int      `json:"verbose"`       // verbosity level
	Fail         bool     `json:"fail"`          // should test fail
	DumpResponse bool     `json:"dump_response"` // enable dump of the response
}

// run test workflow for a single endpoint
// func runTestWorkflow(t *testing.T, c EndpointTestCase) {
func runTestWorkflow(t *testing.T, v TestCase) {
	initServer()
	if v.Verbose > 0 {
		t.Logf("Running test case: %+v", v)
	}
	var err error
	t.Run(v.Description, func(t *testing.T) {

		rr := responseRecorder(t, v)
		if v.Verbose > 1 {
			t.Logf("response code %v", rr.Code)
			body := rr.Body
			var records []ServerError
			if err = json.Unmarshal(body.Bytes(), &records); err == nil {
				for _, rec := range records {
					fmt.Printf("ERROR: %s", rec.DBSError)
				}
			} else {
				fmt.Println("### body\n", body.String())
			}
		}
		// dump response if necessary
		if v.DumpResponse {
			fmt.Printf("Server response:\n%+v\n", rr.Body.String())
		}
		// check response code
		if rr.Code != v.Code {
			msg := fmt.Sprintf("ERROR: wrong response code, expect=%d received=%d", v.Code, rr.Code)
			t.Fatal(msg)
		}

		// check response
		var d []map[string]any
		if v.Method == "GET" {
			data := rr.Body.Bytes()
			err = json.Unmarshal(data, &d)
			//             err = json.NewDecoder(rr.Body).Decode(&d)
			if err != nil {
				t.Fatalf("Failed to decode body, %v", err)
			}
			// check output patterns
			for _, o := range v.Output {
				if o == "" {
					continue
				}
				pat, err := regexp.Compile(o)
				if err != nil {
					t.Fatal(err)
				}
				if matched := pat.MatchString(string(data)); !matched {
					if v.Fail {
						t.Logf("We successfully fail while checking pattern '%s'", o)
					} else {
						msg := fmt.Sprintf("Pattern '%s' does not match received data %s", o, string(data))
						t.Fatal(msg)
					}
				}
			}
		} else if v.Method == "POST" {
			data := rr.Body.Bytes()
			if len(data) != 0 {
				err = json.Unmarshal(data, &d)
				if err != nil {
					msg := fmt.Sprintf("Response data '%s', error %v", string(data), err)
					t.Fatal(msg)
				}
			}
		}
		if v.Fail {
			t.Logf("The %s method to %s endpoint must fail with code %v\n", v.Method, v.Endpoint, v.Code)
		}

	})
}

// TestIntegration provides integration tests
func TestIntegration(t *testing.T) {
	idir := os.Getenv("DBS_INT_TESTS_DIR")
	if idir == "" {
		return
	}
	// loop over files in int dir
	files, err := ioutil.ReadDir(idir)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if f.Name() == "README.md" {
			continue
		}
		fmt.Println("\n+++ run integration tests from", f.Name())
		var testCases []TestCase
		if !strings.HasPrefix(f.Name(), "int_") {
			continue
		}
		fname := filepath.Join(idir, f.Name())
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
		for i, v := range testCases {
			if i == 0 && v.Verbose > 0 {
				t.Logf("reading test file %s", fname)
			}
			runTestWorkflow(t, v)
		}
	}
}
