package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-csvdb"
	"github.com/whosonfirst/go-whosonfirst-log"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {

	var cols = flag.String("columns", "", "Comma-separated list of columns to index")
	var host = flag.String("host", "localhost", "The hostname to listen for requests on")
	var port = flag.Int("port", 8228, "The port number to listen for requests on")
	var cors = flag.Bool("cors", false, "Enable CORS headers")
	var loglevel = flag.String("loglevel", "info", "Log level for reporting")

	flag.Parse()
	args := flag.Args()

	to_index := make([]string, 0)

	for _, c := range strings.Split(*cols, ",") {
		to_index = append(to_index, c)
	}

	l_writer := io.MultiWriter(os.Stdout)

	logger := log.NewWOFLogger("[wof-csvdb-index] ")
	logger.AddLogger(l_writer, *loglevel)

	db, err := csvdb.NewCSVDB(logger)

	if err != nil {
		panic(err)
	}

	for _, path := range args {

		t1 := time.Now()

		err := db.IndexCSVFile(path, to_index)

		if err != nil {

			if err.Error() == "EOF" {
				fmt.Printf("skip %s because it appears to be empty...\n", path)
			} else {
				msg := fmt.Sprintf("failed to %s, because %v", path, err)
				err = errors.New(msg)

				panic(err)
			}
		}

		t2 := time.Since(t1)
		fmt.Printf("time to index %s: %v\n", path, t2)
	}

	handler := func(rsp http.ResponseWriter, req *http.Request) {

		if db.Indexing() {
			http.Error(rsp, "Database is re-indexing, please try again shortly", http.StatusServiceUnavailable)
			return
		}

		query := req.URL.Query()

		k := query.Get("k")
		v := query.Get("v")

		if k == "" {
			http.Error(rsp, "Missing k parameter", http.StatusBadRequest)
			return
		}

		if v == "" {
			http.Error(rsp, "Missing v parameter", http.StatusBadRequest)
			return
		}

		rows, err := db.Where(k, v)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		results := make([]map[string]string, 0)

		for _, row := range rows {
			results = append(results, row.AsMap())
		}

		js, err := json.Marshal(results)

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

	fmt.Printf("wof-csvdb-server running at %s\n", endpoint)

	http.HandleFunc("/", handler)
	err = http.ListenAndServe(endpoint, nil)

	if err != nil {
		logger.Error("failed to start wof-csvdb-server because %v", err)
		os.Exit(1)
	}

	os.Exit(0)
}
