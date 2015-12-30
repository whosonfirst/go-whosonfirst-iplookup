# go-whosonfirst-iplookup

Go package for doing IP address to Who's On First record lookups

## Caveats

This is not quite ready for you to use yet.

## Usage

```
./bin/wof-iplookup -mmdb GeoLite2-City.mmdb -concordances /usr/local/mapzen/whosonfirst-data/meta/wof-concordances-latest.csv -loglevel debug 8.8.8.8
[wof-iplookup] 21:39:56.340480 [debug] time to index concordances 2.99936672s
[wof-iplookup] 21:39:56.340592 [debug] possible matches for 8.8.8.8: [5375480 6252001]
[wof-iplookup] 21:39:56.340602 [debug] concordify geonames 5375480
[wof-iplookup] 21:39:56.340645 [debug] geonames ID (5375480) is WOF ID 85922355
8.8.8.8 is 85922355
```

## See also

* https://github.com/oschwald/maxminddb-golang
* https://dev.maxmind.com/geoip/geoip2/geolite2/
* https://github.com/whosonfirst/p5-Whosonfirst-MaxMind-Writer