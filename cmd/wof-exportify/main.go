package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/sfomuseum/go-flags/multi"
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

	var ids multi.MultiInt64
	flag.Var(&ids, "id", "One or more Who's On First IDs. If left empty the value of the -i flag will be used.")

	flag.Usage = func() {

		fmt.Fprintf(os.Stderr, "Exportify one or more Who's On First IDs.\n\n")
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

	// please to be exposing reader.STDIN_SCHEME

	if *reader_uri == "stdin://" {

		err = exportStdin(ctx, r, wr, ex)

		if err != nil {
			log.Fatalf("Failed to export from STDIN, %v", err)
		}

		return
	}

	for _, id := range ids {

		err := exportId(ctx, r, wr, ex, id)

		if err != nil {
			log.Fatalf("Failed to export record for '%d', %v", id, err)
		}
	}
}

func exportStdin(ctx context.Context, r reader.Reader, wr writer.Writer, ex export.Exporter) error {

	fh, err := r.Read(ctx, "-")

	if err != nil {
		return err
	}

	body, err := io.ReadAll(fh)

	if err != nil {
		return err
	}

	return exportBytes(ctx, body, wr, ex)
}

func exportId(ctx context.Context, r reader.Reader, wr writer.Writer, ex export.Exporter, id int64) error {

	body, err := wof_reader.LoadBytes(ctx, r, id)

	if err != nil {
		return err
	}

	return exportBytes(ctx, body, wr, ex)
}

func exportBytes(ctx context.Context, body []byte, wr writer.Writer, ex export.Exporter) error {

	new_body, err := ex.Export(ctx, body)

	if err != nil {
		return err
	}

	_, err = wof_writer.WriteBytes(ctx, wr, new_body)

	if err != nil {
		return err
	}

	return nil
}
