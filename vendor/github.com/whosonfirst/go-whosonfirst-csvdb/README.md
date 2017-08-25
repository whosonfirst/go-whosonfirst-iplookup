# go-whosonfirst-csvdb

Experimental in-memory database for CSV files.

## Caveats

This is not sophisticated. It is not meant to be sophisticated. It is meant to be easy and fast. It might also be too soon for you to play with depending on how you feel about "things in flux".

Things this package does not do yet (or maybe ever) in no particular order:

* Complex queries
* Query operators besides testing equality (`=`)
* Pagination - this package is still design primarily for datasets that will not return a lot of results for a given query
* Proper logging
* Serializing indexes to disk (or loading them)

## Utilities

### wof-csvdb-index

This is a little bit of a misnomer as it's mostly a testing tool right now. Oh well...

In this example we'll index three columns from the [wof-concordances-latest.csv](https://github.com/whosonfirst-data/whosonfirst-data/blob/master/meta/wof-concordances-latest.csv) file (specifically `wof:id` and `gp:id` and `gn:id`) and then perform a couple queries against the index:

```
./bin/wof-csvdb-index -columns wof:id,gp:id,gn:id /usr/local/mapzen/whosonfirst-data/meta/wof-concordances-latest.csv
> indexes: 3 keys: 239979 rows: 573437 time to index: 2.572705702s
> query <col>=<id>
> gp:id=3534
search for gp:id=3534
where gp:id=3534 1 results (6.789µs)

> query <col>=<id>
> gp:id=55992163
search for gp:id=55992163
where gp:id=55992163 2 results (5.035µs)

> query <col>=<id>
> wd:id=Q30
search for wd:id=Q30
where wd:id=Q30 0 results (2.698µs)

> query <col>=<id>
> gn:id=6252001
search for gn:id=6252001
where gn:id=6252001 1 results (7.443µs)
```

You can also use `go-whosonfirst-csvdb` with non- Who's On First datasets. For example we can index the `objects.csv` meta file that it part of the [collection metadata repository for the Cooper Hewitt, Smithsonian Design Museum](https://github.com/cooperhewitt/collection):

```
./bin/wof-csvdb-index -columns woe:country_id,decade ../../cooperhewitt/objects.csv
> indexes: 2 keys: 134423 rows: 118 time to index: 5.684592543s
> query <col>=<id>
> decade=1920
search for decade=1920
where decade=1920 5092 results (220.538µs)

> query <col>=<id>
> decade=1980
search for decade=1980
where decade=1980 2668 results (135.227µs)

> query <col>=<id>
> type_id=35268079
search for type_id=35268079
where type_id=35268079 0 results (3.107µs)

> query <col>=<id>
> woe:country_id=23424977
search for woe:country_id=23424977
where woe:country_id=23424977 44193 results (2.986562ms)
```

### wof-csvdb-server

A small HTTP pony for querying a CSV file and getting the results back as JSON.

```
$> ./bin/wof-csvdb-server -columns wof:id,gp:id,gn:id /usr/local/mapzen/whosonfirst-data/meta/wof-concordances-latest.csv
time to index /usr/local/mapzen/whosonfirst-data/meta/wof-concordances-latest.csv: 2.667267366s
wof-csvdb-server running at localhost:8228
```

And then:

```
curl -s 'http://localhost:8228?k=gp:id&v=3534' | python -mjson.tool
[
    {
        "dbp:id": "Montreal",
        "fb:id": "en.montreal",
        "fct:id": "03c06bce-8f76-11e1-848f-cfd5bf3ef515",
        "gn:id": "6077243",
        "gp:id": "3534",
        "nyt:id": "N59179828586486930801",
        "tgn:id": "7013051",
        "wd:id": "Q340",
        "wof:id": "101736545"
    }
]
```

_Note that as of this writing the `wof-csvdb-server` does not offer any kind of introspection so you need to know what has been indexed before you issue a query._

## See also

* https://github.com/whosonfirst/go-whosonfirst-csv
