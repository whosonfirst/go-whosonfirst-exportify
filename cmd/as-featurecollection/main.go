package main

import (
	_ "github.com/whosonfirst/go-writer-featurecollection"
)

import (
	"context"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-iterate/emitter"
	"github.com/whosonfirst/go-whosonfirst-iterate/iterator"
	"github.com/whosonfirst/go-writer"
	"io"
	"log"
	"net/url"
	"os"
	"strings"
)

const SCHEME_FEATURECOLLECTION string = "featurecollection://"

func main() {

	emitter_schemes := strings.Join(emitter.Schemes(), ",")
	emitter_desc := fmt.Sprintf("A valid whosonfirst/go-whosonfirst-iterator/emitter URI. Supported emitter URI schemes are: %s", emitter_schemes)

	schemes := make([]string, 0)

	for _, s := range writer.Schemes() {

		if s == SCHEME_FEATURECOLLECTION {
			continue
		}

		schemes = append(schemes, s)
	}

	writer_schemes := strings.Join(schemes, ", ")
	writer_desc := fmt.Sprintf("A valid whosonfirst/go-writer URI. Supported writer URI schemes are: %s", writer_schemes)

	iterator_uri := flag.String("iterator-uri", "repo://", emitter_desc)
	writer_uri := flag.String("writer-uri", "stdout://", writer_desc)

	flag.Usage = func() {

		fmt.Fprintf(os.Stderr, "Export one or more WOF records as a GeoJSON FeatureCollection\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options] path-(N) path-(N)\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "For example:\n")
		fmt.Fprintf(os.Stderr, "\t%s -iterator-uri 'repo://?include=properties.mz:is_current=1' /usr/local/data/sfomuseum-data-publicart/\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	uris := flag.Args()

	if strings.HasPrefix(*writer_uri, SCHEME_FEATURECOLLECTION) {
		log.Fatalf("Invalid -writer-uri")
	}

	ctx := context.Background()

	q := url.Values{}
	q.Set("writer", *writer_uri)

	u := url.URL{}
	u.Scheme = "featurecollection"
	u.RawQuery = q.Encode()

	*writer_uri = u.String()

	wr, err := writer.NewWriter(ctx, *writer_uri)

	if err != nil {
		log.Fatalf("Failed to create writer, %v", err)
	}

	defer wr.Close(ctx)

	iter_cb := func(ctx context.Context, fh io.ReadSeeker, args ...interface{}) error {
		_, err := wr.Write(ctx, "", fh)
		return err
	}

	iter, err := iterator.NewIterator(ctx, *iterator_uri, iter_cb)

	if err != nil {
		log.Fatalf("Failed to create iterator, %v", err)
	}

	err = iter.IterateURIs(ctx, uris...)

	if err != nil {
		log.Fatal("Failed to iterate URIs, %v", err)
	}

}