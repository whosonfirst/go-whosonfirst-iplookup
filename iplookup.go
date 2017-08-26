package iplookup

import (
	_ "fmt"
	"github.com/oschwald/maxminddb-golang"
	"github.com/whosonfirst/go-whosonfirst-log"
	"github.com/whosonfirst/go-whosonfirst-mmdb"
	"net"
)

type Response interface {
	WOFId() int64
}

type SPRResponse struct {
	Response
	spr mmdb.SPRRecord
}

func (r *SPRResponse) WOFId() int64 {
	return r.spr.Id
}

func SPRRecordToReponse(i interface{}) (Response, error) {

	s := i.(mmdb.SPRRecord)

	r := SPRResponse{
		spr: s,
	}

	return &r, nil
}

type IPLookupToResponse func(i interface{}) (Response, error)

type IPLookup struct {
	mmdb     *maxminddb.Reader
	callback IPLookupToResponse
	Logger   *log.WOFLogger
}

func NewIPLookup(db string, cb IPLookupToResponse) (*IPLookup, error) {

	logger := log.SimpleWOFLogger()

	mmdb, err := maxminddb.Open(db)

	if err != nil {
		return nil, err
	}

	ip := IPLookup{
		mmdb:     mmdb,
		callback: cb,
		Logger:   logger,
	}

	return &ip, nil
}

func (ip *IPLookup) QueryString(str_addr string) (Response, error) {

	addr := net.ParseIP(str_addr)
	return ip.Query(addr)
}

func (ip *IPLookup) Query(addr net.IP) (Response, error) {

	var i interface{}
	err := ip.mmdb.Lookup(addr, &i)

	if err != nil {
		return nil, err
	}

	r, err := ip.callback(i)

	if err != nil {
		return nil, err
	}

	return r, nil
}
