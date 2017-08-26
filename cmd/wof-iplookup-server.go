package main

import (
	"flag"
	"fmt"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/whosonfirst/go-whosonfirst-iplookup"
	"github.com/whosonfirst/go-whosonfirst-iplookup/http"
	"github.com/whosonfirst/go-whosonfirst-iplookup/provider"
	"github.com/whosonfirst/go-whosonfirst-log"
	"io"
	gohttp "net/http"
	"os"
)

func main() {

	var db = flag.String("db", "", "The path to your IP lookup database file")
	var host = flag.String("host", "localhost", "The hostname to listen for requests on")
	var port = flag.Int("port", 8668, "The port number to listen for requests on")
	var loglevel = flag.String("loglevel", "status", "")

	flag.Parse()

	writer := io.MultiWriter(os.Stdout)

	logger := log.NewWOFLogger()
	logger.AddLogger(writer, *loglevel)

	pr, err := provider.NewMMDBProvider(*db, iplookup.SPRRecordToResult)

	if err != nil {
		logger.Error("failed to create IPLookup because %v", err)
		os.Exit(1)
	}

	lookuphandler, err := http.LookupHandler(pr)

	if err != nil {
		logger.Fatal("failed to create Lookup handler because %s", err)
	}

	pinghandler, err := http.PingHandler()

	if err != nil {
		logger.Fatal("failed to create Ping handler because %s", err)
	}

	endpoint := fmt.Sprintf("%s:%d", *host, *port)
	logger.Status("wof-iplookup-server running at %s", endpoint)

	mux := gohttp.NewServeMux()
	mux.Handle("/", lookuphandler)
	mux.Handle("/ping", pinghandler)

	err = gracehttp.Serve(&gohttp.Server{Addr: endpoint, Handler: mux})

	if err != nil {
		logger.Fatal("failed to start server because %s", err)
	}

	os.Exit(0)
}
