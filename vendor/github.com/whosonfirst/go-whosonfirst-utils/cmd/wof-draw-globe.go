package main

import (
	"errors"
	"flag"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/globe" // for to make DrawPreparedPaths public
	"github.com/whosonfirst/go-whosonfirst-crawl"
	"github.com/whosonfirst/go-whosonfirst-csv"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func DrawFeature(feature []byte, gl *globe.Globe) error {

	geom_type := gjson.GetBytes(feature, "geometry.type")

	if !geom_type.Exists() {
		return errors.New("Geometry is missing a type property")
	}

	coords := gjson.GetBytes(feature, "geometry.coordinates")

	if !coords.Exists() {
		return errors.New("Geometry is missing a coordinates property")
	}

	switch geom_type.String() {

	case "Point":

		lonlat := coords.Array()
		lat := lonlat[1].Float()
		lon := lonlat[0].Float()

		gl.DrawDot(lat, lon, 0.01)

	// http://geojson.org/geojson-spec.html#id4

	case "Polygon":

		paths := make([][]*globe.Point, 0)

		for _, ring := range coords.Array() {

			path := make([]*globe.Point, 0)

			for _, r := range ring.Array() {

				lonlat := r.Array()
				lat := lonlat[1].Float()
				lon := lonlat[0].Float()

				pt := globe.NewPoint(lat, lon)
				path = append(path, &pt)
			}

			paths = append(paths, path)
		}

		gl.DrawPaths(paths)

	// http://geojson.org/geojson-spec.html#id7

	case "MultiPolygon":

		for _, polys := range coords.Array() {

			paths := make([][]*globe.Point, 0)

			for _, ring := range polys.Array() {

				path := make([]*globe.Point, 0)

				for _, r := range ring.Array() {

					lonlat := r.Array()
					lat := lonlat[1].Float()
					lon := lonlat[0].Float()

					pt := globe.NewPoint(lat, lon)
					path = append(path, &pt)
				}

				paths = append(paths, path)
			}

			gl.DrawPaths(paths)
		}

	default:
		return errors.New("Unsupported geometry type")
	}

	return nil
}

func DrawRow(path string, row map[string]string, g *globe.Globe, throttle chan bool, remote bool) error {

	<-throttle

	defer func() {
		throttle <- true
	}()

	rel_path, ok := row["path"]

	if !ok {
		log.Println("Missing path")
		return nil
	}

	var feature []byte

	if remote {

		root := "https://whosonfirst.mapzen.com/data"
		uri := filepath.Join(root, rel_path)

		rsp, err := http.Get(uri)

		if err != nil {
			log.Fatal(err)
		}

		defer rsp.Body.Close()

		feature, err = ioutil.ReadAll(rsp.Body)

		if err != nil {
			log.Fatal("failed to read %s, because %s\n", uri, err)
		}

	} else {

		meta := filepath.Dir(path)
		root := filepath.Dir(meta)
		data := filepath.Join(root, "data")

		abs_path := filepath.Join(data, rel_path)

		fh, err := os.Open(abs_path)

		if err != nil {
			log.Fatal("failed to open %s, because %s\n", abs_path, err)
		}

		defer fh.Close()

		feature, err = ioutil.ReadAll(fh)

		if err != nil {
			log.Fatal("failed to read %s, because %s\n", abs_path, err)
		}
	}

	return DrawFeature(feature, g)
}

func main() {

	outfile := flag.String("out", "", "Where to write globe")
	size := flag.Int("size", 1600, "The size of the globe (in pixels)")
	mode := flag.String("mode", "meta", "... (default is 'meta' for one or more meta files)")

	remote := flag.Bool("remote", false, "...")
	feature := flag.Bool("feature", false, "...")
	rotate := flag.Bool("rotate", false, "...")

	center := flag.String("center", "", "")
	center_lat := flag.Float64("latitude", 37.755244, "")
	center_lon := flag.Float64("longitude", -122.447777, "")

	flag.Parse()

	if *center != "" {

		latlon := strings.Split(*center, ",")

		lat, err := strconv.ParseFloat(latlon[0], 64)

		if err != nil {
			log.Fatal(err)
		}

		lon, err := strconv.ParseFloat(latlon[1], 64)

		if err != nil {
			log.Fatal(err)
		}

		*center_lat = lat
		*center_lon = lon
	}

	green := color.NRGBA{0x00, 0x64, 0x3c, 192}
	g := globe.New()
	g.DrawGraticule(10.0)

	t1 := time.Now()

	if *mode == "meta" {

		max_fh := 10
		throttle := make(chan bool, max_fh)

		for i := 0; i < max_fh; i++ {
			throttle <- true
		}

		for _, path := range flag.Args() {

			reader, err := csv.NewDictReaderFromPath(path)

			if err != nil {
				log.Fatal(err)
			}

			for {
				row, err := reader.Read()

				if err == io.EOF {
					break
				}

				if err != nil {
					log.Println(err, path)
					break
				}

				if *feature {
					DrawRow(path, row, g, throttle, *remote)
					continue
				}

				str_lat, ok := row["geom_latitude"]

				if !ok {
					continue
				}

				str_lon, ok := row["geom_longitude"]

				if !ok {
					continue
				}

				lat, err := strconv.ParseFloat(str_lat, 64)

				if err != nil {
					log.Println(err, str_lat)
					continue
				}

				lon, err := strconv.ParseFloat(str_lon, 64)

				if err != nil {
					log.Println(err, str_lon)
					continue
				}

				g.DrawDot(lat, lon, 0.01, globe.Color(green))
			}
		}

	} else if *mode == "repo" {

		for _, path := range flag.Args() {

			var cb = func(path string, info os.FileInfo) error {

				if info.IsDir() {
					return nil
				}

				is_wof, err := uri.IsWOFFile(path)

				if err != nil {
					log.Printf("unable to determine whether %s is a WOF file, because %s\n", path, err)
					return err
				}

				if !is_wof {
					return nil
				}

				is_alt, err := uri.IsAltFile(path)

				if err != nil {
					log.Printf("unable to determine whether %s is an alt (WOF) file, because %s\n", path, err)
					return err
				}

				if is_alt {
					return nil
				}

				fh, err := os.Open(path)

				if err != nil {
					log.Printf("failed to open %s, because %s\n", path, err)
					return err
				}

				defer fh.Close()

				feature, err := ioutil.ReadAll(fh)

				if err != nil {
					log.Printf("failed to read %s, because %s\n", path, err)
					return err
				}

				return DrawFeature(feature, g)
			}

			cr := crawl.NewCrawler(path)
			cr.Crawl(cb)
		}

	} else {

		log.Fatal("Invalid mode")
	}

	t2 := time.Since(t1)

	log.Printf("time to read all the things %v\n", t2)

	t3 := time.Now()

	if *rotate {

		// Initialize palette (#ffffff, #000000, #ff0000)
		var palette color.Palette = color.Palette{}
		palette = append(palette, color.White)
		palette = append(palette, color.Black)
		palette = append(palette, color.RGBA{0xff, 0x00, 0x00, 0xff})

		var images []*image.Paletted
		var delays []int

		coords := [][]float64{
			[]float64{45.0, 0.0},
			[]float64{45.0, 45.0},
			[]float64{45.0, 90.0},
			[]float64{45.0, 135.0},
			[]float64{45.0, 180.0},
			[]float64{45.0, -135.0},
			[]float64{45.0, -90.0},
			[]float64{45.0, -45.0},
		}

		ch := make(chan *image.Paletted)

		max_proc := 1 // apparently anything will invoke the OOM killer...
		throttle := make(chan bool, max_proc)

		for i := 0; i < max_proc; i++ {
			throttle <- true
		}

		for _, latlon := range coords {

			lat := latlon[0]
			lon := latlon[1]

			go func(lat float64, lon float64, throttle chan bool) {

				defer func() {
					throttle <- true
				}()

				t1 := time.Now()

				g.CenterOn(lat, lon)
				im := g.Image(*size)

				t2 := time.Since(t1)
				log.Printf("time to render %v\n", t2)

				pm := image.NewPaletted(im.Bounds(), palette)
				draw.FloydSteinberg.Draw(pm, im.Bounds(), im, image.ZP)

				// images = append(images, pm)
				// delays = append(delays, 200)

				ch <- pm

			}(lat, lon, throttle)
		}

		count := len(coords)

		for i := count; i > 0; {

			select {
			case pm := <-ch:

				images = append(images, pm)
				delays = append(delays, 20)

				i -= 1

				log.Print("count is ", i)
			default:
				// pass
			}
		}

		fh, err := os.OpenFile(*outfile, os.O_WRONLY|os.O_CREATE, 0600)

		if err != nil {
			log.Fatal(err)
		}

		defer fh.Close()

		gif.EncodeAll(fh, &gif.GIF{
			Image: images,
			Delay: delays,
		})

	} else {

		g.CenterOn(*center_lat, *center_lon)

		err := g.SavePNG(*outfile, *size)

		if err != nil {
			log.Fatal(err)
		}

	}

	t4 := time.Since(t3)

	log.Printf("time to draw all the things %v\n", t4)

}
