package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/sfomuseum/go-csvdict"
	"github.com/sfomuseum/go-flags/multi"
	"github.com/whosonfirst/go-reader"
	export "github.com/whosonfirst/go-whosonfirst-export/v2"
	wof_reader "github.com/whosonfirst/go-whosonfirst-reader"
	wof_writer "github.com/whosonfirst/go-whosonfirst-writer/v3"
	"github.com/whosonfirst/go-writer/v3"
)

func main() {

	reader_uri := flag.String("reader-uri", "", "A valid whosonfirst/go-reader URI")
	writer_uri := flag.String("writer-uri", "", "A valid whosonfirst/go-writer URI")

	lookup_key := flag.String("lookup-key", "wof:id", "The column in a CSV row to use to lookup a corresponding Who's On First record.")

	var str_fields multi.MultiString
	flag.Var(&str_fields, "string-field", "Zero or more fields in a CSV row to assign to a WOF record as string values.")

	var int_fields multi.MultiString
	flag.Var(&int_fields, "int-field", "Zero or more fields in a CSV row to assign to a WOF record as int values.")

	var int64_fields multi.MultiString
	flag.Var(&int64_fields, "int64-field", "Zero or more fields in a CSV row to assign to a WOF record as int64 values.")
	
	exporter_uri := flag.String("exporter-uri", "whosonfirst://", "A valid whosonfirst/go-whosonfirst-export URI")

	flag.Usage = func() {

		fmt.Fprintf(os.Stderr, "Upate one or more Who's On First records with matching entries in a CSV file.\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options] path(N) path(N)\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "For example:\n")
		fmt.Fprintf(os.Stderr, "\t%s -reader-uri repo:///usr/local/data/sfomuseum-data-architecture -writer-uri repo:///usr/local/data/sfomuseum-data-architecture -int-field sfo:level galleries-with-level.csv\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	ctx := context.Background()

	ex, err := export.NewExporter(ctx, *exporter_uri)

	if err != nil {
		log.Fatal(err)
	}

	r, err := reader.NewReader(ctx, *reader_uri)

	if err != nil {
		log.Fatalf("Failed to create new reader for '%s', %v", *reader_uri, err)
	}

	wr, err := writer.NewWriter(ctx, *writer_uri)

	if err != nil {
		log.Fatalf("Failed to create new writer for '%s', %v", *writer_uri, err)
	}

	paths := flag.Args()

	for _, path := range paths {

		csv_r, err := csvdict.NewReaderFromPath(path)

		if err != nil {
			log.Fatalf("Failed to create CSV reader for %s, %v", path, err)
		}

		for {

			row, err := csv_r.Read()

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("Failed to read row, %v", err)
			}

			str_id, exists := row[*lookup_key]

			if !exists {
				log.Fatalf("Row missing '%s' key", *lookup_key)
			}

			wof_id, err := strconv.ParseInt(str_id, 10, 64)

			if err != nil {
				log.Fatalf("Failed to parse '%s' as WOF Id, %v", str_id, err)
			}

			body, err := wof_reader.LoadBytes(ctx, r, wof_id)

			if err != nil {
				log.Printf("Failed to load '%d', %v. Skipping", wof_id, err)
				continue
			}

			updates := make(map[string]interface{})

			for _, field := range str_fields {

				str_v, exists := row[field]

				if !exists {
					log.Printf("Missing '%s' key in CSV for '%d', skipping", field, wof_id)
					continue
				}

				path := fmt.Sprintf("properties.%s", field)
				updates[path] = str_v
			}

			for _, field := range int_fields {

				str_v, exists := row[field]

				if !exists {
					log.Printf("Missing '%s' key in CSV for '%d', skipping", field, wof_id)
					continue
				}

				v, err := strconv.Atoi(str_v)

				if err != nil {
					log.Fatalf("Failed to convert string value for '%s' (%s) to int, %v", field, str_v, err)
				}

				path := fmt.Sprintf("properties.%s", field)
				updates[path] = v
			}

			for _, field := range int64_fields {

				str_v, exists := row[field]

				if !exists {
					log.Printf("Missing '%s' key in CSV for '%d', skipping", field, wof_id)
					continue
				}

				v, err := strconv.ParseInt(str_v, 10, 64)

				if err != nil {
					log.Fatalf("Failed to convert string value for '%s' (%s) to int64, %v", field, str_v, err)
				}

				path := fmt.Sprintf("properties.%s", field)
				updates[path] = v
			}

			has_changed, new_body, err := export.AssignPropertiesIfChanged(ctx, body, updates)

			if err != nil {
				log.Fatalf("Failed to assign new properties for %d, %v", wof_id, err)
			}

			if !has_changed {
				continue
			}

			new_body, err = ex.Export(ctx, new_body)

			if err != nil {
				log.Fatalf("Failed to export new record for '%d', %v", wof_id, err)
			}

			_, err = wof_writer.WriteBytes(ctx, wr, new_body)

			if err != nil {
				log.Fatalf("Failed to write record for '%d', %v", wof_id, err)
			}

		}

	}
}
