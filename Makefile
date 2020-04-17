BINARY := lilbib
VERSION ?= $(shell git describe --always --dirty --tags 2> /dev/null)

.PHONY: all
all: clean release run

.PHONY: build
build:
	go build -o $(BINARY) cmd/lilbib/main.go

.PHONY: run
run: build
	./$(BINARY)

.PHONY: release
release: linux windows

.PHONY: clean
clean:
	rm -rf $(BINARY) release

.PHONY: linux
linux:
	mkdir -p release/$(BINARY)-$(VERSION)-$@-amd64
	cp -r web release/$(BINARY)-$(VERSION)-$@-amd64
	GOOS=$@ GOARCH=amd64 go build -o release/$(BINARY)-$(VERSION)-$@-amd64/ ./...
	tar -czf release/$(BINARY)-$(VERSION)-$@-amd64.tar.gz release/$(BINARY)-$(VERSION)-$@-amd64

.PHONY: windows
windows:
	mkdir -p release/$(BINARY)-$(VERSION)-$@-amd64
	cp -r web release/$(BINARY)-$(VERSION)-$@-amd64
	GOOS=$@ GOARCH=amd64 go build -o release/$(BINARY)-$(VERSION)-$@-amd64/ ./...
	zip -qr release/$(BINARY)-$(VERSION)-$@-amd64.zip release/$(BINARY)-$(VERSION)-$@-amd64
