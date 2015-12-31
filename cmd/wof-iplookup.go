package main

import (
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-iplookup"
	"github.com/whosonfirst/go-whosonfirst-log"
	"io"
	"net"
	"os"
)

func main() {

	var mmdb = flag.String("mmdb", "", "")
	var concordances = flag.String("concordances", "", "")
	var loglevel = flag.String("loglevel", "warning", "")

	flag.Parse()
	args := flag.Args()

	writer := io.MultiWriter(os.Stdout)

	logger := log.NewWOFLogger("[wof-iplookup] ")
	logger.AddLogger(writer, *loglevel)

	lookup, err := iplookup.NewIPLookup(*mmdb, *concordances, logger)

	if err != nil {
		logger.Error("failed to create IPLookup because %v", err)
		os.Exit(1)
	}

	for _, addr := range args {

		ip := net.ParseIP(addr)
		wofid, err := lookup.Query(ip)

		if err != nil {
			logger.Error("failed to lookup %s, because %v", addr, err)
		}

		logger.Debug("%s is %d", addr, wofid)
		fmt.Println(wofid)
	}
}
