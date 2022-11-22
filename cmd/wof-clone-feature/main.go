package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/sfomuseum/go-flags/multi"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/whosonfirst/go-reader"
	export "github.com/whosonfirst/go-whosonfirst-export/v2"
	wof_reader "github.com/whosonfirst/go-whosonfirst-reader"
	wof_writer "github.com/whosonfirst/go-whosonfirst-writer/v3"
	"github.com/whosonfirst/go-writer/v3"
	"log"
	"os"
	"path/filepath"
)

func main() {

	source := flag.String("s", "", "A valid path to the root directory of the Who's On First data repository. If empty (and -reader-uri or -writer-uri are empty) the current working directory will be used and appended with a 'data' subdirectory.")

	reader_uri := flag.String("reader-uri", "", "A valid whosonfirst/go-reader URI. If empty the value of the -s flag will be used in combination with the fs:// scheme.")
	writer_uri := flag.String("writer-uri", "", "A valid whosonfirst/go-writer URI. If empty the value of the -s flag will be used in combination with the fs:// scheme.")

	exporter_uri := flag.String("exporter-uri", "whosonfirst://", "A valid whosonfirst/go-whosonfirst-export URI.")

	var str_properties multi.KeyValueString
	flag.Var(&str_properties, "string-property", "One or more {KEY}={VALUE} properties to append to the new record where {KEY} is a valid tidwall/gjson path and {VALUE} is a string value.")

	var int_properties multi.KeyValueInt64
	flag.Var(&str_properties, "int-property", "One or more {KEY}={VALUE} properties to append to the new record where {KEY} is a valid tidwall/gjson path and {VALUE} is a int(64) value.")

	var float_properties multi.KeyValueFloat64
	flag.Var(&str_properties, "float-property", "One or more {KEY}={VALUE} properties to append to the new record where {KEY} is a valid tidwall/gjson path and {VALUE} is a float(64) value.")

	src_id := flag.Int64("id", 0, "The feature being cloned.")

	supersedes := flag.Bool("supersedes", false, "The new feature supersedes the feature being cloned.")
	superseded := flag.Bool("superseded", false, "The new feature is superseded by the feature being cloned.")

	flag.Usage = func() {

		fmt.Fprintf(os.Stderr, "Clone and optionally supersede a Who's On First record.\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "For example:\n")
		fmt.Fprintf(os.Stderr, "\t%s -s . -id 1234 -superseded\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\t%s -reader-uri fs:///usr/local/data/whosonfirst-data-admin-ca/data -id 1234 -superseded\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *supersedes && *superseded {
		log.Fatalf("New record can both supersede and be superseded")
	}

	if *reader_uri == "" || *writer_uri == "" {

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

		if *reader_uri == "" {
			*reader_uri = abs_path
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

	r, err := reader.NewReader(ctx, *reader_uri)

	if err != nil {
		log.Fatalf("Failed to create reader for '%s', %v", *reader_uri, err)
	}

	wr, err := writer.NewWriter(ctx, *writer_uri)

	if err != nil {
		log.Fatalf("Failed to create writer for '%s', %v", *writer_uri, err)
	}

	// Load the record being cloned

	src_body, err := wof_reader.LoadBytes(ctx, r, *src_id)

	if err != nil {
		log.Fatalf("Failed to load record, %v", err)
	}

	// Create the new record

	new_body := src_body

	new_body, err = sjson.DeleteBytes(new_body, "properties.wof:id")

	if err != nil {
		log.Fatalf("Failed to remove wof:id from new record, %v", err)
	}

	new_updates := make(map[string]interface{})

	for _, p := range str_properties {
		new_updates[p.Key()] = p.Value()
	}

	for _, p := range int_properties {
		new_updates[p.Key()] = p.Value()
	}

	for _, p := range float_properties {
		new_updates[p.Key()] = p.Value()
	}

	if *supersedes {
		new_updates["properties.wof:supersedes"] = []int64{*src_id}
	}

	if *superseded {
		new_updates["properties.wof:superseded_by"] = []int64{*src_id}
		new_updates["properties.mz:is_current"] = 0
	}

	new_body, err = export.AssignProperties(ctx, new_body, new_updates)

	if err != nil {
		log.Fatalf("Failed to assign properties to new record, %v", err)
	}

	new_body, err = ex.Export(ctx, new_body)

	if err != nil {
		log.Fatalf("Failed to export new record, %v", err)
	}

	id_rsp := gjson.GetBytes(new_body, "properties.wof:id")

	if !id_rsp.Exists() {
		log.Fatalf("failed to derive new properties.wof:id property for record superseding '%d'", src_id)
	}

	new_id := id_rsp.Int()

	new_body, err = ex.Export(ctx, new_body)

	if err != nil {
		log.Fatalf("Failed to export new record, %v", err)
	}

	_, err = wof_writer.WriteBytes(ctx, wr, new_body)

	if err != nil {
		log.Fatalf("Failed to write new record, %v", err)
	}

	// Update the old record

	src_updates := make(map[string]interface{})

	if *supersedes {
		src_updates["properties.wof:superseded_by"] = []int64{new_id}
		new_updates["properties.mz:is_current"] = 0
	}

	if *superseded {
		src_updates["properties.wof:supersedes"] = []int64{new_id}
	}

	has_changed, src_body, err := export.AssignPropertiesIfChanged(ctx, src_body, src_updates)

	if err != nil {
		log.Fatalf("Failed to update source properties, %v", err)
	}

	if has_changed {

		src_body, err = ex.Export(ctx, src_body)

		if err != nil {
			log.Fatalf("Failed to export updated source record, %v", err)

		}

		_, err = wof_writer.WriteBytes(ctx, wr, src_body)

		if err != nil {
			log.Fatalf("Failed to write updated source record, %v", err)
		}

	}

	log.Printf("Created new record %d\n", new_id)
}
