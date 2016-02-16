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
	"strings"
)

type Response struct {
	IPAddress string `json:"ip"`
	WOFId     int64  `json:"wofid"`
}

func main() {

	var db = flag.String("db", "", "The path to your IP lookup database file")
	var source = flag.String("source", "maxmind", "The source of the IP lookups")
	var host = flag.String("host", "localhost", "The hostname to listen for requests on")
	var port = flag.Int("port", 8668, "The port number to listen for requests on")
	var cors = flag.Bool("cors", false, "Enable CORS headers")
	var loglevel = flag.String("loglevel", "status", "")

	flag.Parse()

	writer := io.MultiWriter(os.Stdout)

	logger := log.NewWOFLogger("[wof-iplookup] ")
	logger.AddLogger(writer, *loglevel)

	lookup, err := iplookup.NewIPLookup(*db, *source, logger)

	if err != nil {
		logger.Error("failed to create IPLookup because %v", err)
		os.Exit(1)
	}

	handler := func(rsp http.ResponseWriter, req *http.Request) {

		query := req.URL.Query()
		ip := query.Get("ip")
		raw := query.Get("raw")

		if ip == "" {
			remote := strings.Split(req.RemoteAddr, ":")
			ip = remote[0]
		}

		if ip == "" {
			http.Error(rsp, "Missing IP address", http.StatusInternalServerError)
			return
		}

		if ip == "127.0.0.1" {
			http.Error(rsp, "We are all localhost", http.StatusInternalServerError)
			return
		}

		logger.Debug("parse IP %s", ip)

		addr := net.ParseIP(ip)
		var answer interface{}

		if raw != "" {

			qrsp, err := lookup.QueryRaw(addr)

			if err != nil {
				http.Error(rsp, err.Error(), http.StatusInternalServerError)
				return
			}

			answer = qrsp

		} else {
			wofid, err := lookup.QueryId(addr)

			if err != nil {
				http.Error(rsp, err.Error(), http.StatusInternalServerError)
				return
			}

			answer = Response{ip, wofid}
		}

		js, err := json.Marshal(answer)

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

	logger.Status("wof-iplookup-server running at %s", endpoint)

	http.HandleFunc("/", handler)
	http.ListenAndServe(endpoint, nil)
}
