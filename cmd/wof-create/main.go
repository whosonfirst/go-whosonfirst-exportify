package main

// go run -mod vendor cmd/create/main.go -writer-uri stdout:// -geometry '{"type":"Point", "coordinates":[20.414944,42.032833]}' -string-property 'properties.src:geom=wikidata'

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/paulmach/orb/geojson"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/whosonfirst/go-reader"
	export "github.com/whosonfirst/go-whosonfirst-export/v2"
	exportify "github.com/whosonfirst/go-whosonfirst-exportify"
	wof_reader "github.com/whosonfirst/go-whosonfirst-reader"
	hierarchy "github.com/whosonfirst/go-whosonfirst-spatial-hierarchy"
	_ "github.com/whosonfirst/go-whosonfirst-spatial-sqlite"
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
	"github.com/whosonfirst/go-whosonfirst-spatial/filter"
	wof_writer "github.com/whosonfirst/go-whosonfirst-writer"
	"github.com/whosonfirst/go-writer"
)

//go:embed stub.geojson
var stub []byte

func main() {

	fs := flagset.NewFlagSet("create")

	source := fs.String("s", "", "A valid path to the root directory of the Who's On First data repository. If empty (and -reader-uri or -writer-uri are empty) the current working directory will be used and appended with a 'data' subdirectory.")

	parent_reader_uri := fs.String("parent-reader-uri", "", "A valid whosonfirst/go-reader URI. If empty the value of the -s fs will be used in combination with the fs:// scheme.")
	writer_uri := fs.String("writer-uri", "", "A valid whosonfirst/go-writer URI. If empty the value of the -s fs will be used in combination with the fs:// scheme.")

	exporter_uri := fs.String("exporter-uri", "whosonfirst://", "A valid whosonfirst/go-whosonfirst-export URI.")

	spatial_database_uri := fs.String("spatial-database-uri", "", "A valid whosonfirst/go-whosonfirst-spatial/database URI.")

	var str_properties multi.KeyValueString
	fs.Var(&str_properties, "string-property", "One or more {KEY}={VALUE} flags where {KEY} is a valid tidwall/gjson path and {VALUE} is a string value.")

	var int_properties multi.KeyValueInt64
	fs.Var(&int_properties, "int-property", "One or more {KEY}={VALUE} flags where {KEY} is a valid tidwall/gjson path and {VALUE} is a int(64) value.")

	var float_properties multi.KeyValueFloat64
	fs.Var(&float_properties, "float-property", "One or more {KEY}={VALUE} flags where {KEY} is a valid tidwall/gjson path and {VALUE} is a float(64) value.")

	str_geom := fs.String("geometry", "", "A valid GeoJSON geometry")

	resolve_hierarchy := fs.Bool("resolve-hierarchy", false, "Attempt to resolve parent ID and hierarchy using point-in-polygon lookups. If true the -spatial-database-uri flag must also be set")

	fs.Usage = func() {

		fmt.Fprintf(os.Stderr, "Create a new Who's On First record.\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options] \n\n", os.Args[0])
		// fmt.Fprintf(os.Stderr, "For example:\n")
		// fmt.Fprintf(os.Stderr, "\t%s ...\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n")
		fs.PrintDefaults()
	}

	flagset.Parse(fs)

	if *parent_reader_uri == "" || *writer_uri == "" {

		var path string

		if *source == "" {

			cwd, err := os.Getwd()

			if err != nil {
				log.Fatalf("Failed to determine current working directory, %v", err)
			}

			path = cwd

		} else {

			abs_source, err := filepath.Abs(*source)

			if err != nil {
				log.Fatalf("Failed to derive absolute path for '%s', %v", *source, err)
			}

			path = abs_source
		}

		abs_path := fmt.Sprintf("fs://%s/data", path)

		if *parent_reader_uri == "" {
			*parent_reader_uri = abs_path
		}

		if *writer_uri == "" {
			*writer_uri = abs_path
		}
	}

	ctx := context.Background()

	ex, err := export.NewExporter(ctx, *exporter_uri)

	if err != nil {
		log.Fatalf("Failed create exporter for '%s', %v", *exporter_uri, err)
	}

	wr, err := writer.NewWriter(ctx, *writer_uri)

	if err != nil {
		log.Fatalf("Failed to create writer for '%s', %v", *writer_uri, err)
	}

	geom, err := geojson.UnmarshalGeometry([]byte(*str_geom))

	if err != nil {
		log.Fatalf("Failed to unmarshal geometry, %v", err)
	}

	opts := &exportify.UpdateFeatureOptions{
		StringProperties:  str_properties,
		Int64Properties:   int_properties,
		Float64Properties: float_properties,
		Geometry:          geom,
	}

	body, _, err := exportify.UpdateFeature(ctx, stub, opts)

	if err != nil {
		log.Fatalf("Failed to update properties, %v", err)
	}

	parent_rsp := gjson.GetBytes(body, "properties.wof:parent_id")
	parent_id := parent_rsp.Int()

	if parent_id > 0 {

		parent_r, err := reader.NewReader(ctx, *parent_reader_uri)

		if err != nil {
			log.Fatalf("Failed to create reader for '%s', %v", *parent_reader_uri, err)
		}

		parent_body, err := wof_reader.LoadBytesFromID(ctx, parent_r, parent_id)

		if err != nil {
			log.Fatalf("Failed to load parent record (%d), %v", parent_id, err)
		}

		to_copy := []string{
			"properties.wof:hierarchy",
			"properties.wof:country",
		}

		for _, path := range to_copy {

			rsp := gjson.GetBytes(parent_body, path)
			body, err = sjson.SetBytes(body, path, rsp.Value())

			if err != nil {
				log.Fatalf("Failed to copy '%s' parent value, %v", path, err)
			}
		}
	}

	// START OF pip/hierarchy stuff

	if *resolve_hierarchy {

		spatial_db, err := database.NewSpatialDatabase(ctx, *spatial_database_uri)

		if err != nil {
			log.Fatalf("Failed to create new spatial database for '%s', %v", *spatial_database_uri, err)
		}

		resolver, err := hierarchy.NewPointInPolygonHierarchyResolver(ctx, spatial_db, nil)

		if err != nil {
			log.Fatalf("Failed to create new hierarchy resolver, %v", err)
		}

		inputs := &filter.SPRInputs{}

		results_cb := hierarchy.FirstButForgivingSPRResultsFunc
		update_cb := hierarchy.DefaultPointInPolygonHierarchyResolverUpdateCallback()

		new_body, err := resolver.PointInPolygonAndUpdate(ctx, inputs, results_cb, update_cb, body)

		if err != nil {
			log.Fatalf("Failed to do point in polygon operation, %v", err)
		}

		body = new_body
	}

	// END OF pip/hierarchy stuff

	new_body, err := ex.Export(ctx, body)

	if err != nil {
		log.Fatalf("Failed to export new record, %v", err)
	}

	id_rsp := gjson.GetBytes(new_body, "properties.wof:id")

	err = wof_writer.WriteFeatureBytes(ctx, wr, new_body)

	if err != nil {
		log.Fatalf("Failed to write new record, %v", err)
	}

	fmt.Printf("%d\n", id_rsp.Int())
}
