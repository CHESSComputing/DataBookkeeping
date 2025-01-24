# DataBookkeeping Service

![build status](https://github.com/CHESSComputing/DataBookkeeping/actions/workflows/go.yml/badge.svg)
[![go report card](https://goreportcard.com/badge/github.com/CHESSComputing/DataBookkeeping)](https://goreportcard.com/report/github.com/CHESSComputing/DataBookkeeping)
[![godoc](https://godoc.org/github.com/CHESSComputing/DataBookkeeping?status.svg)](https://godoc.org/github.com/CHESSComputing/DataBookkeeping)

Data Bookkeeping service

### APIs

#### public APIs
- `/datasets` get all datasets
- `/files` get files for a given did
- `/dataset/*name` get dataset with given name
- `/file/*name` get file with given name
- `/provenance` get provenance information about given did

#### Example
Here are examples of GET HTTP requests
```
# look-up all datasets
curl -v http://localhost:8310/datasets

# look-up concrete dataset=/x/y/z
dataset=/x/y/z
curl -v http://localhost:8310/dataset$dataset

# look-up files from a dataset
curl -v "http://localhost:8310/file?dataset=$dataset"
```

#### protected APIs
- HTTP POST requests
    - `/dataset` create new dataset data
    - `/file` create new file data
- HTTP PUT requests
    - `/dataset` update dataset data
    - `/file` update file data
- HTTP DELETE requests
    - `/dataset/*name` delete dataset
    - `/file/*name` delete file

#### Example

Here is an example of HTTP POST request
```
# record.json
{
    "parent_did": "/beamline=aaa/btr=bbb/cycle=ccc/sample_name=sss",
    "did": "/beamline=aaa/btr=bbb/cycle=ccc/sample_name=sss/test=child",
    "processing": "processing string, e.g. glibc-123-python-123",
    "osinfo": {"name": "linux-cc7", "kernel": "1-2-3", "version": "cc7-123"},
    "environments": [
      {"name": "galaxy", "version": "version", "details": "details",
          "parent_environment": "conda-123", "os_name": "linux-cc7"},
      {"name": "conda-123", "version": "version", "details": "details",
          "parent_environment": null, "os_name": "linux-cc7",
          "packages": [
              {"name": "numpy", "version": "123"},
              {"name": "matplotlib", "version": "987"}
          ]
      }
    ],
    "scripts": [
      {"name": "reader", "options": "-reader_options", "parent_script": null, "order_idx": 1},
      {"name": "chap", "options": "-chap_options", "parent_script": "myscript", "order_idx": 2}
    ],
    "input_files": [
      {"name": "/tmp/file1.png"},
      {"name": "/tmp/file2.png"}
    ],
    "output_files": [
      {"name": "/tmp/file1.png"}
    ],
    "site": "Cornell",
    "buckets": ["bucketABC"]
}

# inject new record
curl -v -X POST -H "Authorization: Bearer $token" \
    -H "Content-type: application/json" \
    -d@./record.json \
    http://localhost:8310/dataset
```

For more (and up-to-date examples) please see `data` integration area of this
repository and look-up JSON input in `int_provenance.json` file.
