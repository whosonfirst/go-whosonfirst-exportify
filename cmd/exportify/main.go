package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/sfomuseum/go-flags/multi"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-export/v2"
	wof_reader "github.com/whosonfirst/go-whosonfirst-reader"
	wof_writer "github.com/whosonfirst/go-whosonfirst-writer"
	"github.com/whosonfirst/go-writer"
	"log"
	"path/filepath"
)

func main() {

	source := flag.String("s", "", "...")
	id := flag.String("i", "", "...")

	reader_uri := flag.String("reader-uri", "", "...")
	writer_uri := flag.String("writer-uri", "", "...")

	exporter_uri := flag.String("exporter-uri", "whosonfirst://", "...")

	var ids multi.MultiInt64
	flag.Var(&ids, "id", "...")

	flag.Parse()

	if *source != "" {

		abs_source, err := filepath.Abs(*source)

		if err != nil {
			log.Fatalf("Failed to derive absolute path for '%s', %v", *source, err)
		}

		*reader_uri = fmt.Sprintf("fs://%s/data", abs_source)
		*writer_uri = *reader_uri
	}

	if *id != "" {

		err := ids.Set(*id)

		if err != nil {
			log.Fatalf("Failed to assign '%d' (-id) flag, %v", *id, err)
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

	for _, id := range ids {

		err := exportId(ctx, r, wr, ex, id)

		if err != nil {
			log.Fatalf("Failed to export record for '%d', %v", id, err)
		}
	}
}

func exportId(ctx context.Context, r reader.Reader, wr writer.Writer, ex export.Exporter, id int64) error {

	body, err := wof_reader.LoadBytesFromID(ctx, r, id)

	if err != nil {
		return err
	}

	new_body, err := ex.Export(ctx, body)

	if err != nil {
		return err
	}

	err = wof_writer.WriteFeatureBytes(ctx, wr, new_body)

	if err != nil {
		return err
	}

	return nil
}
