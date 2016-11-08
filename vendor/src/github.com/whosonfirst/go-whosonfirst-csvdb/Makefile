prep:
	if test -d pkg; then rm -rf pkg; fi

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

self:   prep
	if test -d src/github.com/whosonfirst/go-whosonfirst-csvdb; then rm -rf src/github.com/whosonfirst/go-whosonfirst-csvdb; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-csvdb/
	cp csvdb.go src/github.com/whosonfirst/go-whosonfirst-csvdb/
	cp -r vendor/src/* src/

deps:   
	@GOPATH=$(shell pwd) go get -u "github.com/whosonfirst/go-whosonfirst-csv"
	@GOPATH=$(shell pwd) go get -u "github.com/whosonfirst/go-whosonfirst-log"
	@GOPATH=$(shell pwd) go get -u "github.com/whosonfirst/go-whosonfirst-utils"
	@GOPATH=$(shell pwd) go get -u "github.com/go-fsnotify/fsnotify"

fmt:
	go fmt *.go
	go fmt cmd/*.go

bin: 	self
	@GOPATH=$(shell pwd) go build -o bin/wof-csvdb-index cmd/wof-csvdb-index.go
	@GOPATH=$(shell pwd) go build -o bin/wof-csvdb-server cmd/wof-csvdb-server.go

vendor-deps: deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor/src; then rm -rf vendor/src; fi
	cp -r src vendor/src
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src
