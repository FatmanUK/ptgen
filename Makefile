# Variables
GO = go
APP_BINARY = ptgen_relaxed
V_TAG = $(shell git tag --list 'v*' | tail -n1)
VERSION = $(if $(V_TAG),$(V_TAG),v0.0.0)

BUILD_DATE = $(shell date +%Y%m%d)
BUILD_TIMESTAMP = $(shell date +%Y%m%dT%H%M00Z)
CKSUM_SCRIPT = import hashlib; print(hashlib.sha1(open('./$(APP_BINARY)','rb').read()).hexdigest())

BUILD_ENV=CGO_ENABLED=0
BUILD_LDFLAGS=-s -w -X main.APP_NAME=$(APP_BINARY) -X main.VERSION=$(VERSION)

# Default target
all: clean podman-build

.PHONY: clean test lint mod-download build

clean:
	rm ./$(APP_BINARY)

# Run the app locally using go run
run:
	$(GO) run ./main

# Build static, stripped binaries for minimal final images
mod-download:
	$(GO) mod download

build: mod-download
	$(BUILD_ENV) $(GO) build \
		-ldflags="$(BUILD_LDFLAGS)" \
		-o ./$(APP_BINARY) \
		./main

# Run unit tests
test:
	$(GO) test -race -v ./...

# Lint code
lint:
	$(GO) fmt ./...
