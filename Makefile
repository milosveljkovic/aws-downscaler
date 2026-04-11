.PHONY: run test build

run:
	go run -ldflags="-X main.version=dev-$(shell date +%Y-%m-%d-%H:%M:%S)" cmd/aws-downscaler/main.go --help

test:
	go test -v -timeout 30s ./...

build:
	go build -o aws-downscaler ./cmd/aws-downscaler
# 	./build.sh (multi platform build)