package http

import (
	"encoding/json"
	"github.com/whosonfirst/go-whosonfirst-iplookup"
	"net"
	gohttp "net/http"
	"strings"
)

func LookupHandler(pr iplookup.Provider) (gohttp.Handler, error) {

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		query := req.URL.Query()
		ip := query.Get("ip")

		if ip == "" {

			remote := strings.Split(req.RemoteAddr, ":")
			ip = remote[0]
		}

		if ip == "" {
			gohttp.Error(rsp, "Missing IP address", gohttp.StatusInternalServerError)
			return
		}

		if ip == "127.0.0.1" {
			gohttp.Error(rsp, "We are all localhost", gohttp.StatusInternalServerError)
			return
		}

		addr := net.ParseIP(ip)

		r, err := pr.Query(addr)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		js, err := json.Marshal(r)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		rsp.Header().Set("Access-Control-Allow-Origin", "*")
		rsp.Header().Set("Content-Type", "application/json")

		rsp.Write(js)
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
