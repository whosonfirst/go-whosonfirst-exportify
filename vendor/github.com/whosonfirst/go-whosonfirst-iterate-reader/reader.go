package reader

import (
	"context"
	"fmt"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-iterate/v2/emitter"
	"github.com/whosonfirst/go-whosonfirst-iterate/v2/filters"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"net/url"
)

func init() {
	ctx := context.Background()
	emitter.RegisterEmitter(ctx, "reader", NewReaderEmitter)
}

// GitEmitter implements the `Emitter` interface for crawling records with a `whosonfirst/go-reader.Reader` instance.
type ReaderEmitter struct {
	emitter.Emitter
	// reader is the `whosonfirst/go-reader.Reader` instance used to reade documents
	reader reader.Reader
	// filters is a `filters.Filters` instance used to include or exclude specific records from being crawled.
	filters filters.Filters
}

// NewGitEmitter() returns a new `GitEmitter` instance configured by 'uri' in the form of:
//
//	reader://?{PARAMETERS}
//
// Where {PATH} is an optional path on disk where a repository will be clone to (default is to clone repository in memory) and {PARAMETERS} may be:
// * `?include=` Zero or more `aaronland/go-json-query` query strings containing rules that must match for a document to be considered for further processing.
// * `?exclude=` Zero or more `aaronland/go-json-query`	query strings containing rules that if matched will prevent a document from being considered for further processing.
// * `?include_mode=` A valid `aaronland/go-json-query` query mode string for testing inclusion rules.
// * `?exclude_mode=` A valid `aaronland/go-json-query` query mode string for testing exclusion rules.
// * `?reader=` A valid `whosonfirst/go-reader` URI used to create the underlying reader instance.
func NewReaderEmitter(ctx context.Context, uri string) (emitter.Emitter, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	q := u.Query()

	reader_uri := q.Get("reader")

	r, err := reader.NewReader(ctx, reader_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new reader, %w", err)
	}

	f, err := filters.NewQueryFiltersFromURI(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive query filters from URI, %w", err)
	}

	idx := &ReaderEmitter{
		reader:  r,
		filters: f,
	}

	return idx, nil
}

// WalkURI() reads 'path' using the underlying reader (if not excluded by any filters specified when `idx` was
// created) and invokes 'index_cb'.
func (idx *ReaderEmitter) WalkURI(ctx context.Context, index_cb emitter.EmitterCallbackFunc, path string) error {

	id, uri_args, err := uri.ParseURI(path)

	if err != nil {
		return fmt.Errorf("Failed to parse '%s', %w", path, err)
	}

	rel_path, err := uri.Id2RelPath(id, uri_args)

	if err != nil {
		return fmt.Errorf("Failed to derived relative path for '%s', %w", path, err)
	}

	fh, err := idx.reader.Read(ctx, rel_path)

	if err != nil {
		return fmt.Errorf("Failed to read path (%s) for '%s', %w", rel_path, path, err)
	}

	defer fh.Close()

	if idx.filters != nil {

		ok, err := idx.filters.Apply(ctx, fh)

		if err != nil {
			return fmt.Errorf("Failed to apply filters to %s, %w", rel_path, err)
		}

		if !ok {
			return nil
		}

		_, err = fh.Seek(0, 0)

		if err != nil {
			return fmt.Errorf("Failed to reset filehandle for %s, %w", rel_path, err)
		}
	}

	return index_cb(ctx, rel_path, fh)
}
