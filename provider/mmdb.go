package provider

import (
	"github.com/oschwald/maxminddb-golang"
	"github.com/whosonfirst/go-whosonfirst-iplookup"
	"net"
)

type MMDB struct {
	iplookup.Provider
	handler iplookup.InflateResponseFunc
	db      *maxminddb.Reader
}

func NewMMDB(path string, handler iplookup.InflateResponseFunc) (iplookup.Provider, error) {

	db, err := maxminddb.Open(db)

	if err != nil {
		return nil, err
	}

	mm := MMDB{
		db:      db,
		handler: handler,
	}

	return &mm, nil
}

func (mm *MMDB) QueryString(str_addr string) (Result, error) {

	addr := net.ParseIP(str_addr)
	return mm.Query(addr)
}

func (mm *MMDB) Query(addr net.Addr) (Result, error) {

	var i interface{}
	err := ip.mmdb.Lookup(addr, &i)

	if err != nil {
		return nil, err
	}

	r, err := ip.handler(i)

	if err != nil {
		return nil, err
	}

	return r, nil
}
