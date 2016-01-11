package iplookup

import (
	_ "errors"
	_ "fmt"
	"github.com/oschwald/maxminddb-golang"
	"github.com/whosonfirst/go-whosonfirst-log"
	"net"
)

// See also: https://github.com/whosonfirst/p5-Whosonfirst-MaxMind-Writer/blob/master/lib/Whosonfirst/MaxMind/Types.pm
// see the way we're prefixing `whosonfirst` with maxmindb... yeah, I'm not sure either...

type WOFResponse struct {
	WhosonfirstId uint64 `maxminddb:"whosonfirst_id"`
}

type MaxMindResponse struct {
	Country struct {
		ISOCode       string `maxminddb:"iso_code"`
		GeonameId     uint64 `maxminddb:"geoname_id"`
		WhosonfirstId uint64 `maxminddb:"whosonfirst_id"`
	} `maxminddb:"country"`
	City struct {
		GeonameId     uint64 `maxminddb:"geoname_id"`
		WhosonfirstId uint64 `maxminddb:"whosonfirst_id"`
	} `maxminddb:"city"`
}

type IPLookup struct {
	mmdb   *maxminddb.Reader
	source string
	logger *log.WOFLogger
}

func NewIPLookup(db string, source string, logger *log.WOFLogger) (*IPLookup, error) {

	mmdb, err := maxminddb.Open(db)

	if err != nil {
		return nil, err
	}

	ip := IPLookup{
		mmdb:   mmdb,
		source: source,
		logger: logger,
	}

	return &ip, nil
}

func (ip *IPLookup) Query(addr net.IP) (int64, error) {

	/*
	   Please to be reading ip.source and adjusting accordingly...
	   As of this writing I am waiting for the 'append a WOF ID to
	   to standard geolite2 mmdb' databases to complete...
	   (20160110/thisisaaronland)
	*/
	return ip.query_wof(addr)
}

func (ip *IPLookup) query_wof(addr net.IP) (int64, error) {

	var rsp WOFResponse
	err := ip.mmdb.Lookup(addr, &rsp)

	if err != nil {
		return -1, err
	}

	wofid := int64(rsp.WhosonfirstId)
	return wofid, nil
}

/* See notes in Query */

func (ip *IPLookup) query_maxmind(addr net.IP) (int64, error) {

	var rsp MaxMindResponse
	err := ip.mmdb.Lookup(addr, &rsp)

	if err != nil {
		return -1, err
	}

	candidate := rsp.City.WhosonfirstId

	if candidate == 0 {
		candidate = rsp.Country.WhosonfirstId
	}

	wofid := int64(candidate)
	return wofid, nil
}
