cli:
	go build -mod vendor -o bin/exportify cmd/exportify/main.go
	go build -mod vendor -o bin/ensure-properties cmd/ensure-properties/main.go
