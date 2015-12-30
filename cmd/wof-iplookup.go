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

	lookup, _ := iplookup.NewLookup(*mmdb, *concordances)

	for _, addr := range args {

		ip := net.ParseIP(addr)
		r, _ := lookup.Query(ip)

		fmt.Printf("%s is %d\n", addr, r)
	}
}
