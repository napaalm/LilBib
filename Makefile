BINARY := lilbib
VERSION ?= $(shell git describe --always --dirty --tags 2> /dev/null)
PLATFORMS := windows linux
os = $(word 1, $@)

.PHONY: all
all: clean release run

build:
	mkdir -p build
	go build -o build/$(BINARY) -v

.PHONY: run
run: build
	./build/lilbib

.PHONY: release
release: windows linux

.PHONY: clean
clean:
	go clean
	rm -rf build release

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	mkdir -p release
	GOOS=$(os) GOARCH=amd64 go build -o release/$(BINARY)-$(VERSION)-$(os)-amd64
