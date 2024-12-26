package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/planar"
	"github.com/whosonfirst/go-whosonfirst-iterate/v2/emitter"
	"github.com/whosonfirst/go-whosonfirst-iterate/v2/iterator"
	_ "github.com/whosonfirst/go-whosonfirst-iterate-git/v2"	
	_ "github.com/whosonfirst/go-writer-featurecollection/v3"
	"github.com/whosonfirst/go-writer/v3"
)

const SCHEME_FEATURECOLLECTION string = "featurecollection://"

func main() {

	emitter_schemes := strings.Join(emitter.Schemes(), ",")
	emitter_desc := fmt.Sprintf("A valid whosonfirst/go-whosonfirst-iterator/v2 URI. Supported emitter URI schemes are: %s", emitter_schemes)

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

	as_multipoints := flag.Bool("as-multipoints", false, "Output geometries as a MultiPoint array")

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

	iter_cb := func(ctx context.Context, path string, fh io.ReadSeeker, args ...interface{}) error {

		if *as_multipoints {

			body, err := io.ReadAll(fh)

			if err != nil {
				return err
			}

			f, err := geojson.UnmarshalFeature(body)

			if err != nil {
				return err
			}

			geom := f.Geometry

			f.Geometry = geom

			switch geom.GeoJSONType() {
			case "Point":

				points := []orb.Point{
					geom.(orb.Point),
				}

				geom = orb.MultiPoint(points)

			case "MultiPoint":
				// pass
			case "MultiPolygon":

				points := make([]orb.Point, 0)

				for _, poly := range geom.(orb.MultiPolygon) {
					pt, _ := planar.CentroidArea(poly)
					points = append(points, pt)
				}

				geom = orb.MultiPoint(points)

			case "Polygon":

				pt, _ := planar.CentroidArea(geom)
				points := []orb.Point{pt}

				geom = orb.MultiPoint(points)

			default:
				return fmt.Errorf("unsupported geometry type %s", geom.GeoJSONType())
			}

			f.Geometry = geom

			body, err = f.MarshalJSON()

			if err != nil {
				return err
			}

			fh = bytes.NewReader(body)
		}

		_, err := wr.Write(ctx, "", fh)
		return err
	}

	iter, err := iterator.NewIterator(ctx, *iterator_uri, iter_cb)

	if err != nil {
		log.Fatalf("Failed to create iterator, %v", err)
	}

	err = iter.IterateURIs(ctx, uris...)

	if err != nil {
		log.Fatalf("Failed to iterate URIs, %v", err)
	}

}
