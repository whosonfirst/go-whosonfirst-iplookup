# go-whosonfirst-iplookup

Go package for doing IP address to Who's On First record lookups

## Usage

### Example

```
import (
	"github.com/whosonfirst/go-whosonfirst-iplookup"
	"github.com/whosonfirst/go-whosonfirst-log"
	"io"
	"net"
	"os"
)

mmdb := "GeoLite2-City.mmdb"
source := "concordances#wof-concordances-latest.csv"

// Or maybe something like this
// db := "whosonfirst-city-latest.mmdb"
// source := "whosonfirst"
// See the documentation on "sources" below for details

addr := "142.213.160.134"
ip := net.ParseIP(addr)

logger := log.NewWOFLogger("[wof-iplookup] ")
logger.AddLogger(writer, "warning")

// Note the lack of error-handling

lookup, _ := iplookup.NewIPLookup(mmdb, source, logger)
wofid, _ := lookup.QueryId(ip)
```
## Sources

### maxmind

When you specify source as `maxmind` you are telling the code "I have a standard MaxMind GeoLite2 database that has been _augmented_ with Who's On First IDs.

Please remember that these do _not_ include the default GeoLite2 that MaxMind distributes. You can either build a WOF-enabled GeoLite2 database using the [p5-Whosonfirst-MaxMind-Writer](https://github.com/whosonfirst/p5-Whosonfirst-MaxMind-Writer) package or by downloading copies that Who's On First maintains. Links to the latter are included below.

### whosonfirst

When you specify source as `whosonfirst` you are telling the code "I have a non-standard MaxMind database that has been built using the [p5-Whosonfirst-MaxMind-Writer](https://github.com/whosonfirst/p5-Whosonfirst-MaxMind-Writer) package or that you downloaded from Who's On First. Like the WOF-enabled GeoLite2 databases links are included below.

_These non-standard MaxMind databases are still considered experimental. While they will always contain a `whosonfirst_id` property it is still possible that other things may change._

### concordances

First the `concordances` source is actually `concordances#/some/path/to-a/concordances-file.csv`, most likely the [wof-concordances-latest.csv](https://github.com/whosonfirst/whosonfirst-data/blob/master/meta/wof-concordances-latest.csv) file included in the `whosonfirst-data` repository.

Okay, now that that's out of the way when you specify source as `concordances` you are telling the code "I have standard MaxMind GeoLite2 database that doesn't contain any Who's On First information _but_ I do have this handy CSV file that maps Geonames IDs to Who's On First IDs so please use that, okay?".

## Utilities

### wof-iplookup

```
$> ./bin/wof-iplookup -h
Usage of ./bin/wof-iplookup:
  -db string
      The path to your IP lookup database file
  -json
	Dump the raw query response as JSON
  -loglevel string
    	     (default "warning")
  -raw
	Return the raw data
  -source string
    	  The source of the IP lookups (default "maxmind")
```

Perform an IP lookup for a list of IP addresses passed on the command line. By default this emits a single Who's On First ID on a new line for each IP address passed in as an argument.

Here's an example using a `concordances` file as the input source and verbose logging, just so you can see what's going on under the hood:

```
$> ./bin/wof-iplookup -loglevel debug -db /usr/local/mapzen/mmdb/wof-mm-city.mmdb -source concordances#/usr/local/mapzen/whosonfirst-data/meta/wof-concordances-latest.csv 8.8.8.4
[wof-iplookup] 16:22:59.667597 [debug] create new IP lookup using /usr/local/mapzen/mmdb/wof-mm-city.mmdb (concordances#/usr/local/mapzen/whosonfirst-data/meta/wof-concordances-latest.csv)
[wof-iplookup] 16:22:59.667968 [debug] loading concordances database /usr/local/mapzen/whosonfirst-data/meta/wof-concordances-latest.csv
[wof-iplookup] 16:22:59.667984 [debug] indexing gn:id
[wof-iplookup] 16:23:01.935802 [debug] time to index concordances: 2.267799616s
[wof-iplookup] 16:23:01.935829 [debug] lookup 8.8.8.4
85922355
```

If you want to see the complete response coming back from the database you can pass in the `-raw` flag which will print an encoded JSON string for each IP address passed in as an argument. 

Here's an example using a `maxmind` IP data as the input source:

```
$> ./bin/wof-iplookup -json -db /usr/local/mapzen/mmdb/wof-mm-city.mmdb -source maxmind 142.213.160.134 | python -mjson.tool
{
    "City": {
        "GeonameId": 0,
        "WhosonfirstId": 0
    },
    "Country": {
        "GeonameId": 6251999,
        "WhosonfirstId": 85633041
    }
}
```

### wof-iplookup-server

```
$> ./bin/wof-iplookup-server  -h
Usage of ./bin/wof-iplookup-server:
  -cors
	Enable CORS headers
  -db string
      The path to your IP lookup database file
  -host string
    	The hostname to listen for requests on (default "localhost")
  -loglevel string
    	     (default "status")
  -port int
    	The port number to listen for requests on (default 8668)
  -source string
    	  The source of the IP lookups (default "maxmind")
```

A handy HTTP pony for performing IP lookups as a service.

```
$> ./bin/wof-iplookup-server -loglevel info -db ~/usr/local/mapzen/whosonfirst-city-20160111.mmdb -source whosonfirst 
[wof-iplookup] 16:51:21.949622 [status] wof-iplookup-server running at localhost:8668
```

And then:

```
$> curl -s 'http://localhost:8668?ip=205.193.117.158' | python -mjson.tool
{
    "ip": "205.193.117.158",
    "wofid": 101735873
}
```

You can also pass in a `raw=ANYTHING` parameter to get all the data for a given record. Like this:

```
curl -s 'http://localhost:8668?ip=205.193.117.158&raw=1' | python -mjson.tool
{
    "continent_id": 102191575,
    "country_id": 85633041,
    "disputed_id": 0,
    "geom_bbox": "-76.3555857,44.9617738,-75.2465783,45.5376514",
    "geom_latitude": 45.293133,
    "geom_longitude": -75.775424,
    "geoname_id": 6094817,
    "lbl_latitude": 45.209415,
    "lbl_longitude": -75.783876,
    "localadmin_id": 0,
    "locality_id": 101735873,
    "macroregion_id": 0,
    "mm_latitude": 45.3548,
    "mm_longitude": -75.5773,
    "name": "Ottawa",
    "placetype": "locality",
    "region_id": 85682057,
    "whosonfirst_id": 101735873
}
```

## Caveats

### Metadata

Neither the `wof-iplookup-server` or the `ip-lookup` tool (when run with the `-raw` flag) return any additional metadata data for a WOF record besides a Who's On First ID. For example the default GeoLite2 databases contain lots of useful place names and the WOF-enabled derivatives contain bounding box and other useful geographic information. Once the basic database wrangling settles down a bit it will make sense to start bubbling that information back up to consumer applications.

## See also

* https://whosonfirst.mapzen.com/mmdb/
* https://github.com/whosonfirst/p5-Whosonfirst-MaxMind-Writer
* https://github.com/oschwald/maxminddb-golang
* https://dev.maxmind.com/geoip/geoip2/geolite2/
