.PHONY: cli

GOMOD=vendor

cli:
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/wof-assign-geometry cmd/wof-assign-geometry/main.go
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/wof-assign-parent cmd/wof-assign-parent/main.go
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/wof-exportify cmd/wof-exportify/main.go
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/wof-create cmd/wof-create/main.go
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/wof-deprecate cmd/wof-deprecate/main.go
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/wof-cessate cmd/wof-cessate/main.go
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/wof-superseded-by cmd/wof-superseded-by/main.go
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/wof-ensure-properties cmd/wof-ensure-properties/main.go
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/wof-deprecate-and-supersede cmd/wof-deprecate-and-supersede/main.go
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/wof-merge-featurecollection cmd/wof-merge-featurecollection/main.go
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/wof-supersede-with-parent cmd/wof-supersede-with-parent/main.go
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/wof-as-featurecollection cmd/wof-as-featurecollection/main.go
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/wof-as-csv cmd/wof-as-csv/main.go
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/wof-as-jsonl cmd/wof-as-jsonl/main.go
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/wof-rename-property cmd/wof-rename-property/main.go
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/wof-remove-properties cmd/wof-remove-properties/main.go
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/wof-clone-feature cmd/wof-clone-feature/main.go
