all: build

test:
	go test -count=1 -v ./...
	go test -count=10 ./...

build: test
	go build -a -ldflags="-s -w -X main.version=$(VERSION)" ./cmd/healthcheck
	go build -a -ldflags="-s -w -X main.version=$(VERSION)" ./cmd/server
