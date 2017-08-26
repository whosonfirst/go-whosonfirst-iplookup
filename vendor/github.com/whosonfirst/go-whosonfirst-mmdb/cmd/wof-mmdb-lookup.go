package main

import (
	"encoding/json"
	"flag"
	"github.com/whosonfirst/go-whosonfirst-csv"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/properties/whosonfirst"
	"github.com/whosonfirst/go-whosonfirst-log"
	"github.com/whosonfirst/go-whosonfirst-mmdb"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

func main() {

	var concordances = flag.String("concordances", "", "")
	var data_root = flag.String("data-root", "/usr/local/data", "")
	var repo = flag.String("repo", "whosonfirst-data", "")

	flag.Parse()

	logger := log.SimpleWOFLogger()

	fh, err := os.Open(*concordances)

	if err != nil {
		logger.Fatal("failed to open %s, because %s", *concordances, err)
	}

	lookup := make(map[int64][]*mmdb.SPRRecord)

	root := filepath.Join(*data_root, *repo)
	data := filepath.Join(root, "data")

	reader, err := csv.NewDictReader(fh)

	if err != nil {
		logger.Fatal("failed to open csv reader because %s", err)
	}

	for {
		row, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			logger.Fatal("failed to process %s because %s", *concordances, err)
		}

		str_gnid, ok := row["gn:id"]

		if !ok {
			logger.Fatal("missing gn:id key")
		}

		gnid, err := strconv.ParseInt(str_gnid, 10, 64)

		if err != nil {
			logger.Fatal("failed to parse %s because %s", str_gnid, err)
		}

		str_wofid, ok := row["wof:id"]

		if !ok {
			logger.Fatal("missing wof:id key")
		}

		if str_wofid == "-1" {
			continue
		}

		wofid, err := strconv.ParseInt(str_wofid, 10, 64)

		if err != nil {
			logger.Fatal("failed to parse %s because %s", str_wofid, err)
		}

		abs_path, err := uri.Id2AbsPath(data, wofid)

		if err != nil {
			logger.Fatal("failed to determine absolute path for %d because %s", wofid, err)
		}

		f, err := feature.LoadWOFFeatureFromFile(abs_path)

		if err != nil {
			logger.Fatal("failed to load %s because %s", abs_path, err)
		}

		to_process := make(map[int64]geojson.Feature)
		to_process[wofid] = f

		/*
		hiers := whosonfirst.Hierarchy(f)

		for _, hier := range hiers {

			for _, id := range hier {

				_, ok := to_process[id]

				if ok {
					continue
				}

				abs_path, err := uri.Id2AbsPath(data, id)

				if err != nil {
					logger.Fatal("failed to determine absolute path for %d because %s", wofid, err)
				}

				f, err := feature.LoadWOFFeatureFromFile(abs_path)

				if err != nil {
					logger.Fatal("failed to load %s because %s", abs_path, err)
				}

				to_process[id] = f
			}
		}
		*/

		mm_records := make([]*mmdb.SPRRecord, 0)

		for id, f := range to_process {

			mm, err := FeatureToSPRRecord(f)

			if err != nil {
				logger.Fatal("failed to create MM record for %d because %s", id, err)
			}

			mm_records = append(mm_records, mm)
		}

		lookup[gnid] = mm_records
	}

	enc, err := json.Marshal(lookup)	

	if err != nil {
		logger.Fatal("failed to marshal lookup because %s", err)
	}

	writer := os.Stdout
	writer.Write(enc)
}

func FeatureToSPRRecord(f geojson.Feature) (*mmdb.SPRRecord, error) {

	id, _ := strconv.ParseInt(f.Id(), 10, 64)
	name := f.Name()
	pt := f.Placetype()

	centroid, err := whosonfirst.Centroid(f)

	if err != nil {
		return nil, err
	}

	bboxes, err := f.BoundingBoxes()

	if err != nil {
		return nil, err
	}

	coord := centroid.Coord()
	mbr := bboxes.MBR()
	
	mm := mmdb.SPRRecord{
		Id:           id,
		Name:         name,
		Placetype:    pt,
		Latitude:     coord.Y,
		Longitude:    coord.X,
		MinLatitude:  mbr.Min.Y,
		MinLongitude: mbr.Min.X,
		MaxLatitude:  mbr.Max.Y,
		MaxLongitude: mbr.Max.X,
	}

	return &mm, nil
}
