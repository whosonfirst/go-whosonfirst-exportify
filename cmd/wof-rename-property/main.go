package main

import (
	"context"
	"flag"
	"io"
	"io/ioutil"
	"log"

	_ "github.com/sfomuseum/go-flags/multi"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	export "github.com/whosonfirst/go-whosonfirst-export/v2"
	"github.com/whosonfirst/go-whosonfirst-iterate/v2/iterator"
	wof_writer "github.com/whosonfirst/go-whosonfirst-writer/v3"
	"github.com/whosonfirst/go-writer/v3"
)

func main() {

	iterator_uri := flag.String("indexer-uri", "repo://", "A valid whosonfirst/go-whosonfirst-iterate/v2 URI.")
	exporter_uri := flag.String("exporter-uri", "whosonfirst://", "A valid whosonfirst/go-whosonfirst-export URI.")
	writer_uri := flag.String("writer-uri", "null://", "A valid whosonfirst/go-writer URI.")

	old_property := flag.String("old-property", "", "The fully qualified path of the property to rename.")
	new_property := flag.String("new-property", "", "The fully qualified path of the property to be (re)named.")

	flag.Parse()

	ctx := context.Background()

	ex, err := export.NewExporter(ctx, *exporter_uri)

	if err != nil {
		log.Fatalf("Failed to create exporter for '%s', %v", *exporter_uri, err)
	}

	wr, err := writer.NewWriter(ctx, *writer_uri)

	if err != nil {
		log.Fatalf("Failed to create writer for '%s', %v", *writer_uri, err)
	}

	cb := func(ctx context.Context, path string, fh io.ReadSeeker, args ...interface{}) error {

		body, err := ioutil.ReadAll(fh)

		if err != nil {
			return err
		}

		old_rsp := gjson.GetBytes(body, *old_property)

		if !old_rsp.Exists() {
			return nil
		}

		body, err = sjson.SetBytes(body, *new_property, old_rsp.Value())

		if err != nil {
			return err
		}

		body, err = sjson.DeleteBytes(body, *old_property)

		if err != nil {
			return err
		}

		new_body, err := ex.Export(ctx, body)

		if err != nil {
			return err
		}

		_, err = wof_writer.WriteBytes(ctx, wr, new_body)

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
