# go-whosonfirst-mmdb

Go package for working with Who's On First documents and MaxMind DB files.

## Install

You will need to have both `Go` (specifically a version of Go more recent than 1.6 so let's just assume you need [Go 1.8](https://golang.org/dl/) or higher) and the `make` programs installed on your computer. Assuming you do just type:

```
make bin
```

All of this package's dependencies are bundled with the code in the `vendor` directory.

## Important

Too soon. Move along. Really, it's too soon.

## Tools

## wof-mmdb-lookup

```
 ./bin/wof-mmdb-lookup -concordances /usr/local/data-ext/maxmind-data/201708/wof-geonames.csv > lookup.json
```

This tool is designed to take a CSV file mapping Geonames ID to Who's On First ID and produce a JSON file that can be consumed by the [p5-Whosonfirst-MaxMind-Writer](https://github.com/whosonfirst/p5-Whosonfirst-MaxMind-Writer) library. For example:

```
$> less /usr/local/data-ext/maxmind-data/201708/wof-geonames.csv
gn:id,wof:id,mm:city,mm:country,wof:possible
49518,85632303,,RW,
51537,85632379,,SO,
69543,85632499,,YE,
99237,85632191,,IQ,
102358,85632253,,SA,
130758,85632361,,IR,
146669,85632437,,CY,
149590,85632227,,TZ,
163843,85632413,,SY,
174982,85632773,,AM,
192950,85632329,,KE,
...and so on
```

Which produces this:

```
$> less lookup.json
{"102358":[{"wof:id":85632253,"wof:name":"Saudi Arabia","wof:placetype":"country","wof:latitude":23.806678,"wof:longitude":44.700847,"geom:min_latitude":16.370945,"geom:min_longitude":34.572765,"geom:max_latitude":32.121348,"geom:max_longitude":55.637565},{"wof:id":102191569,"wof:name":"Asia","wof:placetype":"continent","wof:latitude":49.512481,"wof:longitude":94.464337,"geom:min_latitude":-12.199965,"geom:min_longitude":-180,"geom:max_latitude":81.288804,"geom:max_longitude":180}],"1036973":[{"wof:id":85632729,"wof:name":"Mozambique","wof:placetype":"country","wof:latitude":-13.885531,"wof:longitude":37.837456,"geom:min_latitude":-26.86816149999993,"geom:min_longitude":30.21555,"geom:max_latitude":-10.47719478599993,"geom:max_longitude":40.84875106800007},{"wof:id":102191573,"wof:name":"Africa","wof:placetype":"continent","wof:latitude":21.638471 ... and so on
```

In these examples the `wof-geonames.csv` file is being produced by the `wof-mmdb-build-concordances` script in the [py-mapzen-whosonfirst-maxmind
](https://github.com/whosonfirst/py-mapzen-whosonfirst-maxmind) but that functionality will probably be moved in to this package.

_Note: As of this writing the `p5-Whosonfirst-MaxMind-Writer` package has not been updated to consume files like `lookup.json`._

## See also

* https://github.com/whosonfirst/p5-Whosonfirst-MaxMind-Writer
* https://github.com/whosonfirst/py-mapzen-whosonfirst-maxmind