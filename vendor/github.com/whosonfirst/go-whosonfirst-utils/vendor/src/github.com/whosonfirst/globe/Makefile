geodata: land.geodata.go countries.geodata.go

world.topojson:
	wget -O $@ https://unpkg.com/world-atlas@1.1.4/world/110m.json

%.world.geojson: world.topojson
	topo2geo --in $< $*=$@

%.geodata.go: %.world.geojson buildgeodata.go
	go run buildgeodata.go -input $< -output $@ -var $*
	gofmt -s -w $@

tools:
	npm install -g topojson
