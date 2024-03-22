# DataBookkeeping Service

![build status](https://github.com/CHESSComputing/DataBookkeeping/actions/workflows/go.yml/badge.svg)
[![go report card](https://goreportcard.com/badge/github.com/CHESSComputing/DataBookkeeping)](https://goreportcard.com/report/github.com/CHESSComputing/DataBookkeeping)
[![godoc](https://godoc.org/github.com/CHESSComputing/DataBookkeeping?status.svg)](https://godoc.org/github.com/CHESSComputing/DataBookkeeping)

Data Bookkeeping service

### APIs

#### public APIs
- `/datasets` get all datasets
- `/files` get all files
- `/dataset/*name` get dataset with given name
- `/file/*name` get file with given name

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
  "buckets": [
    "bucketABC"
  ],
  "did": "/a/b/c",
  "files": [
    "/path/file1.png",
    "/path/file2.png",
    "/path/file3.png"
  ],
  "processing": "glibc",
  "site": "Cornell"
}

# inject new record
curl -v -X POST -H "Authorization: Bearer $token" \
    -H "Content-type: application/json" \
    -d@./record.json \
    http://localhost:8310/dataset
```
