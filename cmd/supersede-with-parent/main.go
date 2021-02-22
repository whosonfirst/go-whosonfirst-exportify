package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/sfomuseum/go-flags/multi"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-export/v2"
	"github.com/whosonfirst/go-whosonfirst-exportify"
	wof_reader "github.com/whosonfirst/go-whosonfirst-reader"
	"github.com/whosonfirst/go-writer"
	"log"
	"os"
)

func main() {

	reader_uri := flag.String("reader-uri", "", "A valid whosonfirst/go-reader URI.")
	writer_uri := flag.String("writer-uri", "", "A valid whosonfirst/go-writer URI. If empty the value of the -reader-uri flag will be assumed.")
	parent_reader_uri := flag.String("parent-reader-uri", "", "A valid whosonfirst/go-reader URI. If empty the value of the -reader-uri flag will be assumed.")

	exporter_uri := flag.String("exporter-uri", "whosonfirst://", "A valid whosonfirst/go-whosonfirst-export URI.")

	var ids multi.MultiInt64
	flag.Var(&ids, "id", "One or more valid Who's On First ID.")

	parent_id := flag.Int64("parent-id", 0, "A valid Who's On First ID.")

	flag.Usage = func() {

		fmt.Fprintf(os.Stderr, "Supersede one or more WOF records with a known parent ID (and hierarchy)\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "For example:\n")
		fmt.Fprintf(os.Stderr, "\t%s -reader-uri fs:///usr/local/data/sfomuseum-data-architecture/data -parent-id 1477855937 -id 1477855955\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nValid options are:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *parent_reader_uri == "" {
		*parent_reader_uri = *reader_uri
	}

	if *writer_uri == "" {
		*writer_uri = *reader_uri
	}

	ctx := context.Background()

	r, err := reader.NewReader(ctx, *reader_uri)

	if err != nil {
		log.Fatalf("Failed to create reader for '%s', %v", *reader_uri, err)
	}

	parent_r, err := reader.NewReader(ctx, *parent_reader_uri)

	if err != nil {
		log.Fatalf("Failed to create reader for '%s', %v", *parent_reader_uri, err)
	}

	wr, err := writer.NewWriter(ctx, *writer_uri)

	if err != nil {
		log.Fatalf("Failed to create new writer for '%s', %v", *writer_uri, err)
	}

	ex, err := export.NewExporter(ctx, *exporter_uri)

	if err != nil {
		log.Fatalf("Failed to create new exporter for '%s', %v", *exporter_uri, err)
	}

	// Parent stuff we only need to set up once

	parent_f, err := wof_reader.LoadBytesFromID(ctx, parent_r, *parent_id)

	if err != nil {
		log.Fatalf("Failed to load '%d', %v", *parent_id, err)
	}

	hier_rsp := gjson.GetBytes(parent_f, "properties.wof:hierarchy")

	if !hier_rsp.Exists() {
		log.Fatalf("Parent (%d) is missing properties.wof:hierarchy", *parent_id)
	}

	parent_hierarchy := hier_rsp.Value()

	inception_rsp := gjson.GetBytes(parent_f, "properties.edtf:inception")

	if !inception_rsp.Exists() {
		log.Fatalf("Parent (%d) is missing properties.edtf:inception", *parent_id)
	}

	cessation_rsp := gjson.GetBytes(parent_f, "properties.edtf:cessation")

	if !cessation_rsp.Exists() {
		log.Fatalf("Parent (%d) is missing properties.edtf:cessation", *parent_id)
	}

	inception := inception_rsp.String()
	cessation := cessation_rsp.String()

	to_update_old := map[string]interface{}{
		"properties.edtf:cessation": inception,
	}

	to_update_new := map[string]interface{}{
		"properties.wof:parent_id":  *parent_id,
		"properties.wof:hierarchy":  parent_hierarchy,
		"properties.edtf:inception": inception,
		"properties.edtf:cessation": cessation,
	}

	// Okay, go

	for _, id := range ids {

		f, err := wof_reader.LoadBytesFromID(ctx, r, id)

		if err != nil {
			log.Fatalf("Failed to load '%d', %v", id, err)
		}

		old_f, new_f, err := export.SupersedeRecord(ctx, ex, f)

		if err != nil {
			log.Fatalf("Failed to supersede record %d, %v", id, err)
		}

		old_f, err = export.AssignProperties(ctx, old_f, to_update_old)

		if err != nil {
			log.Fatalf("Failed to assign properties for new record, ")
		}

		name_rsp := gjson.GetBytes(new_f, "properties.wof:name")

		if !name_rsp.Exists() {
			log.Fatalf("Failed to retrieve wof:name for new record")
		}

		name := name_rsp.String()
		label := fmt.Sprintf("%s (%s)", name, inception)

		to_update_new["properties.wof:label"] = label

		new_f, err = export.AssignProperties(ctx, new_f, to_update_new)

		if err != nil {
			log.Fatalf("Failed to assign updated properties for new record, ")
		}

		err = exportify.ExportWithWriter(ctx, ex, wr, old_f)

		if err != nil {
			log.Fatalf("Failed to export '%d', %v", id, err)
		}

		err = exportify.ExportWithWriter(ctx, ex, wr, new_f)

		if err != nil {
			log.Fatalf("Failed to export new feature, %v", id, err)
		}
	}

}
