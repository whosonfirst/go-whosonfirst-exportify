package exportify

import (
	"context"

	export "github.com/whosonfirst/go-whosonfirst-export/v2"
	wof_writer "github.com/whosonfirst/go-whosonfirst-writer"
	"github.com/whosonfirst/go-writer"
)

func ExportWithWriter(ctx context.Context, ex export.Exporter, wr writer.Writer, body []byte) error {

	var err error

	body, err = ex.Export(ctx, body)

	if err != nil {
		return err
	}

	err = wof_writer.WriteFeatureBytes(ctx, wr, body)

	if err != nil {
		return err
	}

	return nil
}
