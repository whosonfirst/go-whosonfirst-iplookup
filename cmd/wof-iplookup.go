package main

import (
	"flag"
	"fmt"
	iplookup "github.com/whosonfirst/go-whosonfirst-iplookup"
	"net"
)

func main() {

	var mmdb = flag.String("mmdb", "", "")
	var concordances = flag.String("concordances", "", "")

	flag.Parse()
	args := flag.Args()

	lookup, err := iplookup.NewLookup(*mmdb, *concordances)

	if err != nil {
	   panic(err)
	}

	for _, addr := range args {

		ip := net.ParseIP(addr)
		r, err := lookup.Query(ip)

		if err != nil {
		   panic(err)
		}

		fmt.Printf("%s is %d\n", addr, r)
	}
}
