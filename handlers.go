package main

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/CHESSComputing/DataBookkeeping/dbs"
	"github.com/CHESSComputing/DataBookkeeping/utils"
	server "github.com/CHESSComputing/golib/server"
	"github.com/gin-gonic/gin"
)

type NameRequest struct {
	Name string `uri:"name" json:"name"`
}

// ApisHandler provides all server routes
func ApisHandler(c *gin.Context) {
	server.ApisHandler(c, _routes)
}

// FileHandler provides access to /files and /file/:name end-point
func FileHandler(c *gin.Context) {
	ApiHandler(c, "file")
}

// DatasetHandler provides access to GET /datasets and /dataset/:name end-point
func DatasetHandler(c *gin.Context) {
	ApiHandler(c, "dataset")
}

// ApiHandler represents generic API handler for GET/POST/PUT/DELETE requests of a specific API
func ApiHandler(c *gin.Context, api string) {
	r := c.Request
	if r.Method == "POST" {
		DBSPostHandler(c, api)
	} else if r.Method == "PUT" {
		DBSPutHandler(c, api)
	} else if r.Method == "DELETE" {
		DBSDeleteHandler(c, api)
	} else {
		DBSGetHandler(c, api)
	}
}

// helper function to get DBS API
func getApi(c *gin.Context, a string) (*dbs.API, error) {
	r := c.Request
	w := c.Writer
	// all outputs will be added to output list
	sep := ","
	if r.Header.Get("Accept") == "application/ndjson" {
		sep = ""
	}
	if sep != "" {
		w.Header().Add("Content-Type", "application/json")
	} else {
		w.Header().Add("Content-Type", "application/ndjson")
	}

	var api *dbs.API
	params := make(dbs.Record)
	if r.Method == "GET" {
		// for example /file?dataset=/x/y/z we'll parse URL query
		// r.URL.Query() returns map[string][]string
		for k, values := range r.URL.Query() {
			var vals []string
			for _, v := range values {
				vals = append(vals, v)
			}
			params[k] = vals
		}
	}
	if r.Method == "GET" || r.Method == "DELETE" {
		api = &dbs.API{
			Writer:      w,
			Params:      params,
			Separator:   sep,
			Api:         a,
			ContentType: r.Header.Get("Content-Type"),
		}
	} else { // all other HTTP requests POST/PUT may contain payload

		headerContentType := r.Header.Get("Content-Type")
		if headerContentType != "application/json" {
			msg := fmt.Sprintf("unsupported Content-Type: '%s'", headerContentType)
			e := dbs.Error(dbs.ContentTypeErr, dbs.ContentTypeErrorCode, msg, "web.DBSPostHandler")
			responseMsg(w, r, e, http.StatusUnsupportedMediaType)
			return nil, errors.New(msg)
		}
		defer r.Body.Close()
		//         var params dbs.Record
		if utils.VERBOSE > 0 {
			dn, _ := r.Header["Cms-Authn-Dn"]
			log.Printf("DBSPostHandler: API=%s, dn=%s, uri=%s", a, dn, requestURI(r))
		}
		cby := createBy(r)
		body := r.Body
		// handle gzip content encoding
		if r.Header.Get("Content-Encoding") == "gzip" {
			r.Header.Del("Content-Length")
			reader, err := gzip.NewReader(r.Body)
			if err != nil {
				msg := "unable to get gzip reader"
				log.Println(msg, err)
				e := dbs.Error(err, dbs.ReaderErrorCode, msg, "web.DBSPostHandler")
				responseMsg(w, r, e, http.StatusInternalServerError)
				return nil, errors.New(msg)
			}
			body = utils.GzipReader{reader, r.Body}
		} else {
			data, err := io.ReadAll(r.Body)
			if err != nil {
				msg := "unable to get io reader"
				log.Println(msg, err)
				e := dbs.Error(err, dbs.ReaderErrorCode, msg, "web.DBSPostHandler")
				responseMsg(w, r, e, http.StatusInternalServerError)
				return nil, errors.New(msg)
			}
			body = ioutil.NopCloser(bytes.NewBuffer(data))
		}
		api = &dbs.API{
			Reader:      body,
			Writer:      w,
			Params:      params,
			Separator:   sep,
			CreateBy:    cby,
			Api:         a,
			ContentType: r.Header.Get("Content-Type"),
		}
	}
	/*
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			w.Header().Set("Content-Encoding", "gzip")
			api.ContentType = "gzip"
			gw := gzip.NewWriter(w)
			defer gw.Close()
			api.Writer = utils.GzipWriter{GzipWriter: gw, Writer: w}
		}
	*/

	// many APIs carry /api/*name RESTful end-point and we need to
	// get name out of it
	var rest NameRequest
	if err := c.ShouldBindUri(&rest); err == nil {
		if a == "dataset" && rest.Name != "" {
			api.Params["dataset"] = rest.Name
		} else if a == "file" && rest.Name != "" {
			api.Params["logical_file_name"] = rest.Name
		}
	}

	if utils.VERBOSE > 0 {
		log.Println("Call DBS API", api.String())
	}
	return api, nil
}

// DBSGetHandler is a generic Get handler to call DBS Get APIs.
//
//gocyclo:ignore
func DBSGetHandler(c *gin.Context, a string) {
	r := c.Request
	w := c.Writer
	api, err := getApi(c, a)
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
		api.ContentType = "gzip"
		gw := gzip.NewWriter(w)
		defer gw.Close()
		api.Writer = utils.GzipWriter{GzipWriter: gw, Writer: w}
	}
	if err != nil {
		responseMsg(w, r, err, http.StatusBadRequest)
	}
	if a == "dataset" {
		err = api.GetDataset()
	} else if a == "file" {
		err = api.GetFile()
	} else {
		err = dbs.NotImplementedApiErr
	}
	if err != nil {
		responseMsg(w, r, err, http.StatusBadRequest)
		return
	}
}

// POST handler
//
// DBSPostHandler is a generic handler to call DBS Post APIs
//
//gocyclo:ignore
func DBSPostHandler(c *gin.Context, a string) {
	r := c.Request
	w := c.Writer
	api, err := getApi(c, a)
	if err != nil {
		responseMsg(w, r, err, http.StatusBadRequest)
	}
	if a == "dataset" {
		err = api.InsertDataset()
	} else if a == "file" {
		err = api.InsertFile()
	} else {
		err = dbs.NotImplementedApiErr
	}
	if err != nil {
		responseMsg(w, r, err, http.StatusBadRequest)
		return
	}
}

// DBSPutHandler is a generic handler to call DBS put APIs
//
//gocyclo:ignore
func DBSPutHandler(c *gin.Context, a string) {
	r := c.Request
	w := c.Writer
	api, err := getApi(c, a)
	if err != nil {
		responseMsg(w, r, err, http.StatusBadRequest)
	}
	if a == "dataset" {
		err = api.UpdateDataset()
	} else if a == "file" {
		err = api.UpdateFile()
	} else {
		err = dbs.NotImplementedApiErr
	}
	if err != nil {
		responseMsg(w, r, err, http.StatusBadRequest)
		return
	}
}

// DBSDeleteHandler is a generic handler to call DBS delete APIs
//
//gocyclo:ignore
func DBSDeleteHandler(c *gin.Context, a string) {
	r := c.Request
	w := c.Writer
	api, err := getApi(c, a)
	if err != nil {
		responseMsg(w, r, err, http.StatusBadRequest)
	}
	if a == "dataset" {
		err = api.DeleteDataset()
	} else if a == "file" {
		err = api.DeleteFile()
	} else {
		err = dbs.NotImplementedApiErr
	}
	if err != nil {
		responseMsg(w, r, err, http.StatusBadRequest)
		return
	}
}
