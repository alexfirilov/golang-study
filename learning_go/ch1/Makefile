.DEFAULT_GOAL := clean

.PHONY: fmt vet build clean

fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build: vet
	go build -o hello_world

clean: build
	go clean -i -r -cache
	rm -f hello_world
