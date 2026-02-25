.PHONY: build test build-plugins run

build:
	go build -o bin/gateway cmd/gateway/main.go

build-plugins:
	GOOS=wasip1 GOARCH=wasm go build -o plugins/example.wasm plugins/example/main.go

test:
	go test ./...

run: build-plugins build
	./bin/gateway
