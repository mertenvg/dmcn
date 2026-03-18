.PHONY: proto test test-cover lint clean build build-bridge build-web build-linux

VERSION ?= $(shell git describe --tags --always 2>/dev/null || echo dev)
LDFLAGS  = -ldflags "-X main.version=$(VERSION)"

build:
	go build $(LDFLAGS) -o bin/dmcn-node ./cmd/dmcn-node
	go build $(LDFLAGS) -o bin/dmcn-bridge ./cmd/dmcn-bridge
	go build $(LDFLAGS) -o bin/dmcn-web ./cmd/dmcn-web

build-bridge:
	go build $(LDFLAGS) -o bin/dmcn-bridge ./cmd/dmcn-bridge

build-web:
	cd cmd/dmcn-web/web && npm ci && npm run build && cd ../../..
	go build $(LDFLAGS) -o bin/dmcn-web ./cmd/dmcn-web

GOOS   ?= linux
GOARCH ?= amd64
build-linux:
	mkdir -p bin/$(GOOS)-$(GOARCH)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o bin/$(GOOS)-$(GOARCH)/dmcn-node ./cmd/dmcn-node
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o bin/$(GOOS)-$(GOARCH)/dmcn-bridge ./cmd/dmcn-bridge
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o bin/$(GOOS)-$(GOARCH)/dmcn-web ./cmd/dmcn-web

clean:
	rm -f *.out bin/dmcn-* bin/linux-amd64/*
	rm -rf cmd/dmcn-web/web/dist cmd/dmcn-web/web/node_modules

proto:
	buf generate

test:
	go test ./internal/... -timeout 120s

test-cover:
	@echo "=== M1: Crypto ==="
	go test ./internal/core/crypto/... -coverprofile=coverage-crypto.out
	go tool cover -func=coverage-crypto.out | tail -1
	@echo "=== M1: Identity ==="
	go test ./internal/core/identity/... -coverprofile=coverage-identity.out
	go tool cover -func=coverage-identity.out | tail -1
	@echo "=== M1: Message ==="
	go test ./internal/core/message/... -coverprofile=coverage-message.out
	go tool cover -func=coverage-message.out | tail -1
	@echo "=== M2: Keystore ==="
	go test ./internal/keystore/... -coverprofile=coverage-keystore.out
	go tool cover -func=coverage-keystore.out | tail -1
	@echo "=== M2: Relay ==="
	go test ./internal/relay/... -coverprofile=coverage-relay.out
	go tool cover -func=coverage-relay.out | tail -1
	@echo "=== M2: Node + Integration ==="
	go test ./internal/node/... -timeout 120s -coverprofile=coverage-node.out
	go tool cover -func=coverage-node.out | tail -1
	@echo "=== M3: Bridge ==="
	go test ./internal/bridge/... -timeout 120s -coverprofile=coverage-bridge.out
	go tool cover -func=coverage-bridge.out | tail -1

lint:
	buf lint
	go vet ./internal/...
