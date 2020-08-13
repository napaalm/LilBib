BINARY := lilbib
VERSION ?= $(shell git describe --always --dirty --tags 2> /dev/null)

.PHONY: all
all: clean release run

.PHONY: build
build:
	go build -ldflags "-X main.Version=$(VERSION)" -o $(BINARY) cmd/lilbib/main.go

.PHONY: test
test:
	go test -ldflags "-X main.Version=$(VERSION)" ./...

sandbox/config: config | sandbox/
	mkdir $@
	cp $^/config_test.toml $@/config.toml

.PHONY: sandbox/web
sandbox/web: | sandbox/
	rm -rf $@
	cp -r web $@

sandbox/$(BINARY): build | sandbox/
	cp $(BINARY) $@

.PHONY: sandbox
sandbox: sandbox/web sandbox/config sandbox/$(BINARY)

/tmp/lilbib-database.lock: database/lilbib_example.sql
	echo 'Starting database... (if it does not work add yourself to the group "docker")'
	./database/example-database.sh

.PHONY: database
database: /tmp/lilbib-database.lock

.PHONY: run
run: build sandbox database
	cd sandbox; ./$(BINARY)

.PHONY: release
release: linux windows

.PHONY: clean
clean:
	rm -rf sandbox $(BINARY) release

.PHONY: linux
linux:
	mkdir -p release/$(BINARY)-$(VERSION)-$@-amd64
	mkdir -p release/$(BINARY)-$(VERSION)-$@-amd64/config
	cp config/config.toml release/$(BINARY)-$(VERSION)-$@-amd64/config
	cp -r web release/$(BINARY)-$(VERSION)-$@-amd64
	cp database/lilbib.sql release/$(BINARY)-$(VERSION)-$@-amd64
	GOOS=$@ GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION)" -o release/$(BINARY)-$(VERSION)-$@-amd64/ ./...
	cd release; tar -czf $(BINARY)-$(VERSION)-$@-amd64.tar.gz $(BINARY)-$(VERSION)-$@-amd64

.PHONY: windows
windows:
	mkdir -p release/$(BINARY)-$(VERSION)-$@-amd64
	mkdir -p release/$(BINARY)-$(VERSION)-$@-amd64/config
	cp config/config.toml release/$(BINARY)-$(VERSION)-$@-amd64/config
	cp -r web release/$(BINARY)-$(VERSION)-$@-amd64
	cp database/lilbib.sql release/$(BINARY)-$(VERSION)-$@-amd64
	GOOS=$@ GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION)" -o release/$(BINARY)-$(VERSION)-$@-amd64/ ./...
	cd release; zip -qr $(BINARY)-$(VERSION)-$@-amd64.zip $(BINARY)-$(VERSION)-$@-amd64

%/:
	mkdir -p $*
