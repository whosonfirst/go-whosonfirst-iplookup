prep:
	if test -d pkg; then rm -rf pkg; fi

self:	prep
	if test -d src/github.com/whosonfirst/go-whosonfirst-iplookup; then rm -rf src/github.com/whosonfirst/go-whosonfirst-iplookup; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-iplookup
	cp iplookup.go src/github.com/whosonfirst/go-whosonfirst-iplookup/

fmt:
	go fmt cmd/*.go
	go fmt *.go

deps:
	@GOPATH=$(shell pwd) go get -u "github.com/oschwald/maxminddb-golang"
	@GOPATH=$(shell pwd) go get -u "github.com/whosonfirst/go-whosonfirst-csvdb"
	@GOPATH=$(shell pwd) go get -u "github.com/whosonfirst/go-whosonfirst-log"

bin:	self
	@GOPATH=$(shell pwd) go	build -o bin/wof-iplookup cmd/wof-iplookup.go
	@GOPATH=$(shell pwd) go	build -o bin/wof-iplookup-server cmd/wof-iplookup-server.go
