package provider

import (
	"github.com/oschwald/maxminddb-golang"
	"github.com/whosonfirst/go-whosonfirst-iplookup"
	"net"
)

type MMDBProvider struct {
	iplookup.Provider
	handler iplookup.InflateResponseFunc
	db      *maxminddb.Reader
}

func NewMMDBProvider(path string, handler iplookup.InflateResponseFunc) (iplookup.Provider, error) {

	db, err := maxminddb.Open(path)

	if err != nil {
		return nil, err
	}

	pr := MMDBProvider{
		db:      db,
		handler: handler,
	}

	return &pr, nil
}

func (pr *MMDBProvider) QueryString(str_addr string) (iplookup.Result, error) {

	addr := net.ParseIP(str_addr)
	return pr.Query(addr)
}

func (pr *MMDBProvider) Query(addr net.IP) (iplookup.Result, error) {

	var i interface{}
	err := pr.db.Lookup(addr, &i)

	if err != nil {
		return nil, err
	}

	r, err := pr.handler(i)

	if err != nil {
		return nil, err
	}

	return r, nil
}
