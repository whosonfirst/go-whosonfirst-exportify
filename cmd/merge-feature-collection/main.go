package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/sfomuseum/go-flags/multi"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-export/v2"
	wof_reader "github.com/whosonfirst/go-whosonfirst-reader"
	wof_writer "github.com/whosonfirst/go-whosonfirst-writer"
	"github.com/whosonfirst/go-writer"
	"io/ioutil"
	"log"
	"os"
)

func open(ctx context.Context, path string) ([]byte, error) {

	fh, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer fh.Close()

	return ioutil.ReadAll(fh)
}

func main() {

	reader_uri := flag.String("reader-uri", "", "A valid whosonfirst/go-reader URI")
	writer_uri := flag.String("writer-uri", "", "A valid whosonfirst/go-writer URI")

	exporter_uri := flag.String("exporter-uri", "whosonfirst://", "A valid whosonfirst/go-whosonfirst-export URI")

	var to_append multi.MultiString
	flag.Var(&to_append, "path", "One or more valid tidwall/gjson paths. These will be copied from the source GeoJSON feature to the corresponding WOF record.")

	flag.Usage = func() {

		fmt.Fprintf(os.Stderr, "Upate one or more Who's On First records with matching entries in a GeoJSON FeatureCollection file.\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options] path(N) path(N)\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "For example:\n")
		fmt.Fprintf(os.Stderr, "\t%s -reader-uri fs:///usr/local/data/whosonfirst-data-admin-ca/data -writer-uri fs:///usr/local/data/whosonfirst-data-admin-ca/data -path geometry -path 'properties.example:property' /usr/local/data/updates.geojson\n\n", os.Args[0])
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

		fc_b, err := open(ctx, path)

		if err != nil {
			log.Fatalf("Failed to open '%s', %v", path, err)
		}

		f_rsp := gjson.GetBytes(fc_b, "features")

		for idx, qgis_f := range f_rsp.Array() {

			id_rsp := qgis_f.Get("properties.wof:id")

			if !id_rsp.Exists() {
				log.Fatalf("Missing wof:id property for updated feature '%d'", idx)
			}

			wof_id := id_rsp.Int()

			wof_f, err := wof_reader.LoadBytesFromID(ctx, r, wof_id)

			if err != nil {
				log.Fatalf("Failed to load '%d', %v", wof_id, err)
			}

			for _, path := range to_append {

				v := qgis_f.Get(path)

				if !v.Exists() {
					log.Fatalf("Missing '%s' path in updated feature for '%d'", path, wof_id)
				}

				wof_f, err = sjson.SetBytes(wof_f, path, v.Value())

				if err != nil {
					log.Fatalf("Failed to set '%s' for '%d', %v", path, wof_id, err)
				}
			}

			wof_f, err = ex.Export(ctx, wof_f)

			if err != nil {
				log.Fatalf("Failed to export new record for '%d', %v", wof_id, err)
			}

			err = wof_writer.WriteFeatureBytes(ctx, wr, wof_f)

			if err != nil {
				log.Fatalf("Failed to write record for '%d', %v", wof_id, err)
			}

		}

	}
}
