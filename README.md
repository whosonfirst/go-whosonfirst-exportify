# go-whosonfirst-exportify

Tools (written in Go) for exporting Who's On First records.

## Important

This is work in progress. Documentation is incomplete.

## Tools

```
> make cli
go build -mod vendor -o bin/exportify cmd/exportify/main.go
go build -mod vendor -o bin/ensure-properties cmd/ensure-properties/main.go
go build -mod vendor -o bin/deprecate-and-supersede cmd/deprecate-and-supersede/main.go
go build -mod vendor -o bin/merge-feature-collection cmd/merge-feature-collection/main.go
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

#### Notes

This tool will:

* Copy all the properties from the source WOF record in to the new WOF record.
* Create a new `wof:id` property for the new WOF record.
* Update the `wof:superseded_by` and `wof:supersedes` properties for the old and new WOF records respectively.
* Set the `mz:is_current` property to be "0" for the old WOF record.
* Set the `edtf:deprecate` property to be the current "YYYY-MM-DD" for the old WOF record.
* Assign or update any properties defined by the `-{string|int|float}-properties` flags for the new WOF record.

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
  -query value
    	One or more {PATH}={REGEXP} parameters for filtering records.
  -query-mode string
    	Specify how query filtering should be evaluated. Valid modes are: ALL, ANY (default "ALL")
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
...and so on
```

#### Inline queries

For example

```
$> ./bin/ensure-properties \
	-writer-uri fs:///usr/local/data/sfomuseum-data-architecture/data \
	-query 'properties.mz:is_current=1' \
	-query 'properties.sfomuseum:placetype=gallery' \
	-int-property 'properties.sfo:level=2' \
	/usr/local/data/sfomuseum-data-architecture/

2021/02/10 14:51:59 Updated /usr/local/data/sfomuseum-data-architecture/data/147/785/566/1/1477855661.geojson
2021/02/10 14:51:59 Updated /usr/local/data/sfomuseum-data-architecture/data/147/785/566/3/1477855663.geojson
2021/02/10 14:51:59 Updated /usr/local/data/sfomuseum-data-architecture/data/147/785/582/1/1477855821.geojson
2021/02/10 14:51:59 Updated /usr/local/data/sfomuseum-data-architecture/data/147/785/582/3/1477855823.geojson
2021/02/10 14:51:59 Updated /usr/local/data/sfomuseum-data-architecture/data/147/785/583/1/1477855831.geojson
2021/02/10 14:51:59 Updated /usr/local/data/sfomuseum-data-architecture/data/147/785/593/9/1477855939.geojson
2021/02/10 14:51:59 Updated /usr/local/data/sfomuseum-data-architecture/data/147/785/594/1/1477855941.geojson
2021/02/10 14:51:59 Updated /usr/local/data/sfomuseum-data-architecture/data/147/785/594/3/1477855943.geojson
2021/02/10 14:51:59 Updated /usr/local/data/sfomuseum-data-architecture/data/147/785/594/5/1477855945.geojson
...and so on
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

### merge-feature-collection

```
$> ./bin/merge-feature-collection -h
Upate one or more Who's On First records with matching entries in a GeoJSON FeatureCollection file.

Usage:
	 ./bin/merge-feature-collection [options] path(N) path(N)

For example:
	./bin/merge-feature-collection -reader-uri fs:///usr/local/data/whosonfirst-data-admin-ca/data -writer-uri fs:///usr/local/data/whosonfirst-data-admin-ca/data -path geometry -path 'properties.example:property' /usr/local/data/updates.geojson

Valid options are:
  -exporter-uri string
    	A valid whosonfirst/go-whosonfirst-export URI (default "whosonfirst://")
  -path value
    	One or more valid tidwall/gjson paths. These will be copied from the source GeoJSON feature to the corresponding WOF record.
  -reader-uri string
    	A valid whosonfirst/go-reader URI
  -writer-uri string
    	A valid whosonfirst/go-writer URI
```

For example:

```
$> bin/merge-feature-collection \
	-reader-uri fs:///usr/local/data/sfomuseum-data-publicart/data \
	-writer-uri fs:///usr/local/data/sfomuseum-data-publicart/data \
	-path 'geometry' \
	-path 'properties.sfo:level' \
	/usr/local/sfomuseum/go-sfomuseum-gis/data/publicart-hotel.geojson
```

## See also

* https://github.com/whosonfirst/go-whosonfirst-export
* https://github.com/whosonfirst/go-reader
* https://github.com/whosonfirst/go-writer