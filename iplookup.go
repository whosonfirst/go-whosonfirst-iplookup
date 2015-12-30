package iplookup

import (
	"errors"
	_ "fmt"
	"github.com/oschwald/maxminddb-golang"
	csvdb "github.com/whosonfirst/go-whosonfirst-csvdb"
	"net"
	"strconv"
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
	mmdb         *maxminddb.Reader
	concordances *csvdb.CSVDB
}

func NewLookup(db string, meta string) (*Lookup, error) {

	mmdb, err := maxminddb.Open(db)

	if err != nil {
		return nil, err
	}

	to_index := make([]string, 0)
	to_index = append(to_index, "gn:id")

	concordances := csvdb.NewCSVDB()
	err = concordances.IndexCSVFile(meta, to_index)

	if err != nil {
		return nil, err
	}

	lookup := Lookup{
		mmdb:         mmdb,
		concordances: concordances,
	}

	return &lookup, nil
}

func (l *Lookup) Query(addr net.IP) (int64, error) {

	var rsp Response
	err := l.mmdb.Lookup(addr, &rsp)

	if err != nil {
		return -1, err
	}

	possible := make([]uint64, 0)

	possible = append(possible, rsp.City.GeonameId)
	possible = append(possible, rsp.Country.GeonameId)

	for _, gnid := range possible {

		if gnid == 0 {
			continue
		}

		wofid, err := l.Concordify(gnid)

		if err != nil {
			continue
		}

		return wofid, nil
	}

	return -1, errors.New("Unabled to lookup address")
}

func (l *Lookup) Concordify(gnid uint64) (int64, error) {

	str_gnid := strconv.FormatUint(gnid, 16)
	rows, err := l.concordances.Where("gn:id", str_gnid)

	if err != nil {
		return -1, err
	}

	first := rows[0]
	others := first.AsMap()

	str_wofid, ok := others["wof:id"]

	if !ok {
		return -1, errors.New("Unable to locate concordance")
	}

	wofid, err := strconv.ParseInt(str_wofid, 10, 64)

	if err != nil {
		return -1, err
	}

	return wofid, nil
}
