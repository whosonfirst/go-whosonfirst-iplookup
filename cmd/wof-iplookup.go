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

	var mmdb = flag.String("mmdb", "", "The path to your mmdb file")
	var source = flag.String("source", "maxmind", "Who created this mmdb file?")
	var raw = flag.Bool("raw", false, "Dump the raw query response as JSON")
	var loglevel = flag.String("loglevel", "warning", "")

	flag.Parse()
	args := flag.Args()

	writer := io.MultiWriter(os.Stdout)

	logger := log.NewWOFLogger("[wof-iplookup] ")
	logger.AddLogger(writer, *loglevel)

	lookup, err := iplookup.NewIPLookup(*mmdb, *source, logger)

	if err != nil {
		logger.Error("failed to create IPLookup because %v", err)
		os.Exit(1)
	}

	for _, addr := range args {

		ip := net.ParseIP(addr)

		logger.Debug("lookup %s", addr)

		if *raw {

			rsp, err := lookup.QueryRaw(ip)

			if err != nil {
				logger.Error("failed to lookup %s, because %v", addr, err)
			}

			enc, _ := json.Marshal(rsp)
			fmt.Printf("%s\n", enc)
		} else {

			wofid, err := lookup.Query(ip)

			if err != nil {
				logger.Error("failed to lookup %s, because %v", addr, err)
			}

			fmt.Println(wofid)
		}
	}
}
