# go-whosonfirst-exportify

Tools (written in Go) for exporting Who's On First records.

## Important

This is work in progress. Documentation to follow.

## Tools

```
$> make cli
go build -mod vendor -o bin/exportify cmd/exportify/main.go
```

### exportify

```
$> ./bin/exportify -h
Exportify one or more Who's On First IDs.

Usage:
	 ./bin/exportify [options] wof-id-(N) wof-id-(N)

For example:
	./bin/exportify -s . -i 1234
	./bin/exportify -reader-uri fs:///usr/local/data/whosonfirst-data-admin-ca/data -id 1234 -id 5678

Valid options are:
  -exporter-uri string
    	A valid whosonfirst/go-whosonfirst-export URI. (default "whosonfirst://")
  -i string
    	A valid Who's On First ID.
  -id value
    	One or more Who's On First IDs. If left empty the value of the -i flag will be used.
  -reader-uri string
    	A valid whosonfirst/go-reader URI. If empty the value of the -s flag will be used in combination with the fs:// scheme.
  -s string
    	A valid path to the root directory of the Who's On First data repository. If empty (and -reader-uri or -writer-uri are empty) the current working directory will be used and appended with a 'data' subdirectory.
  -writer-uri string
    	A valid whosonfirst/go-writer URI. If empty the value of the -s flag will be used in combination with the fs:// scheme.
```

## See also

* https://github.com/whosonfirst/go-whosonfirst-export
* https://github.com/whosonfirst/go-reader
* https://github.com/whosonfirst/go-writer