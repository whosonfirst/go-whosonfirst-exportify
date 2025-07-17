# go-whosonfirst-exportify

Tools (written in Go) for exporting Who's On First (WOF) records.

## Important

This package is deprecated. Short of really serious bugs it is no longer being updated and its functionality is gradually being moved in to the [whosonfirst/wof-cli](https://github.com/whosonfirst/wof-cli) package.

## go-whosonfirst-exportify versus go-whosonfirst-export

* `go-whosonfirst-export` is an abstract package library containing only those dependencies necessary for exporting WOF records.
* `go-whosonfirst-exportify` defines a variety of _applications_ that perform operations relating to, or involving, exporting WOF records. These applications use `go-whosonfirst-export`.

At some point the various application might get separated out in to their own packages but for now they are all bundled together which means this package has, potentially, a lot of dependencies.

## Tools

```
$> make cli
go build -mod vendor -ldflags="-s -w" -o bin/wof-assign-geometry cmd/wof-assign-geometry/main.go
go build -mod vendor -ldflags="-s -w" -o bin/wof-assign-parent cmd/wof-assign-parent/main.go
go build -mod vendor -ldflags="-s -w" -o bin/wof-exportify cmd/wof-exportify/main.go
go build -mod vendor -ldflags="-s -w" -o bin/wof-create cmd/wof-create/main.go
go build -mod vendor -ldflags="-s -w" -o bin/wof-deprecate cmd/wof-deprecate/main.go
go build -mod vendor -ldflags="-s -w" -o bin/wof-cessate cmd/wof-cessate/main.go
go build -mod vendor -ldflags="-s -w" -o bin/wof-superseded-by cmd/wof-superseded-by/main.go
go build -mod vendor -ldflags="-s -w" -o bin/wof-ensure-properties cmd/wof-ensure-properties/main.go
go build -mod vendor -ldflags="-s -w" -o bin/wof-deprecate-and-supersede cmd/wof-deprecate-and-supersede/main.go
go build -mod vendor -ldflags="-s -w" -o bin/wof-merge-csv cmd/wof-merge-csv/main.go
go build -mod vendor -ldflags="-s -w" -o bin/wof-merge-featurecollection cmd/wof-merge-featurecollection/main.go
go build -mod vendor -ldflags="-s -w" -o bin/wof-supersede-with-parent cmd/wof-supersede-with-parent/main.go
go build -mod vendor -ldflags="-s -w" -o bin/wof-as-featurecollection cmd/wof-as-featurecollection/main.go
go build -mod vendor -ldflags="-s -w" -o bin/wof-as-csv cmd/wof-as-csv/main.go
go build -mod vendor -ldflags="-s -w" -o bin/wof-as-jsonl cmd/wof-as-jsonl/main.go
go build -mod vendor -ldflags="-s -w" -o bin/wof-rename-property cmd/wof-rename-property/main.go
go build -mod vendor -ldflags="-s -w" -o bin/wof-remove-properties cmd/wof-remove-properties/main.go
go build -mod vendor -ldflags="-s -w" -o bin/wof-clone-feature cmd/wof-clone-feature/main.go
```

As of this writing these tools may contain duplicate, or at least common, code that would be well-served from being moved in to a package or library. That hasn't happened yet.

### wof-as-csv

Export one or more WOF records as a CSV document written to `STDOUT`.

```
$> ./bin/wof-as-csv -h
Export one or more WOF records as a CSV document written to STDOUT

Usage:
	 ./bin/wof-as-csv [options] path-(N) path-(N)

For example:
	./bin/wof-as-csv -field wof:id -field wof:name -field centroid -iterator-uri 'repo://?include=properties.mz:is_current=1' /usr/local/data/sfomuseum-data-publicart/
Valid options are:
  -field value
    	One or more relative 'properties.FIELDNAME' paths to include the CSV output. If the fieldname is 'path' the filename of the current record will be included. If the fieldname is 'centroid' the primary centroid of the current record will be derived and included as 'latitude' and 'longitude' columns.
  -iterator-uri string
    	 (default "repo://")
```

For example:

```
$> ./bin/wof-as-csv  \
	-field wof:id \
	-field centroid \
	-field wof:name \
	-iterator-uri 'repo://?include=properties.sfomuseum:placetype=boardingarea&include=properties.mz:is_current=1' \
	/usr/local/data/sfomuseum-data-architecture/
	
wof:id,latitude,longitude,wof:name
1763588233,37.615069504933075,-122.38306184920478,Boarding Area C
1763588177,37.612307411494186,-122.38518056697221,Boarding Area B
1763588335,37.61888020848815,-122.38446905789152,Boarding Area E
1763588433,37.617858008730636,-122.39137174044615,Boarding Area G
1763588125,37.61734905835808,-122.38196681238652,Boarding Area D
1763588271,37.62010253190792,-122.38804674158527,Boarding Area F
1763588371,37.61289407736701,-122.38934858141738,Boarding Area A
```

### wof-as-featurecollection

**THIS TOOL IS DEPRECATED and is no longer being updated. It has been replaced by https://github.com/whosonfirst/wof-cli/tree/main?tab=readme-ov-file#wof-emit**

Export one or more WOF records as a GeoJSON FeatureCollection

```
$> ./bin/wof-as-featurecollection -h
Export one or more WOF records as a GeoJSON FeatureCollection

Usage:
	 ./bin/wof-as-featurecollection [options] path-(N) path-(N)

For example:
  -as-multipoints
    	Output geometries as a MultiPoint array
	./bin/as-featurecollection -iterator-uri 'repo://?include=properties.mz:is_current=1' /usr/local/data/sfomuseum-data-publicart/
Valid options are:
  -iterator-uri string
    	A valid whosonfirst/go-whosonfirst-iterator/emitter URI. Supported emitter URI schemes are: directory://,featurecollection://,file://,filelist://,geojsonl://,repo:// (default "repo://")
  -writer-uri string
    	A valid whosonfirst/go-writer URI. Supported writer URI schemes are: file://, fs://, io://, null://, stdout:// (default "stdout://")
```

For example:

```
$> ./bin/wof-as-featurecollection \
	-iterator-uri 'repo://?include=properties.mz:is_current=1' \
	/usr/local/data/sfomuseum-data-publicart/ \
| jq '.features[]["properties"]["wof:parent_id"]' \
| sort \
| uniq

1477855657
1477855669
1477855979
1477855987
1477856005
1729791967
1729792389
1729792391
1729792433
1729792437
1729792459
1729792483
1729792489
1729792551
1729792577
1729792581
1729792643
1729792645
1729792679
1729792685
1729792689
1729792691
1729792693
1729792695
1729792699
```

### wof-as-jsonl

**THIS TOOL IS DEPRECATED and is no longer being updated. It has been replaced by https://github.com/whosonfirst/wof-cli/tree/main?tab=readme-ov-file#wof-emit**

Export one or more WOF records as a line-separated JSON.

```
$> ./bin/wof-as-jsonl -h
Export one or more WOF records as a line-separated JSON

Usage:
	 ./bin/wof-as-jsonl [options] path-(N) path-(N)

For example:
	./bin/wof-as-jsonl -iterator-uri 'repo://?include=properties.mz:is_current=1' /usr/local/data/sfomuseum-data-publicart/
Valid options are:
  -as-multipoints
    	Output geometries as a MultiPoint array
  -iterator-uri string
    	A valid whosonfirst/go-whosonfirst-iterator/v2 URI. Supported emitter URI schemes are: directory://,featurecollection://,file://,filelist://,geojsonl://,null://,repo:// (default "repo://")
  -writer-uri string
    	A valid whosonfirst/go-writer URI. Supported writer URI schemes are: fs://, io://, null://, stdout:// (default "stdout://")
```

For example:

```
$> ./bin/wof-as-jsonl \
	-iterator-uri 'repo://?include=properties.wof:placetype=timezone' \
	/usr/local/data/whosonfirst-data-admin-xy \
| wc -l

2021/10/14 09:54:07 time to index paths (1) 2.475059637s
     376
```

### assign-geometry

**THIS TOOL IS DEPRECATED and is no longer being updated. It has been replaced by https://github.com/whosonfirst/wof-cli/tree/main?tab=readme-ov-file#wof-geometry**

Assign the geometry from a given record to one or more other records.

```
$> ./bin/wof-assign-geometry -h
Assign the geometry from a given record to one or more other records.

Usage:
	 ./bin/wof-assign-geometry [options] target-id-(N) target-id-(N)

For example:
	./bin/wof-assign-geometry -reader-uri fs:///usr/local/data/whosonfirst-data-admin-ca/data -source-id 1234 5678

Valid options are:
  -exporter-uri string
    	A valid whosonfirst/go-whosonfirst-export URI (default "whosonfirst://")
  -reader-uri string
    	A valid whosonfirst/go-reader URI.
  -source-id int
    	A valid Who's On First ID.
  -stdin
    	Read target IDs from STDIN
  -writer-uri string
    	A valid whosonfirst/go-writer URI. (default "stdout://")
```

_This tool can (and probably be should) be adapted in the general-purpose assign properties or geometry tool. Today it only handles geometries._

For example:

```
$> /usr/local/whosonfirst/go-whosonfirst-travel/bin/wof-travel-id \
	-supersedes \
	-ids \
	-source fs:///usr/local/data/sfomuseum-data-architecture/data \
	1729813675 \
	
| bin/wof-assign-geometry \
	-stdin \
	-reader-uri fs:///usr/local/data/sfomuseum-data-architecture/data \
	-writer-uri fs:///usr/local/data/sfomuseum-data-architecture/data \
	-source-id 1729813675
```

This is a fairly involved example, so here's a description of what's happening:

* First we are using the `wof-travel-id` tool in the [go-whosonfirst-travel](https://github.com/whosonfirst/go-whosonfirst-travel) package to find all the records that, recursively, ID `1729813675` supersedes.
* We are exporting that list as a line-separated list of IDs (having specified the `-ids` flag).
* We are piping that output to the `wof-assign-geometry` tool and reading the list of IDs from `STDIN` (having specified the `-stdin` flag).
* For each of those IDs we are assigning the geometry from the record with the ID `1729813675` and writing those updates back to disk.

### wof-assign-parent

```
$> ./bin/wof-assign-parent -h
Assign the parent ID and its hierarchy to one or more WOF records

Usage:
	 ./bin/wof-assign-parent [options] wof-id-(N) wof-id-(N)

For example:
Valid options are:
  -exporter-uri string
    	A valid whosonfirst/go-whosonfirst-export URI. (default "whosonfirst://")
  -id value
    	One or more valid Who's On First ID.
  -parent-id int
    	A valid Who's On First ID.
  -parent-reader-uri string
    	A valid whosonfirst/go-reader URI. If empty the value of the -reader-uri flag will be assumed.
  -reader-uri string
    	A valid whosonfirst/go-reader URI.
  -writer-uri string
    	A valid whosonfirst/go-writer URI. If empty the value of the -reader-uri flag will be assumed.
```

For example:

```
$> ./bin/wof-assign-parent \
	-reader-uri fs:///usr/local/data/sfomuseum-data-architecture/data \
	-parent-id 1477855937 \
	-id 1477855939 -id 1477855941 -id 1477855943 -id 1477855945 -id 1477855947 -id 1477855949 1477855955
```

### wof-cessate

"Cessate" one or more Who's On First IDs (assign an `edtf:cessation` property and assign `mz:is_current=0`).

```
> ./bin/wof-cessate -h
"Cessate" one or more Who's On First IDs.

Usage:
	 ./bin/wof-cessate [options]

For example:
	./bin/wof-cessate -s . -i 1234
	./bin/wof-cessate -reader-uri fs:///usr/local/data/whosonfirst-data-admin-ca/data -id 1234 -id 5678

Valid options are:
  -date string
    	A valid EDTF date. If empty then the current date will be used
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
  -supersede-with-copy
    	Supersede this record with a copy of itself.	
  -superseded-by value
    	Zero or more Who's On First IDs that the records being deprecated are superseded by.
  -writer-uri string
    	A valid whosonfirst/go-writer URI. If empty the value of the -s flag will be used in combination with the fs:// scheme.
```

### wof-clone-feature

Clone and optionally supersede a Who's On First record.

```
> ./bin/wof-clone-feature -h
Clone and optionally supersede a Who's On First record.

Usage:
	 ./bin/wof-clone-feature [options]

For example:
	./bin/wof-clone-feature -s . -id 1234 -superseded
	./bin/wof-clone-feature -reader-uri fs:///usr/local/data/whosonfirst-data-admin-ca/data -id 1234 -superseded

Valid options are:
  -exporter-uri string
    	A valid whosonfirst/go-whosonfirst-export URI. (default "whosonfirst://")
  -float-property value
    	One or more {KEY}={VALUE} properties to append to the new record where {KEY} is a valid tidwall/gjson path and {VALUE} is a float(64) value.
  -id int
    	The feature being cloned.
  -int-property value
    	One or more {KEY}={VALUE} properties to append to the new record where {KEY} is a valid tidwall/gjson path and {VALUE} is a int(64) value.
  -reader-uri string
    	A valid whosonfirst/go-reader URI. If empty the value of the -s flag will be used in combination with the fs:// scheme.
  -s string
    	A valid path to the root directory of the Who's On First data repository. If empty (and -reader-uri or -writer-uri are empty) the current working directory will be used and appended with a 'data' subdirectory.
  -string-property value
    	One or more {KEY}={VALUE} properties to append to the new record where {KEY} is a valid tidwall/gjson path and {VALUE} is a string value.
  -superseded
    	The new feature is superseded by the feature being cloned.
  -supersedes
    	The new feature supersedes the feature being cloned.
  -writer-uri string
    	A valid whosonfirst/go-writer URI. If empty the value of the -s flag will be used in combination with the fs:// scheme.
```

Consider this passage from the Wikipedia entry for [O. R. Tambo International Airport](https://en.wikipedia.org/wiki/O._R._Tambo_International_Airport) (WOF ID [102546665](https://spelunker.whosonfirst.org/id/102546665/)):

> It was originally known as Jan Smuts International Airport,[4] after the former South African Prime Minister of the same name. The airport was renamed Johannesburg International Airport in 1994 when the newly elected African National Congress (ANC) government implemented a policy of not naming airports after politicians. This policy was later reversed, and on 27 October 2006 the airport was renamed after anti-apartheid politician Oliver Reginald Tambo.[5]

So given an initial record (`102546665`) we want to create two others. First we create the record for "Johannesburg International Airport", which is superseded by the record for "OR Tambo International Airport":

```
$> ./bin/wof-clone-feature \
	-s /usr/local/data/whosonfirst-data-admin-za \
	-id 102546665 \
	-superseded \
	-string-property 'properties.wof:name=Johannesburg International Airport' \
	-string-property 'properties.edtf:inception=1994' \
	-string-property 'properties.edtf:cessation=2006'
```

Which creates a new record with ID `1746875291`. Note that because this new record was _cloned_ it may still need to be manually updated to reflect its reality, for example `name:*_*` properties may need to be removed or changed.

Next, we create a record for "Jan Smuts International Airport" which is cloned from the newly created record for "Johannesburg International Airport":

```
$> ./bin/wof-clone-feature \
	-s /usr/local/data/whosonfirst-data-admin-za \
	-id 1746875291 \
	-supersedes \
	-string-property 'properties.wof:name=Jan Smuts International Airport' \
	-string-property 'properties.edtf:inception=1952' \
	-string-property 'properties.edtf:cessation=1994'
```	

### wof-create

Create a new Who's On First record.

```
$> ./bin/wof-create -h
Create a new Who's On First record.

Usage:
	 ./bin/wof-create [options] 

Valid options are:
  -exporter-uri string
    	A valid whosonfirst/go-whosonfirst-export URI. (default "whosonfirst://")
  -float-property value
    	One or more {KEY}={VALUE} flags where {KEY} is a valid tidwall/gjson path and {VALUE} is a float(64) value.
  -geometry string
    	A valid GeoJSON geometry
  -int-property value
    	One or more {KEY}={VALUE} flags where {KEY} is a valid tidwall/gjson path and {VALUE} is a int(64) value.
  -parent-reader-uri string
    	A valid whosonfirst/go-reader URI. If empty the value of the -s fs will be used in combination with the fs:// scheme.
  -resolve-hierarchy
    	Attempt to resolve parent ID and hierarchy using point-in-polygon lookups. If true the -spatial-database-uri flag must also be set
  -s string
    	A valid path to the root directory of the Who's On First data repository. If empty (and -reader-uri or -writer-uri are empty) the current working directory will be used and appended with a 'data' subdirectory.
  -spatial-database-uri string
    	A valid whosonfirst/go-whosonfirst-spatial/database URI.
  -string-property value
    	One or more {KEY}={VALUE} flags where {KEY} is a valid tidwall/gjson path and {VALUE} is a string value.
  -writer-uri string
    	A valid whosonfirst/go-writer URI. If empty the value of the -s fs will be used in combination with the fs:// scheme.
```

_This tool should be considered "beta" still._

For example:

```
$> bin/wof-create \
	-geometry '{"type":"Point", "coordinates":[20.414944,42.032833]}' \
	-string-property 'properties.src:geom=wikidata' \
	-string-property 'properties.wof:placetype=campus' \
	-string-property 'properties.wof:placetype_alt=airport' \
	-string-property 'properties.wof:repo=whosonfirst-data-admin-al' \
	-string-property 'properties.wof:name=Kukës International Airport' \
	-string-property 'properties.edtf:inception=2021-04-18' \
	-string-property 'properties.edtf:cessation=..' \
	-string-property 'properties.wof:concordances.icao:code=LAKU' \
	-string-property 'properties.wof:concordances.wk:id=Q1431804' \
	-int-property 'properties.mz:is_current=1' \
	-resolve-hierarchy \
	-spatial-database-uri 'sqlite://?dsn=/usr/local/data/al.db' \
	-writer-uri stdout://

{
  "id": 1730032323,
  "type": "Feature",
  "properties": {
    "date:inception_lower": "2021-04-18",
    "date:inception_upper": "2021-04-18",
    "edtf:cessation": "..",
    "edtf:inception": "2021-04-18",
    "geom:area": 0,
    "geom:bbox": "20.414944,42.032833,20.414944,42.032833",
    "geom:latitude": 42.032833,
    "geom:longitude": 20.414944,
    "src:geom": "wikidata",
    "wof:belongsto": [
      102191581,
      85632405,
      421186339,
      85667797
    ],
    "wof:concordances": {
      "icao:code": "LAKU"
    },
    "wof:country": "AL",
    "wof:created": 1619808522,
    "wof:geomhash": "792d95b83651dcc4aedeb0923d9c05f8",
    "wof:hierarchy": [
      {
        "campus_id": 1730032323,
        "continent_id": 102191581,
        "country_id": 85632405,
        "county_id": 421186339,
        "region_id": 85667797
      }
    ],
    "wof:id": 1730032323,
    "wof:lastmodified": 1619808522,
    "wof:name": "Kukës International Airport",
    "wof:parent_id": 421186339,
    "wof:placetype": "campus",
    "wof:placetype_alt": "airport",
    "wof:repo": "whosonfirst-data-admin-al",
    "wof:superseded_by": [],
    "wof:supersedes": []
  },
  "bbox": [
    20.414944,
    42.032833,
    20.414944,
    42.032833
  ],
  "geometry": {"coordinates":[20.414944,42.032833],"type":"Point"}
}
1730032323
```

### wof-deprecate

Deprecate one or more Who's On First IDs.

```
> ./bin/wof-deprecate -h
Deprecate one or more Who's On First IDs.

Usage:
	 ./bin/wof-deprecate [options]

For example:
	./bin/wof-deprecate -s . -i 1234
	./bin/wof-deprecate -reader-uri fs:///usr/local/data/whosonfirst-data-admin-ca/data -id 1234 -id 5678

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
  -superseded-by value
    	Zero or more Who's On First IDs that the records being deprecated are superseded by.
  -writer-uri string
    	A valid whosonfirst/go-writer URI. If empty the value of the -s flag will be used in combination with the fs:// scheme.
```

For example:

```
$> ./bin/wof-deprecate -s /usr/local/data/sfomuseum-data-collection -i 1511957049

$> git diff
diff --git a/data/151/195/704/9/1511957049.geojson b/data/151/195/704/9/1511957049.geojson
index 916115027f..ec21ca2089 100644
--- a/data/151/195/704/9/1511957049.geojson
+++ b/data/151/195/704/9/1511957049.geojson
@@ -8,6 +8,7 @@
     "date:inception_upper": "1960-01-01",
     "edtf:cessation": "1960-12-31",
     "edtf:date": "1960",
+    "edtf:deprecated": "2021-05-03",
     "edtf:inception": "1960-01-01",
     "geom:area": 0,
     "geom:bbox": "-122.386155,37.616358,-122.386155,37.616358",
@@ -18,7 +19,7 @@
     "millsfield:collection_id": 1511214203,
     "millsfield:subcategory_id": 1511213595,
     "mz:hierarchy_label": 1,
-    "mz:is_current": -1,
+    "mz:is_current": 0,
     "sfomuseum:accession_number": "2012.147.452",
     "sfomuseum:category": "Insignia",
     "sfomuseum:collection": "Aviation Museum",
@@ -65,7 +66,7 @@
       }
     ],
     "wof:id": 1511957049,
-    "wof:lastmodified": 1618962263,
+    "wof:lastmodified": 1620087808,
     "wof:name": "flight officer cap badge: Iraq Petroleum Company",
     "wof:parent_id": 1511214277,
     "wof:placetype": "venue",
```

### wof-deprecate-and-supersede

```
$> ./bin/wof-deprecate-and-supersede -h
Deprecate and supersede one or more Who's On First IDs.

Usage:
	 ./bin/wof-deprecate-and-supersede [options] wof-id-(N) wof-id-(N)

For example:
	./bin/wof-deprecate-and-supersede -s . -i 1234
	./bin/wof-deprecate-and-supersede -reader-uri fs:///usr/local/data/whosonfirst-data-admin-ca/data -id 1234 -id 5678

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
$> ./bin/wof-deprecate-and-supersede \
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

### wof-ensure-properties

**THIS TOOL IS DEPRECATED and is no longer being updated. It has been replaced by https://github.com/whosonfirst/wof-cli/tree/main?tab=readme-ov-file#wof-ensure-property**

```
> ./bin/wof-ensure-properties -h
Usage of ./bin/wof-ensure-properties:
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
$> ./bin/wof-ensure-properties \
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
$> ./bin/wof-ensure-properties \
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

### wof-exportify

**THIS TOOL IS DEPRECATED and is no longer being updated. It has been replaced by https://github.com/whosonfirst/wof-cli/tree/main?tab=readme-ov-file#wof-export**

```
$> ./bin/wof-exportify -h
Exportify one or more Who's On First IDs.

Usage:
	 ./bin/wof-exportify [options] wof-id-(N) wof-id-(N)

For example:
	./bin/wof-exportify -s . -i 1234
	./bin/wof-exportify -reader-uri fs:///usr/local/data/whosonfirst-data-admin-ca/data -id 1234 -id 5678

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

### wof-merge-csv

```
$> ./bin/wof-merge-csv -h
Upate one or more Who's On First records with matching entries in a CSV file.

Usage:
	 ./bin/wof-merge-csv [options] path(N) path(N)

For example:
	./bin/wof-merge-csv -reader-uri repo:///usr/local/data/sfomuseum-data-architecture -writer-uri repo:///usr/local/data/sfomuseum-data-architecture -int-field sfo:level galleries-with-level.csv

Valid options are:
  -exporter-uri string
    	A valid whosonfirst/go-whosonfirst-export URI (default "whosonfirst://")
  -int-field value
    	Zero or more fields in a CSV row to assign to a WOF record as int values.
  -int64-field value
    	Zero or more fields in a CSV row to assign to a WOF record as int64 values.
  -lookup-key string
    	The column in a CSV row to use to lookup a corresponding Who's On First record. (default "wof:id")
  -reader-uri string
    	A valid whosonfirst/go-reader URI
  -string-field value
    	Zero or more fields in a CSV row to assign to a WOF record as string values.
  -writer-uri string
    	A valid whosonfirst/go-writer URI
```

### wof-merge-featurecollection

```
$> ./bin/wof-merge-featurecollection -h
Upate one or more Who's On First records with matching entries in a GeoJSON FeatureCollection file.

Usage:
	 ./bin/wof-merge-featurecollection [options] path(N) path(N)

For example:
	./bin/wof-merge-featurecollection -reader-uri fs:///usr/local/data/whosonfirst-data-admin-ca/data -writer-uri fs:///usr/local/data/whosonfirst-data-admin-ca/data -path geometry -path 'properties.example:property' /usr/local/data/updates.geojson

Valid options are:
  -exporter-uri string
    	A valid whosonfirst/go-whosonfirst-export URI (default "whosonfirst://")
  -include value
    	One or more {PATH}={REGEXP} parameters for filtering records when building a lookup map.
  -include-mode string
    	Specify how query filtering should be evaluated. Valid modes are: ALL, ANY (default "ALL")
  -lookup-key string
    	A valid tidwall/gjson path to use for specifying an alternative (to 'properties.wof:id') lookup key. The value of this key will be mapped to the record's 'wof:id' property.
  -lookup-mode string
    	A valid whosonfirst/go-whosonfirst-index URI. (default "repo://")
  -lookup-source value
    	One or more valid whosonfirst/go-whosonfirst-index sources.
  -path value
    	One or more valid tidwall/gjson paths. These will be copied from the source GeoJSON feature to the corresponding WOF record.
  -reader-uri string
    	A valid whosonfirst/go-reader URI
  -writer-uri string
    	A valid whosonfirst/go-writer URI
```

For example:

```
$> bin/wof-merge-featurecollection \
	-reader-uri fs:///usr/local/data/sfomuseum-data-publicart/data \
	-writer-uri fs:///usr/local/data/sfomuseum-data-publicart/data \
	-path 'geometry' \
	-path 'properties.sfo:level' \
	/usr/local/sfomuseum/go-sfomuseum-gis/data/publicart-hotel.geojson
```

Here's a more complex example:

```
$> ./bin/wof-merge-featurecollection \
	-include 'properties.sfomuseum:placetype=gallery' \
	-include 'properties.mz:is_current=1' \
	-lookup-key 'properties.sfomuseum:map_id' \
	-lookup-source /usr/local/data/sfomuseum-data-architecture/ \
	-path geometry \
	-reader-uri fs:///usr/local/data/sfomuseum-data-architecture/data \
	-writer-uri fs:///usr/local/data/sfomuseum-data-architecture/data \
	/usr/local/sfomuseum/go-sfomuseum-gis/data/galleries-3.geojson
```

In the example above we are:

* Building a lookup map using records in the `/usr/local/data/sfomuseum-data-architecture/` directory. This lookup map will track a specific value in both the source data and the data being merged to its corresponding WOF ID.
* Only including records with `sfomuseum:placetype=gallery` and `mz:is_current=1` properties.
* Specifying that the lookup key is `properties.sfomuseum:map_id` - this value will be mapped to the corresponding record's `wof:id` property
* Using the lookup key property in the data being merged to determine which WOF record (read by the `-reader-uri` flag) should be updated

### wof-rename-property

Rename a property in one or more records. Currently this tool does not support renaming more than one property at a time.

```
$> ./bin/wof-rename-property -h
Usage of ./bin/wof-rename-property:
  -exporter-uri string
    	A valid whosonfirst/go-whosonfirst-export URI. (default "whosonfirst://")
  -indexer-uri string
    	A valid whosonfirst/go-whosonfirst-iterate/emitter URI. (default "repo://")
  -new-property string
    	The fully qualified path of the property to be (re)named.
  -old-property string
    	The fully qualified path of the property to rename.
  -writer-uri string
    	A valid whosonfirst/go-writer URI. (default "null://")
```

For example:

```
$> ./bin/wof-rename-property \
	-writer-uri fs:///usr/local/data/sfomuseum-data-architecture/data \
	-old-property 'properties.wof:supsersedes' \
	-new-property 'properties.wof:supersedes' \
	/usr/local/data/sfomuseum-data-architecture
	
2021/05/11 13:33:55 Updated /usr/local/data/sfomuseum-data-architecture/data/172/981/373/7/1729813737.geojson
2021/05/11 13:33:55 Updated /usr/local/data/sfomuseum-data-architecture/data/172/981/373/5/1729813735.geojson
2021/05/11 13:33:55 Updated /usr/local/data/sfomuseum-data-architecture/data/172/981/373/3/1729813733.geojson
2021/05/11 13:33:55 Updated /usr/local/data/sfomuseum-data-architecture/data/172/981/373/1/1729813731.geojson
2021/05/11 13:33:55 Updated /usr/local/data/sfomuseum-data-architecture/data/172/981/373/9/1729813739.geojson
2021/05/11 13:33:55 Updated /usr/local/data/sfomuseum-data-architecture/data/172/979/194/7/1729791947.geojson
2021/05/11 13:33:55 Updated /usr/local/data/sfomuseum-data-architecture/data/172/979/194/9/1729791949.geojson
2021/05/11 13:33:55 Updated /usr/local/data/sfomuseum-data-architecture/data/172/981/374/1/1729813741.geojson
2021/05/11 13:33:55 Updated /usr/local/data/sfomuseum-data-architecture/data/172/979/194/3/1729791943.geojson
2021/05/11 13:33:55 Updated /usr/local/data/sfomuseum-data-architecture/data/172/981/374/5/1729813745.geojson
2021/05/11 13:33:55 Updated /usr/local/data/sfomuseum-data-architecture/data/172/981/374/7/1729813747.geojson
2021/05/11 13:33:55 Updated /usr/local/data/sfomuseum-data-architecture/data/172/979/195/7/1729791957.geojson
2021/05/11 13:33:55 Updated /usr/local/data/sfomuseum-data-architecture/data/172/979/195/3/1729791953.geojson
2021/05/11 13:33:55 Updated /usr/local/data/sfomuseum-data-architecture/data/172/979/195/9/1729791959.geojson
2021/05/11 13:33:55 Updated /usr/local/data/sfomuseum-data-architecture/data/172/979/195/5/1729791955.geojson
2021/05/11 13:33:55 Updated /usr/local/data/sfomuseum-data-architecture/data/172/979/195/1/1729791951.geojson
2021/05/11 13:33:55 Updated /usr/local/data/sfomuseum-data-architecture/data/172/979/196/1/1729791961.geojson
2021/05/11 13:33:55 Updated /usr/local/data/sfomuseum-data-architecture/data/172/979/196/5/1729791965.geojson
2021/05/11 13:33:55 Updated /usr/local/data/sfomuseum-data-architecture/data/172/979/196/9/1729791969.geojson
2021/05/11 13:33:55 Updated /usr/local/data/sfomuseum-data-architecture/data/172/979/196/7/1729791967.geojson
2021/05/11 13:33:55 Updated /usr/local/data/sfomuseum-data-architecture/data/172/994/578/5/1729945785.geojson
2021/05/11 13:33:55 time to index paths (1) 180.656951ms
```

### wof-supersede-with-parent

**THIS TOOL IS DEPRECATED and is no longer being updated. It has been replaced by https://github.com/whosonfirst/wof-cli/tree/main?tab=readme-ov-file#wof-supersede**

Supersede one or more WOF records with a known parent ID (and hierarchy).

```
$> ./bin/wof-supersede-with-parent -h
Supersede one or more WOF records with a known parent ID (and hierarchy)

Usage:
	 ./bin/wof-supersede-with-parent [options]

For example:
	./bin/wof-supersede-with-parent -reader-uri fs:///usr/local/data/sfomuseum-data-architecture/data -parent-id 1477855937 -id 1477855955

Valid options are:
  -exporter-uri string
    	A valid whosonfirst/go-whosonfirst-export URI. (default "whosonfirst://")
  -id value
    	One or more valid Who's On First ID.
  -parent-id int
    	A valid Who's On First ID.
  -parent-reader-uri string
    	A valid whosonfirst/go-reader URI. If empty the value of the -reader-uri flag will be assumed.
  -reader-uri string
    	A valid whosonfirst/go-reader URI.
  -writer-uri string
    	A valid whosonfirst/go-writer URI. If empty the value of the -reader-uri flag will be assumed.
```	

### wof-superseded-by

The `wof-superseded-by` tool will update the `wof:superseded_by` and `wof:supersedes` properties for one or more sets of Who's On First records. Additionally it will assign the `mz:is_current=0` property for the records being superseded.

```
> ./bin/wof-superseded-by -h
"Supersede" one or more Who's On First IDs.

Usage:
	 ./bin/wof-superseded-by [options]

For example:
	./bin/wof-superseded-by -reader-uri fs:///usr/local/data/sfomuseum-data-enterprise/data -id 1159286017 -by 1159283849

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
  -superseded-by value
    	Zero or more Who's On First IDs that the records being deprecated are superseded by.
  -writer-uri string
    	A valid whosonfirst/go-writer URI. If empty the value of the -s flag will be used in combination with the fs:// scheme.
```

## See also

* https://github.com/whosonfirst/go-whosonfirst-export
* https://github.com/whosonfirst/go-reader
* https://github.com/whosonfirst/go-writer
