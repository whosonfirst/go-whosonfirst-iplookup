package iplookup

import (
	"errors"
	_ "fmt"
	"github.com/oschwald/maxminddb-golang"
	"github.com/whosonfirst/go-whosonfirst-csvdb"
	"github.com/whosonfirst/go-whosonfirst-log"
	"net"
	"strconv"
	"time"
)

type MMDBResponse struct {
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
	concordances *csvdb.CSVDB
	logger       *log.WOFLogger
}

func NewIPLookup(db string, meta string, logger *log.WOFLogger) (*IPLookup, error) {

	mmdb, err := maxminddb.Open(db)

	if err != nil {
		return nil, err
	}

	to_index := make([]string, 0)
	to_index = append(to_index, "gn:id")

	t1 := time.Now()

	concordances := csvdb.NewCSVDB()
	err = concordances.IndexCSVFile(meta, to_index)

	t2 := time.Since(t1)

	logger.Debug("time to index concordances %v", t2)

	if err != nil {
		return nil, err
	}

	ip := IPLookup{
		mmdb:         mmdb,
		concordances: concordances,
		logger:       logger,
	}

	return &ip, nil
}

func (ip *IPLookup) Query(addr net.IP) (int64, error) {

	var rsp MMDBResponse
	err := ip.mmdb.Lookup(addr, &rsp)

	if err != nil {
		return -1, err
	}

	possible := make([]uint64, 0)

	possible = append(possible, rsp.City.GeonameId)
	possible = append(possible, rsp.Country.GeonameId)

	ip.logger.Debug("possible matches for %v: %v", addr, possible)

	for _, gnid := range possible {

		if gnid == 0 {
			continue
		}

		wofid, err := ip.ConcordifyGeonames(gnid)

		if err != nil {
			continue
		}

		return wofid, nil
	}

	return -1, errors.New("Unabled to lookup address")
}

func (ip *IPLookup) ConcordifyGeonames(gnid uint64) (int64, error) {

	ip.logger.Debug("concordify geonames %d", gnid)

	str_gnid := strconv.FormatUint(gnid, 10)

	rows, err := ip.concordances.Where("gn:id", str_gnid)

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

	ip.logger.Debug("geonames ID (%d) is WOF ID %d", gnid, wofid)
	return wofid, nil
}
