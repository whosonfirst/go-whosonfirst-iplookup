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

type Response interface {
	WOFId() int64
}

type WOFResponse struct {
	WhosonfirstId uint64 `maxminddb:"whosonfirst_id"`
}

func (rsp WOFResponse) WOFId() int64 {
	wofid := int64(rsp.WhosonfirstId)
	return wofid
}

type WOFConcordanceResponse struct {
	Country struct {
		// ISOCode       string `maxminddb:"iso_code"`
		GeonameId uint64 `maxminddb:"geoname_id"`
	} `maxminddb:"country"`
	City struct {
		GeonameId uint64 `maxminddb:"geoname_id"`
	} `maxminddb:"city"`
}

func (rsp WOFConcordanceResponse) WOFId() int64 {
	return 0

	/*
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
	*/
}

/*
func (rsp WOFConcordanceResponse) ConcordifyGeonames(gnid uint64) (int64, error) {

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
*/

type MaxMindResponse struct {
	Country struct {
		// ISOCode       string `maxminddb:"iso_code"`
		GeonameId     uint64 `maxminddb:"geoname_id"`
		WhosonfirstId uint64 `maxminddb:"whosonfirst_id"`
	} `maxminddb:"country"`
	City struct {
		GeonameId     uint64 `maxminddb:"geoname_id"`
		WhosonfirstId uint64 `maxminddb:"whosonfirst_id"`
	} `maxminddb:"city"`
}

func (rsp MaxMindResponse) WOFId() int64 {

	candidate := rsp.City.WhosonfirstId

	if candidate == 0 {
		candidate = rsp.Country.WhosonfirstId
	}

	wofid := int64(candidate)
	return wofid
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

	rsp, err := ip.QueryRaw(addr)

	if err != nil {
		return 0, err
	}

	wofid := rsp.WOFId()
	return wofid, nil
}

func (ip *IPLookup) QueryRaw(addr net.IP) (Response, error) {

	var rsp Response
	var err error

	if ip.source == "wof" {
		rsp, err = ip.query_wof(addr)
	} else {
		rsp, err = ip.query_maxmind(addr)
	}

	return rsp, err
}

func (ip *IPLookup) query_wof(addr net.IP) (Response, error) {

	var rsp WOFResponse
	err := ip.mmdb.Lookup(addr, &rsp)

	if err != nil {
		return nil, err
	}

	return rsp, nil
}

func (ip *IPLookup) query_maxmind(addr net.IP) (Response, error) {

	var rsp MaxMindResponse
	err := ip.mmdb.Lookup(addr, &rsp)

	if err != nil {
		return nil, err
	}

	return rsp, nil
}
