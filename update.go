package exportify

import (
	"context"
	"github.com/sfomuseum/go-flags/multi"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type UpdateFeatureOptions struct {
	StringProperties  multi.KeyValueString
	Int64Properties   multi.KeyValueInt64
	Float64Properties multi.KeyValueFloat64
	// JSONProperties multi.KeyValueString
}

func UpdateFeature(ctx context.Context, body []byte, opts *UpdateFeatureOptions) ([]byte, bool, error) {

	changed := false

	for _, p := range opts.StringProperties {

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
				return nil, false, err
			}

			body = new_body
			changed = true
		}
	}

	for _, p := range opts.Int64Properties {

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
				return nil, false, err
			}

			body = new_body
			changed = true
		}
	}

	for _, p := range opts.Float64Properties {

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
				return nil, false, err
			}

			body = new_body
			changed = true
		}
	}

	return body, changed, nil
}
