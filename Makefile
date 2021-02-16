cli:
	go build -mod vendor -o bin/assign-parent cmd/assign-parent/main.go
	go build -mod vendor -o bin/exportify cmd/exportify/main.go
	go build -mod vendor -o bin/ensure-properties cmd/ensure-properties/main.go
	go build -mod vendor -o bin/deprecate-and-supersede cmd/deprecate-and-supersede/main.go
	go build -mod vendor -o bin/merge-feature-collection cmd/merge-feature-collection/main.go
	go build -mod vendor -o bin/supersede-with-parent cmd/supersede-with-parent/main.go
