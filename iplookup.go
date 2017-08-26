package iplookup

import (
	"github.com/oschwald/maxminddb-golang"
	"github.com/whosonfirst/go-whosonfirst-log"
	"github.com/whosonfirst/go-whosonfirst-mmdb"
	"net"
)

type Provider interface {
	QueryString(string) (Result, error)
	Query(string) (Result, error)
}

type Result interface {
	Latitude() float64
	Longitude() float64
}

type InflateResponseFunc func(i interface{}) (Result, error)

// PLEASE REPLACE ME WITH VANILLA go-whosonfirst-spr

type Response interface {
	WOFId() int64
}

type SPRResponse struct {
	Response `json:",omitempty"`
	SPR      mmdb.SPRRecord `json:"spr"`
}

func (r *SPRResponse) WOFId() int64 {
	return r.SPR.Id
}

func SPRRecordToReponse(x interface{}) (Response, error) {

	// this makes me feel dirty but you know we're
	// still just figuring it out so... there you go
	// (20170825/thisisaaronland)

	i := x.(map[string]interface{})

	s := mmdb.SPRRecord{
		Id:           int64(i["wof:id"].(uint64)),
		Name:         i["wof:name"].(string),
		Placetype:    i["wof:placetype"].(string),
		Latitude:     i["wof:latitude"].(float64),
		Longitude:    i["wof:longitude"].(float64),
		MinLatitude:  i["geom:min_latitude"].(float64),
		MinLongitude: i["geom:min_longitude"].(float64),
		MaxLatitude:  i["geom:max_latitude"].(float64),
		MaxLongitude: i["geom:max_longitude"].(float64),
	}

	r := SPRResponse{
		SPR: s,
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
