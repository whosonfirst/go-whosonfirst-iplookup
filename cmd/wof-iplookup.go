package main

import (
	"encoding/json"
	"flag"
	"github.com/whosonfirst/go-whosonfirst-iplookup"
	"github.com/whosonfirst/go-whosonfirst-log"
	"os"
)

func main() {

	var db = flag.String("db", "", "The path to your IP lookup database file")

	flag.Parse()

	logger := log.SimpleWOFLogger()

	lookup, err := iplookup.NewIPLookup(*db, iplookup.SPRRecordToReponse)

	if err != nil {
		logger.Fatal("failed to create IPLookup because %v", err)
	}

	for _, ip := range flag.Args() {

		r, err := lookup.QueryString(ip)

		if err != nil {
			logger.Fatal("unable to query %s because %s", ip, err)
		}

		enc, err := json.Marshal(r)

		if err != nil {
			logger.Fatal("unable to encode %s because %s", ip, err)
		}

		logger.Status("%s becomes %s", ip, string(enc))
	}

	os.Exit(0)
}
