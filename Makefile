#GOPATH:=$(PWD):${GOPATH}
#export GOPATH
flags=-ldflags="-s -w"
# flags=-ldflags="-s -w -extldflags -static"
TAG := $(shell git tag | sed -e "s,v,,g" | sort -r | head -n 1)

all: build

gorelease:
	goreleaser release --snapshot --clean

build:
ifdef TAG
	sed -i -e "s,{{VERSION}},$(TAG),g" main.go
endif
	go clean; rm -rf pkg; go build -o srv ${flags}
ifdef TAG
	sed -i -e "s,$(TAG),{{VERSION}},g" main.go
endif

build_all: build_darwin_amd64 build_darwin_arm64 build_amd64 build_arm64 build_power8 build_windows

build_darwin_amd64:
ifdef TAG
	sed -i -e "s,{{VERSION}},$(TAG),g" main.go
endif
	go clean; rm -rf pkg srv_darwin; GOOS=darwin go build -o srv ${flags}
	mv srv srv_darwin_amd64
ifdef TAG
	sed -i -e "s,$(TAG),{{VERSION}},g" main.go
endif

build_darwin_arm64:
ifdef TAG
	sed -i -e "s,{{VERSION}},$(TAG),g" main.go
endif
	go clean; rm -rf pkg srv_darwin; GOARCH=arm64 GOOS=darwin go build -o srv ${flags}
	mv srv srv_darwin_arm64
ifdef TAG
	sed -i -e "s,$(TAG),{{VERSION}},g" main.go
endif

build_amd64:
ifdef TAG
	sed -i -e "s,{{VERSION}},$(TAG),g" main.go
endif
	go clean; rm -rf pkg srv_linux; GOOS=linux go build -o srv ${flags}
	mv srv srv_amd64
ifdef TAG
	sed -i -e "s,$(TAG),{{VERSION}},g" main.go
endif

build_power8:
ifdef TAG
	sed -i -e "s,{{VERSION}},$(TAG),g" main.go
endif
	go clean; rm -rf pkg srv_power8; GOARCH=ppc64le GOOS=linux go build -o srv ${flags}
ifdef TAG
	sed -i -e "s,$(TAG),{{VERSION}},g" main.go
endif
	mv srv srv_power8

build_arm64:
ifdef TAG
	sed -i -e "s,{{VERSION}},$(TAG),g" main.go
endif
	go clean; rm -rf pkg srv_arm64; GOARCH=arm64 GOOS=linux go build -o srv ${flags}
ifdef TAG
	sed -i -e "s,$(TAG),{{VERSION}},g" main.go
endif
	mv srv srv_arm64

build_windows:
ifdef TAG
	sed -i -e "s,{{VERSION}},$(TAG),g" main.go
endif
	go clean; rm -rf pkg srv.exe; GOARCH=amd64 GOOS=windows go build -o srv.exe ${flags}
ifdef TAG
	sed -i -e "s,$(TAG),{{VERSION}},g" main.go
endif

install:
	go install

clean:
	go clean; rm -rf pkg

MONGO := $(shell ps auxww | grep mongo | egrep -v grep)

mongo:
ifndef MONGO
	$(error "mongo process not found, please start it to proceed ...")
endif

testdb:
	/bin/rm -f /tmp/dbs.db && \
	sqlite3 /tmp/dbs.db < ./static/schema/sqlite.sql && \
	mkdir -p /tmp/${USER} && \
	echo "test" > /tmp/${USER}/test.txt

test : testdb test_code

test_code:
	go test -test.v .

# here is an example for execution of individual test
# go test -v -run TestFilesDB

test: test-errors test-dbs

test-errors:
	cd test && LD_LIBRARY_PATH=${odir} DYLD_LIBRARY_PATH=${odir} go test -v -run TestDBSError
test-dbs:
	cd test && rm -f /tmp/dbs-test.db && \
	sqlite3 /tmp/dbs-test.db < ../static/schema/sqlite.sql && \
	LD_LIBRARY_PATH=${odir} DYLD_LIBRARY_PATH=${odir} \
	DBS_DB_FILE=/tmp/dbs-test.db \
	DBS_API_PARAMETERS_FILE=../static/parameters.json \
	DBS_LEXICON_FILE=../static/lexicon_writer.json \
	go test -v -run TestDBS
