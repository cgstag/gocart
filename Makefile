MAIN_VERSION:=$(shell git describe --abbrev=0 --tags || echo "0.1.0")
VERSION:=${MAIN_VERSION}
PACKAGES:=$(shell /usr/local/go/bin/go list ./... | sed -n '1!p' | grep -v './vendor')
PACKAGE_WITHOUT_E2E:=$(shell /usr/local/go/bin/go list ./... | sed -n '1!p' | grep -v -E '(./vendor|./testdata|./daos|./routes)')
LDFLAGS:="-X main.Version=${VERSION}"

# Debug mode... prints with verbose.
# To use: DEBUG=true make test
ifeq ($(DEBUG),true)
	VERBOSE=-v
endif

default: test

deps:
	go get -v .

deps-all: deps
	go get -v github.com/stretchr/testify/assert
	go get -v github.com/mitchellh/gox
	go get -v github.com/githubnemo/CompileDaemon
	go get -v github.com/golang/dep/cmd/dep
	go get -v github.com/Kount/pq-timeouts
	go get -v github.com/patrickmn/go-cache

test:
	go test $(VERBOSE) -p=1 $(PACKAGES)

cover:
	echo "mode: count" > coverage-all.out
	$(foreach pkg,$(PACKAGES), \
		go test $(VERBOSE) -p=1 -cover -covermode=count -coverprofile=coverage.out ${pkg}; \
		tail -n +2 coverage.out >> coverage-all.out;)
	go tool cover -html=coverage-all.out

unit-test:
	go test $(VERBOSE) $(PACKAGE_WITHOUT_E2E)

unit-cover:
	echo "mode: count" > coverage-all.out
	$(foreach pkg,$(PACKAGE_WITHOUT_E2E), \
		go test $(VERBOSE) -p=1 -cover -covermode=count -coverprofile=coverage.out ${pkg}; \
		tail -n +2 coverage.out >> coverage-all.out;)
	go tool cover -html=coverage-all.out

daos-test:
	go test $(VERBOSE) -p=1 ./daos

daos-cover:
	echo "mode: count" > coverage-all.out
	go test $(VERBOSE) -p=1 -cover -covermode=count -coverprofile=coverage.out ./daos
	tail -n +2 coverage.out >> coverage-all.out
	go tool cover -html=coverage-all.out

run:
	go run -ldflags=${LDFLAGS} main.go

dev:
	CompileDaemon -exclude-dir ".git" -exclude-dir "vendor" -color -build "go build -o _build_hot_reload" -command "./_build_hot_reload"

build: clean
	CGO_ENABLED=0 go build $(VERBOSE) -ldflags=${LDFLAGS} -a -o api_${VERSION} main.go

build-all: clean
	gox -ldflags=${LDFLAGS} -output="{{.Dir}}_{{.OS}}_{{.Arch}}_${VERSION}" -osarch="windows/amd64" -osarch="linux/amd64" -osarch="darwin/amd64"

clean:
	rm -rf _build_hot_reload coverage.out coverage-all.out
