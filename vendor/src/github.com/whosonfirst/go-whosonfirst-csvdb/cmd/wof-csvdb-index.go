package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-csvdb"
	"github.com/whosonfirst/go-whosonfirst-log"
	"io"
	"os"
	"strings"
	"time"
)

func main() {

	var cols = flag.String("columns", "", "Comma-separated list of columns to index")
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

	t1 := time.Now()

	db, err := csvdb.NewCSVDB(logger)

	if err != nil {
		fmt.Printf("failed to create csvdb\n")
		os.Exit(1)
	}

	for _, path := range args {

		err := db.IndexCSVFile(path, to_index)

		if err != nil {
			fmt.Printf("failed to index %s, because %v\n", path, err)
			os.Exit(1)
		}
	}

	t2 := time.Since(t1)

	fmt.Printf("> %v\n", t2)

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("> query <col>=<id>")
	fmt.Printf("> ")

	for scanner.Scan() {

		input := scanner.Text()
		query := strings.Split(input, "=")

		if len(query) != 2 {
			fmt.Println("invalid query")
			continue
		}

		k := query[0]
		v := query[1]

		fmt.Printf("search for %s=%s\n", k, v)

		t1 := time.Now()

		rows, _ := db.Where(k, v)

		t2 := time.Since(t1)

		fmt.Printf("where %s=%s %d results (%v)\n", k, v, len(rows), t2)

		fmt.Println("")
		fmt.Println("> query <col>=<id>")
		fmt.Printf("> ")
	}
}
