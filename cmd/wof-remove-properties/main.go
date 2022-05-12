package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"

	"github.com/sfomuseum/go-flags/multi"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	_ "github.com/whosonfirst/go-whosonfirst-iterate-reader"
	"github.com/whosonfirst/go-whosonfirst-iterate/v2/iterator"
	uri "github.com/whosonfirst/go-whosonfirst-uri"
	wof_writer "github.com/whosonfirst/go-whosonfirst-writer"
	"github.com/whosonfirst/go-writer"
)

func main() {

	iterator_uri := flag.String("indexer-uri", "repo://", "A valid whosonfirst/go-whosonfirst-iterate/v2 URI.")
	writer_uri := flag.String("writer-uri", "null://", "A valid whosonfirst/go-writer URI.")

	var properties multi.MultiString
	flag.Var(&properties, "property", "One or more (fully-qualified) properties to remove")

	flag.Parse()

	ctx := context.Background()

	wr, err := writer.NewWriter(ctx, *writer_uri)

	if err != nil {
		log.Fatalf("Failed to create writer for '%s', %v", *writer_uri, err)
	}

	cb := func(ctx context.Context, path string, fh io.ReadSeeker, args ...interface{}) error {

		_, uri_args, err := uri.ParseURI(path)

		if err != nil {
			return err
		}

		if uri_args.IsAlternate {
			log.Printf("Alternate files (%s) are not supported yet, skipping\n", path)
			return nil
		}

		body, err := io.ReadAll(fh)

		if err != nil {
			return err
		}

		changed := false

		for _, path := range properties {

			rsp := gjson.GetBytes(body, path)

			if !rsp.Exists() {
				continue
			}

			var err error

			body, err = sjson.DeleteBytes(body, path)

			if err != nil {
				return fmt.Errorf("failed to delete %s, %w", path, err)
			}

			changed = true
		}

		if !changed {
			return nil
		}

		err = wof_writer.WriteBytes(ctx, wr, body)

		if err != nil {
			return err
		}

		log.Printf("Updated %s\n", path)
		return nil
	}

	iter, err := iterator.NewIterator(ctx, *iterator_uri, cb)

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	paths := flag.Args()

	err = iter.IterateURIs(ctx, paths...)

	if err != nil {
		log.Fatal(err)
	}
}
