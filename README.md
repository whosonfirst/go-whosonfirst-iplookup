# go-whosonfirst-iplookup

Go package for doing IP address to Who's On First record lookups

## Description

This is basically a thin wrapper around the [MaxMindDB geoIP2](https://dev.maxmind.com/geoip/geoip2/geolite2/) lookup databases and @oschwald's [maxminddb-golang](https://github.com/oschwald/maxminddb-golang) package.

By default MaxMindDB results include [Geonames](http://www.geonames.org) IDs for each location. We (Who's On First) are working on [tools to generate copies of the MaxMindDB databases with WOF IDs](https://github.com/whosonfirst/p5-Whosonfirst-MaxMind-Writer) as well. That work is not complete yet.

In the meantime this package accepts a concordances "meta" file from the [whosonfirst-data](https://github.com/whosonfirst/whosonfirst-data/) repository and will attempt to map a Geonames ID to a Who's On First ID. As of this writing not every Geonames ID has a Who's On First concordance. They will but in the event of a failed lookup the package will return an error.

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
concordances := "wof-concordances-latest.csv"

addr := "142.213.160.134"
ip := net.ParseIP(addr)

logger := log.NewWOFLogger("[wof-iplookup] ")
logger.AddLogger(writer, "warning")

// Note the lack of error-handling

lookup, _ := iplookup.NewIPLookup(mmdb, concordances, logger)
wofid, _ := lookup.Query(ip)
```

## Utilities

### wof-iplookup

Perform IP lookup for a list of IP addresses passed on the command line.

```
$> ./bin/wof-iplookup -mmdb GeoLite2-City.mmdb -concordances /usr/local/mapzen/whosonfirst-data/meta/wof-concordances-latest.csv -loglevel debug 8.8.8.8 8.8.8.4 142.213.160.134
[wof-iplookup] 00:20:49.187094 [debug] time to index concordances 3.03230889s
[wof-iplookup] 00:20:49.187217 [debug] possible matches for 8.8.8.8: [5375480 6252001]
[wof-iplookup] 00:20:49.187244 [debug] concordify geonames 5375480
[wof-iplookup] 00:20:49.187272 [debug] geonames ID (5375480) is WOF ID 85922355
[wof-iplookup] 00:20:49.187284 [debug] 8.8.8.8 is 85922355
85922355
[wof-iplookup] 00:20:49.187310 [debug] possible matches for 8.8.8.4: [5375480 6252001]
[wof-iplookup] 00:20:49.187326 [debug] concordify geonames 5375480
[wof-iplookup] 00:20:49.187335 [debug] geonames ID (5375480) is WOF ID 85922355
[wof-iplookup] 00:20:49.187353 [debug] 8.8.8.4 is 85922355
85922355
[wof-iplookup] 00:20:49.187391 [debug] possible matches for 142.213.160.134: [0 6251999]
[wof-iplookup] 00:20:49.187403 [debug] concordify geonames 6251999
[wof-iplookup] 00:20:49.187415 [debug] geonames ID (6251999) is WOF ID 85633041
[wof-iplookup] 00:20:49.187427 [debug] 142.213.160.134 is 85633041
85633041
```

### wof-iplookup-server

A handy HTTP pony for performing IP lookups as a service.

```
$> ./bin/wof-iplookup-server -mmdb GeoLite2-City.mmdb -concordances /usr/local/mapzen/whosonfirst-data/meta/wof-concordances-latest.csv
```

And then:

```
$> curl -s 'http://localhost:8668?ip=205.193.117.158' | python -mjson.tool
{
    "ip": "205.193.117.158",
    "wofid": 85784763
}
```

## Caveats

* The `wof-iplookup-server` does not return any data (like centroids, hierarchies or geometries) for a WOF record besides its ID. Yet.

## See also

* https://github.com/oschwald/maxminddb-golang
* https://dev.maxmind.com/geoip/geoip2/geolite2/
* https://github.com/whosonfirst/p5-Whosonfirst-MaxMind-Writer