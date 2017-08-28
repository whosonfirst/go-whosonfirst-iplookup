# go-whosonfirst-iplookup

Go package for doing IP address lookups with Who's On First "standard place results"

## Install

You will need to have both `Go` (specifically a version of Go more recent than 1.6 so let's just assume you need [Go 1.8](https://golang.org/dl/) or higher) and the `make` programs installed on your computer. Assuming you do just type:

```
make bin
```

All of this package's dependencies are bundled with the code in the `vendor` directory.

## Interfaces

### Provider

```
type Provider interface {
	QueryString(string) (spr.StandardPlacesResult, error)
	Query(net.IP) (spr.StandardPlacesResult, error)
}
```

## Packages

### http

```
import (
	"github.com/whosonfirst/go-whosonfirst-iplookup/http"
	"github.com/whosonfirst/go-whosonfirst-mmdb/provider"
	gohttp "net/http"
)

func main() {

	pr, _ := provider.NewWOFProvider("example.mmdb")

	lookuphandler, _ := http.LookupHandler(pr)

	gohttp.HandleFunc("/", lookuphandler)
	gohttp.ListenAndServe(":8080", nil)
}
```

## See also

* https://github.com/whosonfirst/go-whosonfirst-mmdb
* https://github.com/whosonfirst/p5-Whosonfirst-MaxMind-Writer
* https://github.com/oschwald/maxminddb-golang
* https://dev.maxmind.com/geoip/geoip2/geolite2/
