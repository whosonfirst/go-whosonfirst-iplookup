package main

// TO DO: please reconcile with wof-geojsonls-dump.go

// wof-api -param api_key=mapzen-xxxxxx -param method=whosonfirst.places.getDescendants -param placetype=venue -param id=102086957 -geojson-ls -async -paginated --geojson-ls-output /usr/local/data-ext/lacity/wof-venues-lacounty.geojson.txt

// wof-geojsonls-dump-filelist -root /usr/local/data/whosonfirst-data-venue-us-ca/data /usr/local/data-ext/lacity/lacounty-venues.txt > /usr/local/data-ext/lacity/lacounty-venues-geojson.txt

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

func main() {

	outfile := flag.String("outfile", "", "Where to write records (default is STDOUT)")
	root := flag.String("root", "", "...")

	lieu := flag.Bool("lieu", false, "...")

	exclude_not_current := flag.Bool("exclude-not-current", false, "Exclude records that have been ...")
	exclude_ceased := flag.Bool("exclude-ceased", false, "Exclude records that have been ...")
	exclude_deprecated := flag.Bool("exclude-deprecated", false, "Exclude records that have been deprecated.")
	exclude_superseded := flag.Bool("exclude-superseded", false, "Exclude records that have been superseded.")

	procs := flag.Int("processes", runtime.NumCPU()*2, "The number of concurrent processes to use")

	flag.Parse()

	var wr *bufio.Writer

	if *outfile != "" {

		fh, err := os.Create(*outfile)

		if err != nil {
			log.Fatal(err)
		}

		wr = bufio.NewWriter(fh)

	} else {
		wr = bufio.NewWriter(os.Stdout)
	}

	mu := new(sync.Mutex)
	wg := new(sync.WaitGroup)

	throttle := make(chan bool, *procs)

	for i := 0; i < *procs; i++ {
		throttle <- true
	}

	for _, filelist := range flag.Args() {

		fh, err := os.Open(filelist)

		if err != nil {
			log.Fatal(err)
		}

		scanner := bufio.NewScanner(fh)

		for scanner.Scan() {

			rel_path := scanner.Text()
			abs_path := filepath.Join(*root, rel_path)

			is_wof, err := uri.IsWOFFile(abs_path)

			if err != nil {
				log.Fatal("unable to determine whether %s is a WOF file, because %s\n", abs_path, err)
			}

			if !is_wof {
				// log.Printf("%s is not a WOF file\n", abs_path)
				continue
			}

			is_alt, err := uri.IsAltFile(abs_path)

			if err != nil {
				log.Fatal("unable to determine whether %s is an alt (WOF) file, because %s\n", abs_path, err)
			}

			if is_alt {
				// log.Printf("%s is an alt (WOF) file\n", abs_path)
				continue
			}

			<-throttle

			wg.Add(1)

			// sudo put this functionality in a package funciton or something...

			go func(abs_abs_path string, wr *bufio.Writer, wg *sync.WaitGroup, throttle chan bool) {

				defer func() {
					wg.Done()
					throttle <- true
				}()

				fh, err := os.Open(abs_path)

				if err != nil {
					// log.Fatal("failed to open %s, because %s\n", abs_path, err)
					return
				}

				defer fh.Close()

				body, err := ioutil.ReadAll(fh)

				if err != nil {
					log.Fatal("failed to read %s, because %s\n", abs_path, err)
				}

				if *exclude_not_current {

					rsp := gjson.GetBytes(body, "properties.mz:is_current")

					if rsp.Exists() {

						is_current := rsp.Int()

						if is_current == 0 {
							return
						}
					}
				}

				if *exclude_ceased || *exclude_not_current {

					rsp := gjson.GetBytes(body, "properties.edtf:cessation")

					if rsp.Exists() {

						cessation := rsp.String()

						if cessation != "" && cessation != "uuuu" {
							return
						}
					}
				}

				if *exclude_deprecated || *exclude_not_current {

					rsp := gjson.GetBytes(body, "properties.edtf:deprecated")

					if rsp.Exists() {

						deprecated := rsp.String()

						if deprecated != "" && deprecated != "uuuu" {
							return
						}
					}
				}

				if *exclude_superseded  || *exclude_not_current {

					rsp := gjson.GetBytes(body, "properties.wof:superseded_by")

					if rsp.Exists() {

						superseded_by := rsp.Array()

						if len(superseded_by) > 0 {
							return
						}
					}
				}

				if *lieu {

					rsp := gjson.GetBytes(body, "properties.wof:id")

					if !rsp.Exists() {
						log.Fatal("WOF record is missing a wof:id property", abs_path)
					}

					source_id := fmt.Sprintf("wof:id=%d", rsp.Int())
					body, err = sjson.SetBytes(body, "id", source_id)

					if err != nil {
						log.Fatal("failed to set source ID for %s, because %s\n", abs_path, err)
					}

					name := gjson.GetBytes(body, "properties.wof:name")

					if !name.Exists() {
						log.Fatal("WOF record is missing a wof:name property", abs_path)
					}

					body, err = sjson.SetBytes(body, "properties.name", name.String())

					if err != nil {
						log.Fatal("failed to set name for %s, because %s\n", abs_path, err)
					}
				}

				var feature interface{}

				err = json.Unmarshal(body, &feature)

				if err != nil {
					log.Fatal("failed to parse %s, because %s\n", abs_path, err)
				}

				body, err = json.Marshal(feature)

				if err != nil {
					log.Fatal("failed to parse %s, because %s\n", abs_path, err)
				}

				mu.Lock()
				defer mu.Unlock()

				_, err = wr.Write(body)

				if err != nil {
					return
				}

				wr.Write([]byte("\n"))
				wr.Flush()

			}(abs_path, wr, wg, throttle)

			wg.Wait()
		}
	}
}
