# go-whosonfirst-iterate-reader

Go package implementing go-whosonfirst-iterate/emitter functionality using whosonfirst/go-reader.Reader instances.

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/whosonfirst/go-whosonfirst-iterate-reader.svg)](https://pkg.go.dev/github.com/whosonfirst/go-whosonfirst-iterate-reader)

## Example

```
import (
	"context"
	"github.com/whosonfirst/go-whosonfirst-iterate/v2/iterator"
	_ "github.com/whosonfirst/go-whosonfirst-iterate-reader"		
	"io"
	"log"
)

func main() {

	ctx := context.Background()
	
	iter_cb := func(ctx context.Context, path string, fh io.ReadSeeker, args ...interface{}) error {
		log.Println(path)	       
		return nil
	}

	iter_uri := "reader://?reader=fs%3A%2F%2F%2Fusr%2Flocal%2Fwhosonfirst%2Fgo-whosonfirst-iterate-reader%2Ffixtures%2Fdata"
	
	iter, _ := iterator.NewIterator(ctx, iter_uri, emitter_cb)

	iter.IterateURIs(ctx, "102/527/513/102527513.geojson")
``

## Tools

### emit

```
$> go run -mod vendor cmd/emit/main.go \
	-emitter-uri 'reader://?reader=fs:///usr/local/data/sfomuseum-data-whosonfirst/data' 1 \

| jq '.["properties"]["wof:name"]'

"Null Island"
```

## See also

* https://github.com/whosonfirst/go-whosonfirst-iterate
* https://github.com/whosonfirst/go-reader 