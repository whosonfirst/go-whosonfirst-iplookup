package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-iplookup"
	"github.com/whosonfirst/go-whosonfirst-log"
	"io"
	"net"
	"net/http"
	"os"
)

func main() {

	var mmdb = flag.String("mmdb", "", "")
	var concordances = flag.String("concordances", "", "")
	var host = flag.String("host", "localhost", "The hostname to listen for requests on")
	var port = flag.Int("port", 8668, "The port number to listen for requests on")
	var cors = flag.Bool("cors", false, "Enable CORS headers")
	var loglevel = flag.String("loglevel", "warning", "")

	flag.Parse()

	writer := io.MultiWriter(os.Stdout)

	logger := log.NewWOFLogger("[wof-iplookup] ")
	logger.AddLogger(writer, *loglevel)

	lookup, err := iplookup.NewIPLookup(*mmdb, *concordances, logger)

	if err != nil {
		logger.Error("failed to create IPLookup because %v", err)
		os.Exit(1)
	}

	handler := func(rsp http.ResponseWriter, req *http.Request) {

		query := req.URL.Query()
		ip := query.Get("ip")

		// TO DO: chunk out port numbers etc.

		if ip == "" {
			ip = req.RemoteAddr
		}

		logger.Debug("parse IP %s", ip)

		addr := net.ParseIP(ip)
		wofid, err := lookup.Query(addr)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		js, err := json.Marshal(wofid)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		// maybe this although it seems like it adds functionality for a lot of
		// features this server does not need - https://github.com/rs/cors
		// (20151022/thisisaaronland)

		if *cors {
			rsp.Header().Set("Access-Control-Allow-Origin", "*")
		}

		rsp.Header().Set("Content-Type", "application/json")
		rsp.Write(js)
	}

	endpoint := fmt.Sprintf("%s:%d", *host, *port)

	logger.Info("wof-iplookup-server running at %s", endpoint)

	http.HandleFunc("/", handler)
	http.ListenAndServe(endpoint, nil)
}
