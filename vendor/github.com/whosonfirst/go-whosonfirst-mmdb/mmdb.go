package mmdb

// is this really an SPR? not really... that's what we're trying to
// figure out... (20170824/thisisaaronland)

type SPRRecord struct {
	Id           int64   `json:"wof:id"`
	Name         string  `json:"wof:name"`
	Placetype    string  `json:"wof:placetype"`
	Latitude     float64 `json:"wof:latitude"`
	Longitude    float64 `json:"wof:longitude"`
	MinLatitude  float64 `json:"geom:min_latitude"`
	MinLongitude float64 `json:"geom:min_longitude"`
	MaxLatitude  float64 `json:"geom:max_latitude"`
	MaxLongitude float64 `json:"geom:max_longitude"`
}
