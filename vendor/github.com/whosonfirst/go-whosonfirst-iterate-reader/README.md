# go-whosonfirst-iterate-reader

Go package implementing go-whosonfirst-iterate/emitter functionality using whosonfirst/go-reader.Reader instances.

## Documentation

Documentation is incomplete.

## Example

```
> go run -mod vendor cmd/emit/main.go \
	-emitter-uri 'reader://?reader=fs:///usr/local/data/sfomuseum-data-whosonfirst/data' 1 \

| jq '.["properties"]["wof:name"]'

"Null Island"
```

## See also

* https://github.com/whosonfirst/go-whosonfirst-iterate
* https://github.com/whosonfirst/go-reader 