package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-iplookup"
	"github.com/whosonfirst/go-whosonfirst-log"
	"io"
	"net"
	"os"
)

func main() {

	var db = flag.String("db", "", "The path to your IP lookup database file")
	var source = flag.String("source", "maxmind", "The source of the IP lookups")
	var raw = flag.Bool("raw", false, "Return the raw data")
	var asjson = flag.Bool("json", false, "Dump the raw query response as JSON")
	var loglevel = flag.String("loglevel", "warning", "")

	flag.Parse()
	args := flag.Args()

	writer := io.MultiWriter(os.Stdout)

	logger := log.NewWOFLogger("[wof-iplookup] ")
	logger.AddLogger(writer, *loglevel)

	lookup, err := iplookup.NewIPLookup(*db, *source, logger)

	if err != nil {
		logger.Error("failed to create IPLookup because %v", err)
		os.Exit(1)
	}

	for _, addr := range args {

		ip := net.ParseIP(addr)
		logger.Debug("lookup %s", addr)

		var result interface{}

		if *raw {

			rsp, err := lookup.QueryRaw(ip)

			if err != nil {
				logger.Error("failed to lookup %s, because %v", addr, err)
			}

			result = rsp

		} else {

			rsp, err := lookup.Query(ip)

			if err != nil {
				logger.Error("failed to lookup %s, because %v", addr, err)
			}

			result = rsp
		}

		if *asjson {
			enc, _ := json.Marshal(result)
			fmt.Printf("%s\n", enc)
		} else {
			fmt.Println(result)
		}
	}
}
