package iplookup

import (
	"errors"
	_ "fmt"
	"github.com/oschwald/maxminddb-golang"
	"net"
)

type Response struct {
	Country struct {
		ISOCode   string `maxminddb:"iso_code"`
		GeonameId uint64 `maxminddb:"geoname_id"`
	} `maxminddb:"country"`
	City struct {
		GeonameId uint64 `maxminddb:"geoname_id"`
	} `maxminddb:"city"`
}

type Lookup struct {
	mmdb *maxminddb.Reader
}

func NewLookup(db string) (*Lookup, error) {

	mmdb, err := maxminddb.Open(db)

	if err != nil {
		return nil, err
	}

	lookup := Lookup{mmdb}
	return &lookup, nil
}

func (l *Lookup) Query(addr net.IP) (int64, error) {

	var rsp Response
	err := l.mmdb.Lookup(addr, &rsp)

	if err != nil {
		return -1, err
	}

	city_id := rsp.City.GeonameId
	country_id := rsp.Country.GeonameId

	if city_id == 0 && country_id == 0 {
		return -1, errors.New("Unable to locate address")
	}

	return 0, nil
}
