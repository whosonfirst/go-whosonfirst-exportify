package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"log"

	"github.com/whosonfirst/go-whosonfirst-export/v2"
	"github.com/whosonfirst/go-whosonfirst-iterate/v2/iterator"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"github.com/whosonfirst/go-writer/v3"
)

func main() {

	writer_uri := flag.String("writer-uri", "stdout://", "")
	iterator_uri := flag.String("iterator-uri", "repo://", "")

	flag.Parse()

	ctx := context.Background()
	uris := flag.Args()

	wr, err := writer.NewWriter(ctx, *writer_uri)

	if err != nil {
		log.Fatalf("Failed to create writer, %v", err)
	}

	export_opts, err := export.NewDefaultOptions(ctx)

	if err != nil {
		log.Fatalf("Failed to create export options, %v", err)
	}

	iter_cb := func(ctx context.Context, path string, r io.ReadSeeker, args ...interface{}) error {

		id, uri_args, err := uri.ParseURI(path)

		if err != nil {
			return fmt.Errorf("Failed to parse URI for %s, %w", path, err)
		}

		body, err := io.ReadAll(r)

		if err != nil {
			return fmt.Errorf("Failed to read %s, %w", path, err)
		}

		var buf bytes.Buffer
		buf_wr := bufio.NewWriter(&buf)

		has_changed, err := export.ExportChanged(body, body, export_opts, buf_wr)

		if err != nil {
			return fmt.Errorf("Failed to export %s, %w", path, err)
		}

		if !has_changed {
			return nil
		}
		
		buf_wr.Flush()
		br := bytes.NewReader(buf.Bytes())

		rel_path, err := uri.Id2RelPath(id, uri_args)

		if err != nil {
			return fmt.Errorf("Failed to derive rel_path for %d (%s), %w", id, path, err)
		}

		_, err = wr.Write(ctx, rel_path, br)

		if err != nil {
			return fmt.Errorf("Failed to write %s (for %s), %w", rel_path, path, err)
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

	wr.Close(ctx)
}
