package iplookup

import (
	_ "errors"
	_ "fmt"
	"github.com/oschwald/maxminddb-golang"
	"github.com/whosonfirst/go-whosonfirst-log"
	"net"
)

// See also
// https://github.com/whosonfirst/p5-Whosonfirst-MaxMind-Writer/blob/master/lib/Whosonfirst/MaxMind/Types.pm

type WOFResponse struct {
     whosonfirst_id uint64 `maxminddb:"whosonfirst_id"`
}

type MaxMindResponse struct {
	Country struct {
		ISOCode   string `maxminddb:"iso_code"`
		GeonameId uint64 `maxminddb:"geoname_id"`
	} `maxminddb:"country"`
	City struct {
		GeonameId uint64 `maxminddb:"geoname_id"`
	} `maxminddb:"city"`
}

type IPLookup struct {
	mmdb         *maxminddb.Reader
	source	     string
	logger       *log.WOFLogger
}

func NewIPLookup(db string, source string, logger *log.WOFLogger) (*IPLookup, error) {

	mmdb, err := maxminddb.Open(db)

	if err != nil {
		return nil, err
	}

	ip := IPLookup{
		mmdb:         mmdb,
		source:	      source,
		logger:       logger,
	}

	return &ip, nil
}

func (ip *IPLookup) Query(addr net.IP) (int64, error) {

     // please to be reading ip.source...

	var rsp WOFResponse
	err := ip.mmdb.Lookup(addr, &rsp)

	if err != nil {
		return -1, err
	}

	return 0, nil
}
