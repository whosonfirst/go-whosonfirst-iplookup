prep:
	if test -d pkg; then rm -rf pkg; fi

self:	prep
	if test -d src/github.com/whosonfirst/go-whosonfirst-iplookup; then rm -rf src/github.com/whosonfirst/go-whosonfirst-iplookup; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-iplookup
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-iplookup/http
	cp iplookup.go src/github.com/whosonfirst/go-whosonfirst-iplookup/
	cp http/*.go src/github.com/whosonfirst/go-whosonfirst-iplookup/http/
	cp -r vendor/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

deps:	rmdeps
	@GOPATH=$(shell pwd) go get -u "github.com/whosonfirst/go-whosonfirst-spr"

vendor-deps: deps
	if test -d vendor; then rm -rf vendor; fi
	cp -r src vendor
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src
fmt:
	go fmt cmd/*.go
	go fmt http/*.go
	go fmt *.go


bin:	self
	# @GOPATH=$(shell pwd) go	build -o bin/wof-iplookup cmd/wof-iplookup.go
