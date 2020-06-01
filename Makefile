BINARY := lilbib
VERSION ?= $(shell git describe --always --dirty --tags 2> /dev/null)

.PHONY: all
all: clean release run

.PHONY: build
build:
	go build -ldflags "-X main.Version=$(VERSION)" -o $(BINARY) cmd/lilbib/main.go

sandbox/config: config | sandbox/
	cp -r $^ $@

sandbox/web: | sandbox/
	ln -s ../web $@

sandbox/$(BINARY): $(BINARY) | sandbox/
	cp $^ $@

.PHONY: sandbox
sandbox: sandbox/web sandbox/config sandbox/$(BINARY)

.PHONY: run
run: build sandbox
	./sandbox/$(BINARY)

.PHONY: release
release: linux windows

.PHONY: clean
clean:
	rm -rf sandbox $(BINARY) release

.PHONY: linux
linux:
	mkdir -p release/$(BINARY)-$(VERSION)-$@-amd64
	cp -r config release/$(BINARY)-$(VERSION)-$@-amd64
	cp -r web release/$(BINARY)-$(VERSION)-$@-amd64
	GOOS=$@ GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION)" -o release/$(BINARY)-$(VERSION)-$@-amd64/ ./...
	tar -czf release/$(BINARY)-$(VERSION)-$@-amd64.tar.gz release/$(BINARY)-$(VERSION)-$@-amd64

.PHONY: windows
windows:
	mkdir -p release/$(BINARY)-$(VERSION)-$@-amd64
	cp -r config release/$(BINARY)-$(VERSION)-$@-amd64
	cp -r web release/$(BINARY)-$(VERSION)-$@-amd64
	GOOS=$@ GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION)" -o release/$(BINARY)-$(VERSION)-$@-amd64/ ./...
	zip -qr release/$(BINARY)-$(VERSION)-$@-amd64.zip release/$(BINARY)-$(VERSION)-$@-amd64

%/:
	mkdir -p $*
