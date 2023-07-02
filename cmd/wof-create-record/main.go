package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/whosonfirst/go-reader"
	export "github.com/whosonfirst/go-whosonfirst-export/v2"
	wofReader "github.com/whosonfirst/go-whosonfirst-reader"
	wofWriter "github.com/whosonfirst/go-whosonfirst-writer/v3"
	"github.com/whosonfirst/go-writer/v3"
)

func main() {
	fs := flagset.NewFlagSet("create")

	parentReaderURI := fs.String("parent-reader-uri", "", "A valid whosonfirst/go-reader URI")
	parentID := fs.Int64("parent-wof-id", -1, "An optional WOF ID which the created record should be parented by")
	writerURI := fs.String("writer-uri", "", "A valid whosonfirst/go-writer URI")
	exporterURI := fs.String("exporter-uri", "whosonfirst://", "A valid whosonfirst/go-whosonfirst-export URI.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Create a new WOF record from a partially prepared record. Useful when you have a new record, but need to give it an ID, place it into the hierarchy and write it into a repo.\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options] file [file ...]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n")
		fs.PrintDefaults()
	}

	flagset.Parse(fs)

	// Fail if there's no files
	if fs.NArg() < 1 {
		fs.Usage()
		log.Fatalf("No files provided")
	}

	ctx := context.Background()

	// Setup the exporter
	ex, err := export.NewExporter(ctx, *exporterURI)
	if err != nil {
		log.Fatalf("Failed create exporter for '%s', %v", *exporterURI, err)
	}

	// Setup the writer
	wr, err := writer.NewWriter(ctx, *writerURI)
	if err != nil {
		log.Fatalf("Failed create writer for '%s', %v", *writerURI, err)
	}

	files := fs.Args()

	for _, file := range files {
		log.Print(file)

		bytes, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("Unable to read file: %s, %v", file, err)
		}

		// If -parent-id is provided, attempt to set the hierarchy and other details from the parent
		if *parentID > 0 {
			if *parentReaderURI == "" {
				log.Fatalf("No parent reader URI provided")
			}

			parentReader, err := reader.NewReader(ctx, *parentReaderURI)
			if err != nil {
				log.Fatalf("Failed to create reader for '%s', %v", *parentReaderURI, err)
			}

			parentBytes, err := wofReader.LoadBytes(ctx, parentReader, *parentID)
			if err != nil {
				log.Fatalf("Failed to load parent record (%d), %v", *parentID, err)
			}

			to_copy := []string{
				"properties.wof:hierarchy",
				"properties.wof:country",
			}

			for _, path := range to_copy {
				rsp := gjson.GetBytes(parentBytes, path)
				bytes, err = sjson.SetBytes(bytes, path, rsp.Value())

				if err != nil {
					log.Fatalf("Failed to copy '%s' parent value, %v", path, err)
				}
			}
		}

		exportBytes, err := ex.Export(ctx, bytes)
		if err != nil {
			log.Fatalf("Failed to export '%s', %v", file, err)
		}

		_, err = wofWriter.WriteBytes(ctx, wr, exportBytes)
		if err != nil {
			log.Fatalf("Failed to write '%s', %v", file, err)
		}
	}

}
