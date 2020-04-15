BINARY := lilbib
VERSION ?= $(shell git describe --always --dirty --tags 2> /dev/null)

.PHONY: all
all: clean release run

.PHONY: build
build:
	go build -o $(BINARY) cmd/lilbib/main.go

.PHONY: run
run: build
	./lilbib

.PHONY: release
release: linux windows

.PHONY: clean
clean:
	rm -f lilbib
	rm -rf release

.PHONY: linux
linux:
	mkdir -p release/$(BINARY)-$(VERSION)-linux-amd64
	cp -r web release/$(BINARY)-$(VERSION)-linux-amd64
	GOOS=linux GOARCH=amd64 go build -o release/$(BINARY)-$(VERSION)-linux-amd64/lilbib cmd/lilbib/main.go
	tar -czf release/$(BINARY)-$(VERSION)-linux-amd64.tar.gz release/$(BINARY)-$(VERSION)-linux-amd64

.PHONY: windows
windows:
	mkdir -p release/$(BINARY)-$(VERSION)-windows-amd64
	cp -r web release/$(BINARY)-$(VERSION)-windows-amd64
	GOOS=linux GOARCH=amd64 go build -o release/$(BINARY)-$(VERSION)-windows-amd64/lilbib.exe cmd/lilbib/main.go
	zip -qr release/$(BINARY)-$(VERSION)-windows-amd64.zip release/$(BINARY)-$(VERSION)-windows-amd64
