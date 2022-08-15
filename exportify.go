package exportify

import (
	"context"
	"fmt"

	export "github.com/whosonfirst/go-whosonfirst-export/v2"
	wof_writer "github.com/whosonfirst/go-whosonfirst-writer/v2"
	"github.com/whosonfirst/go-writer/v2"
)

func ExportWithWriter(ctx context.Context, ex export.Exporter, wr writer.Writer, body []byte) error {

	var err error

	body, err = ex.Export(ctx, body)

	if err != nil {
		return fmt.Errorf("Failed to export body, %w", err)
	}

	_, err = wof_writer.WriteBytes(ctx, wr, body)

	if err != nil {
		return fmt.Errorf("Failed to write bytes, %w", err)
	}

	return nil
}
