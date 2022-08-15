package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/sfomuseum/go-flags/multi"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/whosonfirst/go-reader"
	export "github.com/whosonfirst/go-whosonfirst-export/v2"
	wof_reader "github.com/whosonfirst/go-whosonfirst-reader"
	wof_writer "github.com/whosonfirst/go-whosonfirst-writer/v2"
	"github.com/whosonfirst/go-writer/v2"
)

func main() {

	source := flag.String("s", "", "A valid path to the root directory of the Who's On First data repository. If empty (and -reader-uri or -writer-uri are empty) the current working directory will be used and appended with a 'data' subdirectory.")
	id := flag.String("i", "", "A valid Who's On First ID.")

	reader_uri := flag.String("reader-uri", "", "A valid whosonfirst/go-reader URI. If empty the value of the -s flag will be used in combination with the fs:// scheme.")
	writer_uri := flag.String("writer-uri", "", "A valid whosonfirst/go-writer URI. If empty the value of the -s flag will be used in combination with the fs:// scheme.")

	exporter_uri := flag.String("exporter-uri", "whosonfirst://", "A valid whosonfirst/go-whosonfirst-export URI.")

	var str_properties multi.KeyValueString
	flag.Var(&str_properties, "string-property", "One or more {KEY}={VALUE} properties to append to the new record where {KEY} is a valid tidwall/gjson path and {VALUE} is a string value.")

	var int_properties multi.KeyValueInt64
	flag.Var(&str_properties, "int-property", "One or more {KEY}={VALUE} properties to append to the new record where {KEY} is a valid tidwall/gjson path and {VALUE} is a int(64) value.")

	var float_properties multi.KeyValueFloat64
	flag.Var(&str_properties, "float-property", "One or more {KEY}={VALUE} properties to append to the new record where {KEY} is a valid tidwall/gjson path and {VALUE} is a float(64) value.")

	var ids multi.MultiInt64
	flag.Var(&ids, "id", "One or more Who's On First IDs. If left empty the value of the -i flag will be used.")

	flag.Usage = func() {

		fmt.Fprintf(os.Stderr, "Deprecate and supersede one or more Who's On First IDs.\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options] wof-id-(N) wof-id-(N)\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "For example:\n")
		fmt.Fprintf(os.Stderr, "\t%s -s . -i 1234\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\t%s -reader-uri fs:///usr/local/data/whosonfirst-data-admin-ca/data -id 1234 -id 5678\n\n", os.Args[0])
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

	if *id != "" {

		err := ids.Set(*id)

		if err != nil {
			log.Fatalf("Failed to assign '%s' (-id) flag, %v", *id, err)
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

	props := make([]multi.KeyValueFlag, 0)

	for _, p := range str_properties {
		props = append(props, p)
	}

	for _, p := range int_properties {
		props = append(props, p)
	}

	for _, p := range float_properties {
		props = append(props, p)
	}

	for _, old_id := range ids {

		new_id, err := replaceId(ctx, r, wr, ex, old_id, props...)

		if err != nil {
			log.Fatalf("Failed to export record for '%d', %v", old_id, err)
		}

		log.Printf("%d replaced by %d\n", old_id, new_id)
	}
}

func replaceId(ctx context.Context, r reader.Reader, wr writer.Writer, ex export.Exporter, old_id int64, props ...multi.KeyValueFlag) (int64, error) {

	old_body, err := wof_reader.LoadBytes(ctx, r, old_id)

	if err != nil {
		return -1, err
	}

	// Create the new record

	new_body := old_body

	new_body, err = sjson.DeleteBytes(new_body, "properties.wof:id")

	if err != nil {
		return -1, err
	}

	new_body, err = ex.Export(ctx, new_body)

	if err != nil {
		return -1, err
	}

	id_rsp := gjson.GetBytes(new_body, "properties.wof:id")

	if !id_rsp.Exists() {
		return -1, fmt.Errorf("failed to derive new properties.wof:id property for record superseding '%d'", old_id)
	}

	new_id := id_rsp.Int()

	// Update the new record

	new_body, err = sjson.SetBytes(new_body, "properties.wof:supsersedes", []int64{old_id})

	if err != nil {
		return -1, err
	}

	for _, p := range props {

		path := p.Key()
		new_value := p.Value()

		new_body, err = sjson.SetBytes(new_body, path, new_value)

		if err != nil {
			return -1, err
		}
	}

	// Update the old record

	now := time.Now()
	deprecated := now.Format("2006-01-02")

	old_body, err = sjson.SetBytes(old_body, "properties.edtf:deprecated", deprecated)

	if err != nil {
		return -1, err
	}

	old_body, err = sjson.SetBytes(old_body, "properties.mz:is_current", 0)

	if err != nil {
		return -1, err
	}

	old_body, err = sjson.SetBytes(old_body, "properties.wof:supserseded_by", []int64{new_id})

	if err != nil {
		return -1, err
	}

	// Write records

	old_body, err = ex.Export(ctx, old_body)

	if err != nil {
		return -1, err
	}

	new_body, err = ex.Export(ctx, new_body)

	if err != nil {
		return -1, err
	}

	_, err = wof_writer.WriteBytes(ctx, wr, old_body)

	if err != nil {
		return -1, err
	}

	_, err = wof_writer.WriteBytes(ctx, wr, new_body)

	if err != nil {
		return -1, err
	}

	return new_id, nil
}
