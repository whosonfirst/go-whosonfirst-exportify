package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/sfomuseum/go-csvdict"
	"github.com/sfomuseum/go-flags/multi"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-feature/properties"
	"github.com/whosonfirst/go-whosonfirst-iterate/v2/iterator"
)

func main() {

	var fields multi.MultiCSVString

	iterator_uri := flag.String("iterator-uri", "repo://", "")
	flag.Var(&fields, "field", "One or more relative 'properties.FIELDNAME' paths to include the CSV output. If the fieldname is 'path' the filename of the current record will be included. If the fieldname is 'centroid' the primary centroid of the current record will be derived and included as 'latitude' and 'longitude' columns.")

	flag.Usage = func() {

		fmt.Fprintf(os.Stderr, "Export one or more WOF records as a CSV document written to STDOUT\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options] path-(N) path-(N)\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "For example:\n")
		fmt.Fprintf(os.Stderr, "\t%s -field wof:id -field wof:name -field centroid -iterator-uri 'repo://?include=properties.mz:is_current=1' /usr/local/data/sfomuseum-data-publicart/\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	uris := flag.Args()

	ctx := context.Background()

	writers := []io.Writer{
		os.Stdout,
	}

	mw := io.MultiWriter(writers...)

	var csv_wr *csvdict.Writer

	mu := new(sync.RWMutex)

	iter_cb := func(ctx context.Context, path string, fh io.ReadSeeker, args ...interface{}) error {

		body, err := io.ReadAll(fh)

		if err != nil {
			return fmt.Errorf("Failed to read %s, %w", path, err)
		}

		out := make(map[string]string)

		for _, k := range fields {

			switch k {
			case "path":
				out[k] = path
			case "centroid":

				c, _, err := properties.Centroid(body)

				if err != nil {
					return fmt.Errorf("Failed to derive centroid for %s, %w", path, err)
				}

				out["latitude"] = strconv.FormatFloat(c.Lat(), 'f', -1, 64)
				out["longitude"] = strconv.FormatFloat(c.Lon(), 'f', -1, 64)

			default:

				fq_k := fmt.Sprintf("properties.%s", k)
				rsp := gjson.GetBytes(body, fq_k)
				out[k] = rsp.String()
			}
		}

		if csv_wr == nil {

			fieldnames := make([]string, 0)

			for f, _ := range out {
				fieldnames = append(fieldnames, f)
			}

			wr, err := csvdict.NewWriter(mw, fieldnames)

			if err != nil {
				return fmt.Errorf("Failed to create CSV writer, %w", err)
			}

			csv_wr = wr

			err = csv_wr.WriteHeader()

			if err != nil {
				return fmt.Errorf("Failed to write CSV header, %w", err)
			}
		}

		mu.Lock()
		defer mu.Unlock()

		err = csv_wr.WriteRow(out)

		if err != nil {
			return fmt.Errorf("Failed to write row for %s, %w", path, err)
		}

		return nil
	}

	iter, err := iterator.NewIterator(ctx, *iterator_uri, iter_cb)

	if err != nil {
		log.Fatalf("Failed to create iterator, %v", err)
	}

	err = iter.IterateURIs(ctx, uris...)

	if err != nil {
		log.Fatalf("Failed to iterate URIs, %v", err)
	}

	csv_wr.Flush()
}
