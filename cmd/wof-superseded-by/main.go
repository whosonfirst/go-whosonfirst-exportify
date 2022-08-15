// superseded-by will update the `wof:superseded_by` and `wof:supersedes` properties for one or more sets of Who's On First records.
// Additionally it will assign the `mz:is_current=0` property for the records being superseded.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/sfomuseum/go-flags/multi"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-reader"
	export "github.com/whosonfirst/go-whosonfirst-export/v2"
	exportify "github.com/whosonfirst/go-whosonfirst-exportify"
	wof_reader "github.com/whosonfirst/go-whosonfirst-reader"
	"github.com/whosonfirst/go-writer/v2"
)

func main() {

	source := flag.String("s", "", "A valid path to the root directory of the Who's On First data repository. If empty (and -reader-uri or -writer-uri are empty) the current working directory will be used and appended with a 'data' subdirectory.")
	id := flag.String("i", "", "A valid Who's On First ID.")

	reader_uri := flag.String("reader-uri", "", "A valid whosonfirst/go-reader URI. If empty the value of the -s flag will be used in combination with the fs:// scheme.")
	writer_uri := flag.String("writer-uri", "", "A valid whosonfirst/go-writer URI. If empty the value of the -s flag will be used in combination with the fs:// scheme.")

	exporter_uri := flag.String("exporter-uri", "whosonfirst://", "A valid whosonfirst/go-whosonfirst-export URI.")

	var ids multi.MultiInt64
	flag.Var(&ids, "id", "One or more Who's On First IDs. If left empty the value of the -i flag will be used.")

	var superseded_by multi.MultiInt64
	flag.Var(&superseded_by, "by", "Zero or more Who's On First IDs that the records being deprecated are superseded by.")

	flag.Usage = func() {

		fmt.Fprintf(os.Stderr, "\"Supersede\" one or more Who's On First IDs.\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "For example:\n")
		fmt.Fprintf(os.Stderr, "\t%s -reader-uri fs:///usr/local/data/sfomuseum-data-enterprise/data -id 1159286017 -by 1159283849\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *reader_uri == "" || *writer_uri == "" {

		var path string

		if *source == "" {

			cwd, err := os.Getwd()

			if err != nil {
				log.Fatalf("Failed to determine current working directory, %v", err)
			}

			path = cwd

		} else {

			abs_source, err := filepath.Abs(*source)

			if err != nil {
				log.Fatalf("Failed to derive absolute path for '%s', %v", *source, err)
			}

			path = abs_source
		}

		abs_path := fmt.Sprintf("fs://%s/data", path)

		if *reader_uri == "" {
			*reader_uri = abs_path
		}

		if *writer_uri == "" {
			*writer_uri = abs_path
		}
	}

	if *id != "" {

		err := ids.Set(*id)

		if err != nil {
			log.Fatalf("Failed to assign '%s' (-id) flag, %v", *id, err)
		}
	}

	ctx := context.Background()

	ex, err := export.NewExporter(ctx, *exporter_uri)

	if err != nil {
		log.Fatalf("Failed create exporter for '%s', %v", *exporter_uri, err)
	}

	r, err := reader.NewReader(ctx, *reader_uri)

	if err != nil {
		log.Fatalf("Failed to create reader for '%s', %v", *reader_uri, err)
	}

	wr, err := writer.NewWriter(ctx, *writer_uri)

	if err != nil {
		log.Fatalf("Failed to create writer for '%s', %v", *writer_uri, err)
	}

	for _, id := range ids {

		err := supersededById(ctx, r, wr, ex, id, superseded_by)

		if err != nil {
			log.Fatalf("Failed to supersede record for '%d', %v", id, err)
		}

		err = supersedesId(ctx, r, wr, ex, superseded_by, id)

		if err != nil {
			log.Fatalf("Failed to update wof:supersedes properties for superseding records, %v", err)
		}
	}
}

// supersededById ensures that all the IDs in 'superseded_by' are present in the `wof:superseded_by` property of the record defined by 'id'.
func supersededById(ctx context.Context, r reader.Reader, wr writer.Writer, ex export.Exporter, id int64, superseded_by multi.MultiInt64) error {

	body, err := wof_reader.LoadBytes(ctx, r, id)

	if err != nil {
		return err
	}

	to_update := map[string]interface{}{
		"properties.mz:is_current": 0,
	}

	tmp := make(map[int64]bool)

	for _, id := range superseded_by {
		tmp[id] = true
	}

	rsp := gjson.GetBytes(body, "properties.wof:superseded_by")

	for _, r := range rsp.Array() {
		id := r.Int()
		tmp[id] = true
	}

	new_list := make([]int64, 0)

	for id := range tmp {

		switch id {
		case 0, -1:
			continue
		default:
			new_list = append(new_list, id)
		}
	}

	to_update["properties.wof:superseded_by"] = new_list

	new_body, err := export.AssignProperties(ctx, body, to_update)

	if err != nil {
		return fmt.Errorf("failed to assign properties for %d, %w", id, err)
	}

	err = exportify.ExportWithWriter(ctx, ex, wr, new_body)

	if err != nil {
		return fmt.Errorf("failed to write data for %d, %v", id, err)
	}

	return nil
}

// To do: Reconcile this with equivalent code in cmd/deprecate

// supersedesId ensures that 'id' is present in the `wof:supersedes` property of all the records defined by 'superseded_by'.
func supersedesId(ctx context.Context, r reader.Reader, wr writer.Writer, ex export.Exporter, superseded_by multi.MultiInt64, id int64) error {

	for _, sid := range superseded_by {

		body, err := wof_reader.LoadBytes(ctx, r, sid)

		if err != nil {
			return fmt.Errorf("failed to load record for %d, %w", sid, err)
		}

		to_update := map[string]interface{}{}

		tmp := make(map[int64]bool)
		tmp[id] = true

		rsp := gjson.GetBytes(body, "properties.wof:supersedes")

		for _, r := range rsp.Array() {
			id := r.Int()
			tmp[id] = true
		}

		new_list := make([]int64, 0)

		for id := range tmp {

			switch id {
			case 0, -1:
				continue
			default:
				new_list = append(new_list, id)
			}
		}

		to_update["properties.wof:supersedes"] = new_list

		new_body, err := export.AssignProperties(ctx, body, to_update)

		if err != nil {
			return fmt.Errorf("failed to assign properties for %d, %w", sid, err)
		}

		err = exportify.ExportWithWriter(ctx, ex, wr, new_body)

		if err != nil {
			return fmt.Errorf("failed to write data for %d, %v", sid, err)
		}

	}

	return nil
}
