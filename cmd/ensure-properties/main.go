package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/aaronland/go-json-query"
	"github.com/sfomuseum/go-flags/multi"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/whosonfirst/go-whosonfirst-export/v2"
	"github.com/whosonfirst/go-whosonfirst-index"
	_ "github.com/whosonfirst/go-whosonfirst-index/fs"
	wof_writer "github.com/whosonfirst/go-whosonfirst-writer"
	"github.com/whosonfirst/go-writer"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

func main() {

	indexer_uri := flag.String("indexer-uri", "repo://", "A valid whosonfirst/go-whosonfirst-index URI.")
	exporter_uri := flag.String("exporter-uri", "whosonfirst://", "A valid whosonfirst/go-whosonfirst-export URI.")
	writer_uri := flag.String("writer-uri", "null://", "A valid whosonfirst/go-writer URI.")

	var str_properties multi.KeyValueString
	flag.Var(&str_properties, "string-property", "One or more {KEY}={VALUE} flags where {KEY} is a valid tidwall/gjson path and {VALUE} is a string value.")

	var int_properties multi.KeyValueInt64
	flag.Var(&int_properties, "int-property", "One or more {KEY}={VALUE} flags where {KEY} is a valid tidwall/gjson path and {VALUE} is a int(64) value.")

	var float_properties multi.KeyValueFloat64
	flag.Var(&float_properties, "float-property", "One or more {KEY}={VALUE} flags where {KEY} is a valid tidwall/gjson path and {VALUE} is a float(64) value.")

	var queries query.QueryFlags
	flag.Var(&queries, "query", "One or more {PATH}={REGEXP} parameters for filtering records.")

	valid_query_modes := strings.Join([]string{query.QUERYSET_MODE_ALL, query.QUERYSET_MODE_ANY}, ", ")
	desc_query_modes := fmt.Sprintf("Specify how query filtering should be evaluated. Valid modes are: %s", valid_query_modes)

	query_mode := flag.String("query-mode", query.QUERYSET_MODE_ALL, desc_query_modes)

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

	var qs *query.QuerySet

	if len(queries) > 0 {

		qs = &query.QuerySet{
			Queries: queries,
			Mode:    *query_mode,
		}
	}

	cb := func(ctx context.Context, fh io.Reader, args ...interface{}) error {

		path, err := index.PathForContext(ctx)

		if err != nil {
			return err
		}

		body, err := ioutil.ReadAll(fh)

		if err != nil {
			return err
		}

		if qs != nil {

			matches, err := query.Matches(ctx, qs, body)

			if err != nil {
				return err
			}

			if !matches {
				return nil
			}

		}

		changed := false

		for _, p := range str_properties {

			path := p.Key()
			new_value := p.Value()

			update := true

			old_rsp := gjson.GetBytes(body, path)

			if old_rsp.Exists() {

				old_value := old_rsp.String()

				if old_value == new_value {
					update = false
				}
			}

			if update {

				new_body, err := sjson.SetBytes(body, path, new_value)

				if err != nil {
					return nil
				}

				body = new_body
				changed = true
			}
		}

		for _, p := range int_properties {

			path := p.Key()
			new_value := p.Value()

			update := true

			old_rsp := gjson.GetBytes(body, path)

			if old_rsp.Exists() {

				old_value := old_rsp.Int()

				if old_value == new_value {
					update = false
				}
			}

			if update {

				new_body, err := sjson.SetBytes(body, path, new_value)

				if err != nil {
					return nil
				}

				body = new_body
				changed = true
			}
		}

		for _, p := range float_properties {

			path := p.Key()
			new_value := p.Value()

			update := true

			old_rsp := gjson.GetBytes(body, path)

			if old_rsp.Exists() {

				old_value := old_rsp.Float()

				if old_value == new_value {
					update = false
				}
			}

			if update {

				new_body, err := sjson.SetBytes(body, path, new_value)

				if err != nil {
					return nil
				}

				body = new_body
				changed = true
			}
		}

		if !changed {
			return nil
		}

		new_body, err := ex.Export(ctx, body)

		if err != nil {
			return err
		}

		err = wof_writer.WriteFeatureBytes(ctx, wr, new_body)

		if err != nil {
			return err
		}

		log.Printf("Updated %s\n", path)
		return nil
	}

	i, err := index.NewIndexer(*indexer_uri, cb)

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	paths := flag.Args()

	err = i.Index(ctx, paths...)

	if err != nil {
		log.Fatal(err)
	}
}
