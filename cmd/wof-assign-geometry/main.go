// Assign the geometry from a given record to one or more other records.
package main

// This could (should) be extended to be a basic assign <properties or geometry> tool.

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/whosonfirst/go-reader"
	export "github.com/whosonfirst/go-whosonfirst-export/v2"
	wof_reader "github.com/whosonfirst/go-whosonfirst-reader"
	wof_writer "github.com/whosonfirst/go-whosonfirst-writer"
	"github.com/whosonfirst/go-writer/v2"
)

func main() {

	reader_uri := flag.String("reader-uri", "", "A valid whosonfirst/go-reader URI.")
	writer_uri := flag.String("writer-uri", "stdout://", "A valid whosonfirst/go-writer URI.")

	exporter_uri := flag.String("exporter-uri", "whosonfirst://", "A valid whosonfirst/go-whosonfirst-export URI")

	source_id := flag.Int64("source-id", 0, "A valid Who's On First ID.")

	from_stdin := flag.Bool("stdin", false, "Read target IDs from STDIN")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Assign the geometry from a given record to one or more other records.\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options] target-id-(N) target-id-(N)\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "For example:\n")
		fmt.Fprintf(os.Stderr, "\t%s -reader-uri fs:///usr/local/data/whosonfirst-data-admin-ca/data -source-id 1234 5678\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	// START OF put me in a function...

	str_ids := flag.Args()

	if *from_stdin {

		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			str_ids = append(str_ids, scanner.Text())
		}

		err := scanner.Err()

		if err != nil {
			log.Fatalf("Failed to read input from STDIN, %v", err)
		}
	}

	target_ids := make([]int64, len(str_ids))

	for idx, this_id := range str_ids {

		id, err := strconv.ParseInt(this_id, 10, 64)

		if err != nil {
			log.Fatalf("Failed to parse '%s', %v", this_id, err)
		}

		target_ids[idx] = id
	}

	// END OF put me in a function

	ctx := context.Background()

	r, err := reader.NewReader(ctx, *reader_uri)

	if err != nil {
		log.Fatalf("Failed to create reader for '%s', %v", *reader_uri, err)
	}

	wr, err := writer.NewWriter(ctx, *writer_uri)

	if err != nil {
		log.Fatalf("Failed to create new writer for '%s', %v", *writer_uri, err)
	}

	ex, err := export.NewExporter(ctx, *exporter_uri)

	if err != nil {
		log.Fatalf("Failed to create new exporter for '%s', %v", *exporter_uri, err)
	}

	source_body, err := wof_reader.LoadBytes(ctx, r, *source_id)

	if err != nil {
		log.Fatalf("Failed to load source '%d'", *source_id)
	}

	geom_rsp := gjson.GetBytes(source_body, "geometry")

	if !geom_rsp.Exists() {
		log.Fatal("Source record is missing a geometry")
	}

	source_geom := geom_rsp.Value()

	for _, id := range target_ids {

		target_body, err := wof_reader.LoadBytes(ctx, r, id)

		if err != nil {
			log.Fatalf("Failed to load target '%d'", id)
		}

		new_body, err := sjson.SetBytes(target_body, "geometry", source_geom)

		if err != nil {
			log.Fatalf("Failed to update target geometry for '%d', %v", id, err)
		}

		new_body, err = ex.Export(ctx, new_body)

		if err != nil {
			log.Fatalf("Failed to export target '%d', %v", id, err)
		}

		err = wof_writer.WriteBytes(ctx, wr, new_body)

		if err != nil {
			log.Fatalf("Failed to write target '%d', %v", id, err)
		}
	}
}
