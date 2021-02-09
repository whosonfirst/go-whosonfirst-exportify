# go-whosonfirst-exportify

Tools (written in Go) for exporting Who's On First records.

## Important

This is work in progress. Documentation is incomplete.

## Tools

```
$> make cli
go build -mod vendor -o bin/exportify cmd/exportify/main.go
```

As of this writing these tools may contain duplicate, or at least common, code that would be well-served from being moved in to a package or library. That hasn't happened yet.

### deprecate-and-supersede

```
> ./bin/deprecate-and-supersede -h
Deprecate and supersede one or more Who's On First IDs.

Usage:
	 ./bin/deprecate-and-supersede [options] wof-id-(N) wof-id-(N)

For example:
	./bin/deprecate-and-supersede -s . -i 1234
	./bin/deprecate-and-supersede -reader-uri fs:///usr/local/data/whosonfirst-data-admin-ca/data -id 1234 -id 5678

Valid options are:
  -exporter-uri string
    	A valid whosonfirst/go-whosonfirst-export URI. (default "whosonfirst://")
  -float-property value
    	One or more {KEY}={VALUE} properties to append to the new record where {KEY} is a valid tidwall/gjson path and {VALUE} is a float(64) value.
  -i string
    	A valid Who's On First ID.
  -id value
    	One or more Who's On First IDs. If left empty the value of the -i flag will be used.
  -int-property value
    	One or more {KEY}={VALUE} properties to append to the new record where {KEY} is a valid tidwall/gjson path and {VALUE} is a int(64) value.
  -reader-uri string
    	A valid whosonfirst/go-reader URI. If empty the value of the -s flag will be used in combination with the fs:// scheme.
  -s string
    	A valid path to the root directory of the Who's On First data repository. If empty (and -reader-uri or -writer-uri are empty) the current working directory will be used and appended with a 'data' subdirectory.
  -string-property value
    	One or more {KEY}={VALUE} properties to append to the new record where {KEY} is a valid tidwall/gjson path and {VALUE} is a string value.
  -writer-uri string
    	A valid whosonfirst/go-writer URI. If empty the value of the -s flag will be used in combination with the fs:// scheme.
```

For example:

```
$> ./bin/deprecate-and-supersede \
	-s /usr/local/data/sfomuseum-data-architecture \
	-i 1477856003 \
	-string-property 'properties.wof:placetype=arcade'
	
2021/02/09 13:49:58 1477856003 replaced by 1729791935
```

### ensure-properties

```
> ./bin/ensure-properties -h
Usage of ./bin/ensure-properties:
  -exporter-uri string
    	A valid whosonfirst/go-whosonfirst-export URI. (default "whosonfirst://")
  -float-property value
    	One or more {KEY}={VALUE} flags where {KEY} is a valid tidwall/gjson path and {VALUE} is a float(64) value.
  -indexer-uri string
    	A valid whosonfirst/go-whosonfirst-index URI. (default "repo://")
  -int-property value
    	One or more {KEY}={VALUE} flags where {KEY} is a valid tidwall/gjson path and {VALUE} is a int(64) value.
  -string-property value
    	One or more {KEY}={VALUE} flags where {KEY} is a valid tidwall/gjson path and {VALUE} is a string value.
  -writer-uri string
    	A valid whosonfirst/go-writer URI. (default "null://")
```

For example:

```
$> ./bin/ensure-properties \
	-string-property 'properties.wof:repo=sfomuseum-data-architecture' \
	-writer-uri fs:///usr/local/data/sfomuseum-data-architecture/data \
	/usr/local/data/sfomuseum-data-architecture
	
2021/02/05 14:43:51 Updated /usr/local/data/sfomuseum-data-architecture/data/172/977/140/5/1729771405.geojson
2021/02/05 14:43:51 Updated /usr/local/data/sfomuseum-data-architecture/data/172/977/140/9/1729771409.geojson
2021/02/05 14:43:51 Updated /usr/local/data/sfomuseum-data-architecture/data/172/977/140/3/1729771403.geojson
2021/02/05 14:43:51 Updated /usr/local/data/sfomuseum-data-architecture/data/172/977/141/1/1729771411.geojson
2021/02/05 14:43:51 Updated /usr/local/data/sfomuseum-data-architecture/data/172/977/141/5/1729771415.geojson
2021/02/05 14:43:51 Updated /usr/local/data/sfomuseum-data-architecture/data/172/977/142/1/1729771421.geojson
2021/02/05 14:43:51 Updated /usr/local/data/sfomuseum-data-architecture/data/172/977/141/3/1729771413.geojson
2021/02/05 14:43:51 Updated /usr/local/data/sfomuseum-data-architecture/data/172/977/142/3/1729771423.geojson
2021/02/05 14:43:51 Updated /usr/local/data/sfomuseum-data-architecture/data/172/977/142/7/1729771427.geojson
2021/02/05 14:43:51 Updated /usr/local/data/sfomuseum-data-architecture/data/172/977/142/9/1729771429.geojson
2021/02/05 14:43:51 Updated /usr/local/data/sfomuseum-data-architecture/data/172/977/143/1/1729771431.geojson
2021/02/05 14:43:51 Updated /usr/local/data/sfomuseum-data-architecture/data/172/977/143/3/1729771433.geojson
2021/02/05 14:43:51 Updated /usr/local/data/sfomuseum-data-architecture/data/172/977/143/5/1729771435.geojson
2021/02/05 14:43:51 Updated /usr/local/data/sfomuseum-data-architecture/data/172/977/141/9/1729771419.geojson
2021/02/05 14:43:51 Updated /usr/local/data/sfomuseum-data-architecture/data/172/977/141/7/1729771417.geojson
2021/02/05 14:43:51 Updated /usr/local/data/sfomuseum-data-architecture/data/172/977/145/1/1729771451.geojson
2021/02/05 14:43:51 Updated /usr/local/data/sfomuseum-data-architecture/data/172/977/143/9/1729771439.geojson
2021/02/05 14:43:51 Updated /usr/local/data/sfomuseum-data-architecture/data/172/977/146/3/1729771463.geojson
2021/02/05 14:43:51 Updated /usr/local/data/sfomuseum-data-architecture/data/172/977/143/7/1729771437.geojson
2021/02/05 14:43:51 Updated /usr/local/data/sfomuseum-data-architecture/data/172/977/144/1/1729771441.geojson
2021/02/05 14:43:51 Updated /usr/local/data/sfomuseum-data-architecture/data/172/977/145/3/1729771453.geojson
2021/02/05 14:43:51 Updated /usr/local/data/sfomuseum-data-architecture/data/172/977/144/5/1729771445.geojson
2021/02/05 14:43:51 Updated /usr/local/data/sfomuseum-data-architecture/data/172/977/146/5/1729771465.geojson
2021/02/05 14:43:51 Updated /usr/local/data/sfomuseum-data-architecture/data/172/977/146/7/1729771467.geojson
2021/02/05 14:43:51 Updated /usr/local/data/sfomuseum-data-architecture/data/172/977/145/5/1729771455.geojson
2021/02/05 14:43:51 Updated /usr/local/data/sfomuseum-data-architecture/data/172/977/144/7/1729771447.geojson
2021/02/05 14:43:51 Updated /usr/local/data/sfomuseum-data-architecture/data/172/977/145/7/1729771457.geojson
2021/02/05 14:43:51 Updated /usr/local/data/sfomuseum-data-architecture/data/172/977/146/9/1729771469.geojson
2021/02/05 14:43:51 Updated /usr/local/data/sfomuseum-data-architecture/data/172/977/144/9/1729771449.geojson
2021/02/05 14:43:51 Updated /usr/local/data/sfomuseum-data-architecture/data/172/977/145/9/1729771459.geojson
2021/02/05 14:43:51 Updated /usr/local/data/sfomuseum-data-architecture/data/172/977/148/3/1729771483.geojson
2021/02/05 14:43:51 Updated /usr/local/data/sfomuseum-data-architecture/data/172/977/147/1/1729771471.geojson
2021/02/05 14:43:51 Updated /usr/local/data/sfomuseum-data-architecture/data/172/977/148/1/1729771481.geojson
2021/02/05 14:43:51 Updated /usr/local/data/sfomuseum-data-architecture/data/172/977/147/3/1729771473.geojson
2021/02/05 14:43:51 Updated /usr/local/data/sfomuseum-data-architecture/data/172/977/147/5/1729771475.geojson
2021/02/05 14:43:51 Updated /usr/local/data/sfomuseum-data-architecture/data/172/977/148/5/1729771485.geojson
2021/02/05 14:43:51 Updated /usr/local/data/sfomuseum-data-architecture/data/172/977/147/7/1729771477.geojson
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