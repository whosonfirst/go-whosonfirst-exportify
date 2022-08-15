package main

import (
	"context"
	"flag"
	"io"
	"io/ioutil"
	"log"

	"github.com/sfomuseum/go-flags/multi"
	export "github.com/whosonfirst/go-whosonfirst-export/v2"
	exportify "github.com/whosonfirst/go-whosonfirst-exportify"
	_ "github.com/whosonfirst/go-whosonfirst-iterate-reader"
	"github.com/whosonfirst/go-whosonfirst-iterate/v2/iterator"
	uri "github.com/whosonfirst/go-whosonfirst-uri"
	wof_writer "github.com/whosonfirst/go-whosonfirst-writer"
	"github.com/whosonfirst/go-writer/v2"
)

func main() {

	iterator_uri := flag.String("indexer-uri", "repo://", "A valid whosonfirst/go-whosonfirst-iterate/v2 URI.")
	exporter_uri := flag.String("exporter-uri", "whosonfirst://", "A valid whosonfirst/go-whosonfirst-export URI.")
	writer_uri := flag.String("writer-uri", "null://", "A valid whosonfirst/go-writer URI.")

	var str_properties multi.KeyValueString
	flag.Var(&str_properties, "string-property", "One or more {KEY}={VALUE} flags where {KEY} is a valid tidwall/gjson path and {VALUE} is a string value.")

	var int_properties multi.KeyValueInt64
	flag.Var(&int_properties, "int-property", "One or more {KEY}={VALUE} flags where {KEY} is a valid tidwall/gjson path and {VALUE} is a int(64) value.")

	var float_properties multi.KeyValueFloat64
	flag.Var(&float_properties, "float-property", "One or more {KEY}={VALUE} flags where {KEY} is a valid tidwall/gjson path and {VALUE} is a float(64) value.")

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

		_, uri_args, err := uri.ParseURI(path)

		if err != nil {
			return err
		}

		if uri_args.IsAlternate {
			log.Printf("Alternate files (%s) are not supported yet, skipping\n", path)
			return nil
		}

		body, err := ioutil.ReadAll(fh)

		if err != nil {
			return err
		}

		opts := &exportify.UpdateFeatureOptions{
			StringProperties:  str_properties,
			Int64Properties:   int_properties,
			Float64Properties: float_properties,
		}

		new_body, changed, err := exportify.UpdateFeature(ctx, body, opts)

		if err != nil {
			return err
		}

		if !changed {
			return nil
		}

		new_body, err = ex.Export(ctx, new_body)

		if err != nil {
			return err
		}

		err = wof_writer.WriteBytes(ctx, wr, new_body)

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
