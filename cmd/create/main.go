package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"github.com/sfomuseum/go-flags/multi"
	// "github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-export/v2"
	"github.com/whosonfirst/go-whosonfirst-exportify"
	// wof_reader "github.com/whosonfirst/go-whosonfirst-reader"
	wof_writer "github.com/whosonfirst/go-whosonfirst-writer"
	"github.com/whosonfirst/go-writer"
	"log"
	"os"
	"path/filepath"
)

//go:embed stub.geojson
var stub []byte

func main() {

	source := flag.String("s", "", "A valid path to the root directory of the Who's On First data repository. If empty (and -reader-uri or -writer-uri are empty) the current working directory will be used and appended with a 'data' subdirectory.")

	reader_uri := flag.String("reader-uri", "", "A valid whosonfirst/go-reader URI. If empty the value of the -s flag will be used in combination with the fs:// scheme.")
	writer_uri := flag.String("writer-uri", "", "A valid whosonfirst/go-writer URI. If empty the value of the -s flag will be used in combination with the fs:// scheme.")

	exporter_uri := flag.String("exporter-uri", "whosonfirst://", "A valid whosonfirst/go-whosonfirst-export URI.")

	var str_properties multi.KeyValueString
	flag.Var(&str_properties, "string-property", "One or more {KEY}={VALUE} flags where {KEY} is a valid tidwall/gjson path and {VALUE} is a string value.")

	var int_properties multi.KeyValueInt64
	flag.Var(&int_properties, "int-property", "One or more {KEY}={VALUE} flags where {KEY} is a valid tidwall/gjson path and {VALUE} is a int(64) value.")

	var float_properties multi.KeyValueFloat64
	flag.Var(&float_properties, "float-property", "One or more {KEY}={VALUE} flags where {KEY} is a valid tidwall/gjson path and {VALUE} is a float(64) value.")

	flag.Usage = func() {

		fmt.Fprintf(os.Stderr, "Create a new Who's On First record.\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options] \n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "For example:\n")
		fmt.Fprintf(os.Stderr, "\t%s ...\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

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

	/*
		r, err := reader.NewReader(ctx, *reader_uri)

		if err != nil {
			log.Fatalf("Failed to create reader for '%s', %v", *reader_uri, err)
		}
	*/

	wr, err := writer.NewWriter(ctx, *writer_uri)

	if err != nil {
		log.Fatalf("Failed to create writer for '%s', %v", *writer_uri, err)
	}

	opts := &exportify.UpdateFeatureOptions{
		StringProperties:  str_properties,
		Int64Properties:   int_properties,
		Float64Properties: float_properties,
	}

	body, _, err := exportify.UpdateFeature(ctx, stub, opts)

	if err != nil {
		log.Fatalf("Failed to update properties, %v", err)
	}

	new_body, err := ex.Export(ctx, body)

	if err != nil {
		log.Fatalf("Failed to export new record, %v", err)
	}

	err = wof_writer.WriteFeatureBytes(ctx, wr, new_body)

	if err != nil {
		log.Fatalf("Failed to write new record, %v", err)
	}

	log.Println("OK")
}
