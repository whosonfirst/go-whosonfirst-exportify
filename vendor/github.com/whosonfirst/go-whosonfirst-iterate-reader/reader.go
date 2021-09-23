package reader

import (
	"context"
	"fmt"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-iterate/emitter"
	"github.com/whosonfirst/go-whosonfirst-iterate/filters"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"net/url"
)

func init() {
	ctx := context.Background()
	emitter.RegisterEmitter(ctx, "reader", NewReaderEmitter)
}

type ReaderEmitter struct {
	emitter.Emitter
	reader  reader.Reader
	filters filters.Filters
}

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

func (idx *ReaderEmitter) WalkURI(ctx context.Context, index_cb emitter.EmitterCallbackFunc, path string) error {

	id, uri_args, err := uri.ParseURI(path)

	if err != nil {
		return fmt.Errorf("Failed to parse '%s', %w", path, err)
	}

	rel_path, err := uri.Id2RelPath(id, uri_args)

	if err != nil {
		return fmt.Errorf("Failed to derived relative path for '%s', %v", path, err)
	}

	fh, err := idx.reader.Read(ctx, rel_path)

	if err != nil {
		return fmt.Errorf("Failed to read path (%s) for '%s', %v", rel_path, path, err)
	}

	defer fh.Close()

	if idx.filters != nil {

		ok, err := idx.filters.Apply(ctx, fh)

		if err != nil {
			return err
		}

		if !ok {
			return nil
		}

		_, err = fh.Seek(0, 0)

		if err != nil {
			return err
		}
	}

	ctx = emitter.AssignPathContext(ctx, rel_path)
	return index_cb(ctx, fh)
}
