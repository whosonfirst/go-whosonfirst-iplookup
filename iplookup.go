package iplookup

import (
	"github.com/whosonfirst/go-whosonfirst-mmdb"
	"net"
)

type Provider interface {
	QueryString(string) (Result, error)
	Query(net.IP) (Result, error)
}

// PLEASE REPLACE ME WITH VANILLA go-whosonfirst-spr... maybe?

type Result interface {
	Latitude() float64
	Longitude() float64
}

type InflateResponseFunc func(i interface{}) (Result, error)

type SPRResult struct {
	Result `json:",omitempty"`
	SPR    mmdb.SPRRecord `json:"spr"`
}

func (r *SPRResult) Latitude() float64 {
	return r.SPR.Latitude
}

func (r *SPRResult) Longitude() float64 {
	return r.SPR.Longitude
}

func SPRRecordToResult(x interface{}) (Result, error) {

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

	r := SPRResult{
		SPR: s,
	}

	return &r, nil
}
