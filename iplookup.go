package iplookup

import (
	"errors"
	_ "fmt"
	"github.com/oschwald/maxminddb-golang"
	"github.com/whosonfirst/go-whosonfirst-csvdb"
	"github.com/whosonfirst/go-whosonfirst-log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
	See this... it is the dreaded "global variable" - I don't like
	it either but on the other hand it works. It's only ever instantiated
	when someone specifies a "concordances#/path/to/csvfile" source
	which isn't likely to be ever but since it's a useful piece of
	functionality we're going to keep it around for the time being
	(20160113/thisisaaronland)
*/

var concordances *csvdb.CSVDB

type Response interface {
	WOFId() int64
}

type WOFResponse struct {

	/*
	   See also: https://github.com/whosonfirst/p5-Whosonfirst-MaxMind-Writer/blob/master/lib/Whosonfirst/MaxMind/Types.pm
	   See the way we're prefixing `whosonfirst` with maxmindb... yeah, I'm not sure either...
	*/

	WhosonfirstId uint64 `maxminddb:"whosonfirst_id"`
}

func (rsp WOFResponse) WOFId() int64 {
	wofid := int64(rsp.WhosonfirstId)
	return wofid
}

type WOFConcordanceResponse struct {
	Country struct {
		GeonameId uint64 `maxminddb:"geoname_id"`
	} `maxminddb:"country"`
	City struct {
		GeonameId uint64 `maxminddb:"geoname_id"`
	} `maxminddb:"city"`
}

func (rsp WOFConcordanceResponse) WOFId() int64 {

	possible := make([]uint64, 0)

	possible = append(possible, rsp.City.GeonameId)
	possible = append(possible, rsp.Country.GeonameId)

	for _, gnid := range possible {

		if gnid == 0 {
			continue
		}

		wofid, err := rsp.ConcordifyGeonames(gnid)

		if err != nil {
			continue
		}

		return wofid
	}

	return 0
}

func (rsp WOFConcordanceResponse) ConcordifyGeonames(gnid uint64) (int64, error) {

	str_gnid := strconv.FormatUint(gnid, 10)

	// fmt.Printf("look for %s\n", str_gnid)

	rows, err := concordances.Where("gn:id", str_gnid)
	// fmt.Printf("error is %v\n", err)

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

type MaxMindResponse struct {
	Country struct {
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

	logger.Debug("create new IP lookup using %s (%s)", db, source)

	mmdb, err := maxminddb.Open(db)

	if err != nil {
		return nil, err
	}

	ip := IPLookup{
		mmdb:   mmdb,
		source: source,
		logger: logger,
	}

	if strings.HasPrefix(source, "concordances#") {

		parts := strings.Split(source, "#")

		if len(parts) < 2 {
			return nil, errors.New("concordances string is missing a data source")
		}

		data := parts[1]

		_, err := os.Stat(data)

		if os.IsNotExist(err) {
			logger.Error("%s does not exist", data)
			return nil, err
		}

		ip.logger.Debug("loading concordances database %s", data)

		// Waiting on the 'reload' branch of go-wof-csvdb to be
		// pushed to master (20160113/thisisaaronland)

		db := csvdb.NewCSVDB()

		/*
		   db, err := csvdb.NewCSVDB()

		   if err != nil {
		      return nil, err
		   }
		*/

		t1 := time.Now()

		to_index := []string{"gn:id"}

		/*
			Remember, this still needs to be a Geonames ID. Other
			sources would be nice but will require a bit more thinking
			about where we keep track of what we're looking for and
			their types so not today (20160113/thisisaaronland)
		*/

		if len(parts) == 3 {
			to_index = strings.Split(parts[2], ",")
		}

		ip.logger.Debug("indexing %s", strings.Join(to_index, ","))

		err = db.IndexCSVFile(data, to_index)

		t2 := time.Since(t1)
		ip.logger.Debug("time to index concordances: %v", t2)

		if err != nil {
			return nil, err
		}

		ip.source = "concordances"
		concordances = db
	}

	return &ip, nil
}

func (ip *IPLookup) QueryId(addr net.IP) (int64, error) {

	rsp, err := ip.Query(addr)

	if err != nil {
		return 0, err
	}

	wofid := rsp.WOFId()
	return wofid, nil
}

func (ip *IPLookup) QueryRaw(addr net.IP) (interface{}, error) {

	var rsp interface{}
	err := ip.mmdb.Lookup(addr, &rsp)

	if err != nil {
		return nil, err
	}

	return rsp, err
}

func (ip *IPLookup) Query(addr net.IP) (Response, error) {

	var rsp Response
	var err error

	if ip.source == "whosonfirst" || ip.source == "wof" {
		rsp, err = ip.query_wof(addr)
	} else if ip.source == "concordances" {
		rsp, err = ip.query_concordances(addr)
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

func (ip *IPLookup) query_concordances(addr net.IP) (Response, error) {

	var rsp WOFConcordanceResponse
	err := ip.mmdb.Lookup(addr, &rsp)

	if err != nil {
		return nil, err
	}

	return rsp, nil

}
