prep:
	if test -d pkg; then rm -rf pkg; fi

self:	prep
	if test -d src/github.com/whosonfirst/go-whosonfirst-iplookup; then rm -rf src/github.com/whosonfirst/go-whosonfirst-iplookup; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-iplookup
	cp iplookup.go src/github.com/whosonfirst/go-whosonfirst-iplookup/

deps:
	@GOPATH=$(shell pwd) go get -u "github.com/oschwald/maxminddb-golang"
