package exportify

import (
	"context"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/whosonfirst/go-whosonfirst-export/v2"
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

func SupersedeRecord(ctx context.Context, ex export.Exporter, old_body []byte) ([]byte, []byte, error) {

	id_rsp := gjson.GetBytes(old_body, "properties.wof:id")

	if !id_rsp.Exists() {
		return nil, nil, fmt.Errorf("Failed to derive old properties.wof:id property for record being superseded")
	}

	old_id := id_rsp.Int()

	// Create the new record

	new_body := old_body

	new_body, err := sjson.DeleteBytes(new_body, "properties.wof:id")

	if err != nil {
		return nil, nil, err
	}

	new_body, err = ex.Export(ctx, new_body)

	if err != nil {
		return nil, nil, err
	}

	id_rsp = gjson.GetBytes(new_body, "properties.wof:id")

	if !id_rsp.Exists() {
		return nil, nil, fmt.Errorf("Failed to derive new properties.wof:id property for record superseding '%d'", old_id)
	}

	new_id := id_rsp.Int()

	// Update the new record

	new_body, err = sjson.SetBytes(new_body, "properties.wof:supsersedes", []int64{old_id})

	if err != nil {
		return nil, nil, err
	}

	// Update the old record

	to_update := map[string]interface{}{
		"properties.mz:is_current":      0,
		"properties.wof:supserseded_by": []int64{new_id},
	}

	old_body, err = AssignProperties(ctx, old_body, to_update)

	if err != nil {
		return nil, nil, err
	}

	return old_body, new_body, nil
}

func AssignProperties(ctx context.Context, body []byte, to_assign map[string]interface{}) ([]byte, error) {

	var err error

	for path, v := range to_assign {

		body, err = sjson.SetBytes(body, path, v)

		if err != nil {
			return nil, err
		}
	}

	return body, nil
}
