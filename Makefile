cli:
	go build -mod vendor -o bin/assign-geometry cmd/assign-geometry/main.go
	go build -mod vendor -o bin/assign-parent cmd/assign-parent/main.go
	go build -mod vendor -o bin/exportify cmd/exportify/main.go
	go build -mod vendor -o bin/create cmd/create/main.go
	go build -mod vendor -o bin/deprecate cmd/deprecate/main.go
	go build -mod vendor -o bin/cessate cmd/cessate/main.go
	go build -mod vendor -o bin/ensure-properties cmd/ensure-properties/main.go
	go build -mod vendor -o bin/deprecate-and-supersede cmd/deprecate-and-supersede/main.go
	go build -mod vendor -o bin/merge-feature-collection cmd/merge-feature-collection/main.go
	go build -mod vendor -o bin/supersede-with-parent cmd/supersede-with-parent/main.go
	go build -mod vendor -o bin/as-featurecollection cmd/as-featurecollection/main.go
	go build -mod vendor -o bin/rename-property cmd/rename-property/main.go
