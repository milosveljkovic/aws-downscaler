.PHONY: run test build

run:
	go run cmd/aws-downscaler/main.go

test:
	go test -v -timeout 30s ./...

# todo: build multi platform bins
# build:
# 	./build.sh