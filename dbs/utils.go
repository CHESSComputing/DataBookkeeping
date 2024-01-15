package dbs

// DBS utils module
// Copyright (c) 2023 - Valentin Kuznetsov <vkuznet@gmail.com>
//
import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
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

// ReplaceBinds replaces given pattern in string
func ReplaceBinds(stm string) string {
	regexp, err := regexp.Compile(`:[a-zA-Z_0-9]+`)
	if err != nil {
		log.Fatal(err)
	}
	match := regexp.ReplaceAllString(stm, "?")
	return match
}

// ConvertFloat converts string representation of float scientific number to string int
func ConvertFloat(val string) string {
	if strings.Contains(val, "e+") || strings.Contains(val, "E+") {
		// we got float number, should be converted to int
		v, e := strconv.ParseFloat(val, 64)
		if e != nil {
			log.Println("unable to convert", val, " to float, error", e)
			return val
		}
		return strings.Split(fmt.Sprintf("%f", v), ".")[0]
	}
	return val
}

// PrintSQL prints SQL/args
func PrintSQL(stm string, args []interface{}, msg string) {
	if msg != "" {
		log.Println(msg)
	} else {
		log.Println("")
	}
	log.Printf("### SQL statement ###\n%s\n\n", stm)
	var values string
	for _, v := range args {
		values = fmt.Sprintf("%s\t'%v'\n", values, v)
	}
	log.Printf("### SQL values ###\n%s\n", values)
}

// ListFiles lists files in a given directory
func ListFiles(dir string) []string {
	var out []string
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du1: %v\n", err)
		return nil
	}
	for _, f := range entries {
		if !f.IsDir() {
			out = append(out, f.Name())
		}
	}
	return out
}

// CastString function to check and cast interface{} to string data-type
func CastString(val interface{}) (string, error) {
	switch v := val.(type) {
	case string:
		return v, nil
	}
	msg := fmt.Sprintf("wrong data type for %v type %T", val, val)
	return "", errors.New(msg)
}

// CastInt function to check and cast interface{} to int data-type
func CastInt(val interface{}) (int, error) {
	switch v := val.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	}
	msg := fmt.Sprintf("wrong data type for %v type %T", val, val)
	return 0, errors.New(msg)
}

// CastInt64 function to check and cast interface{} to int64 data-type
func CastInt64(val interface{}) (int64, error) {
	switch v := val.(type) {
	case int:
		return int64(v), nil
	case int64:
		return v, nil
	}
	msg := fmt.Sprintf("wrong data type for %v type %T", val, val)
	return 0, errors.New(msg)
}
