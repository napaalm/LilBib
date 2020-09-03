BINARY := lilbib
VERSION ?= $(shell git describe --always --dirty --tags 2> /dev/null)

.PHONY: all
all: clean release run

.PHONY: build
build: | build/
	go build -ldflags "-X main.Version=$(VERSION)" -o build/$(BINARY) cmd/lilbib/main.go

.PHONY: test
test:
	go test -ldflags "-X main.Version=$(VERSION)" ./...

.PHONY: tidy
tidy:
	./scripts/tidy.sh web/template

build/web/template/%.html: web/template/%.html | build/web/template/
	tidy -qi -w 0 -omit --output-bom no --hide-comments yes --vertical-space auto --indent no -o $@ $^; [ $$? -eq 2 ] && { echo "Template HTML code is not valid!"; exit 1; } || { exit 0; }

build/web/template/%.xml: web/template/%.xml | build/web/template/
	tidy -qi -w 0 -xml -omit --output-bom no --hide-comments yes --vertical-space auto --indent no -o $@ $^; [ $$? -eq 2 ] && { echo "Template XML code is not valid!"; exit 1; } || { exit 0; }

build/web/static: web/static | build/web/static/
	cp -r $^ $@

.PHONY: web
web: $(addprefix build/,$(wildcard web/template/*.html) $(wildcard web/template/*.xml) web/static)

sandbox/config: config | sandbox/
	mkdir $@
	cp $^/config_test.toml $@/config.toml

.PHONY: sandbox/web
sandbox/web: web | sandbox/
	rm -rf $@
	ln -s ../build/web $@

sandbox/$(BINARY): build | sandbox/
	cp build/$(BINARY) $@

.PHONY: sandbox
sandbox: sandbox/web sandbox/config sandbox/$(BINARY)

/tmp/lilbib-database.lock: database/lilbib_example.sql
	echo 'Starting database... (if it does not work add yourself to the group "docker")'
	./database/example-database.sh

.PHONY: database
database: /tmp/lilbib-database.lock

.PHONY: run
run: database sandbox
	cd sandbox; ./$(BINARY)

.PHONY: release
release: linux windows

.PHONY: clean
clean:
	rm -rf build sandbox release


.PHONY: linux
linux: web
	mkdir -p release/$(BINARY)-$(VERSION)-$@-amd64
	mkdir -p release/$(BINARY)-$(VERSION)-$@-amd64/config
	cp config/config.toml release/$(BINARY)-$(VERSION)-$@-amd64/config
	cp -r build/web release/$(BINARY)-$(VERSION)-$@-amd64
	cp database/lilbib.sql release/$(BINARY)-$(VERSION)-$@-amd64
	GOOS=$@ GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION)" -o release/$(BINARY)-$(VERSION)-$@-amd64/ ./...
	cd release; tar -czf $(BINARY)-$(VERSION)-$@-amd64.tar.gz $(BINARY)-$(VERSION)-$@-amd64

.PHONY: windows
windows: web
	mkdir -p release/$(BINARY)-$(VERSION)-$@-amd64
	mkdir -p release/$(BINARY)-$(VERSION)-$@-amd64/config
	cp config/config.toml release/$(BINARY)-$(VERSION)-$@-amd64/config
	cp -r build/web release/$(BINARY)-$(VERSION)-$@-amd64
	cp database/lilbib.sql release/$(BINARY)-$(VERSION)-$@-amd64
	GOOS=$@ GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION)" -o release/$(BINARY)-$(VERSION)-$@-amd64/ ./...
	cd release; zip -qr $(BINARY)-$(VERSION)-$@-amd64.zip $(BINARY)-$(VERSION)-$@-amd64

%/:
	mkdir -p $*
