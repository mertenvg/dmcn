.PHONY: proto test test-cover lint clean build build-bridge build-web build-release

VERSION  ?= $(shell git describe --tags --always 2>/dev/null || echo dev)
LDFLAGS   = -ldflags "-X main.version=$(VERSION)"
DIST      = dist/$(VERSION)
BINARIES  = dmcn-node dmcn-bridge dmcn-web
PLATFORMS = linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 windows/arm64

build:
	go build $(LDFLAGS) -o bin/dmcn-node ./cmd/dmcn-node
	go build $(LDFLAGS) -o bin/dmcn-bridge ./cmd/dmcn-bridge
	go build $(LDFLAGS) -o bin/dmcn-web ./cmd/dmcn-web

build-bridge:
	go build $(LDFLAGS) -o bin/dmcn-bridge ./cmd/dmcn-bridge

build-web:
	cd cmd/dmcn-web/web && npm ci && npm run build && cd ../../..
	go build $(LDFLAGS) -o bin/dmcn-web ./cmd/dmcn-web

build-release:
	@rm -rf $(DIST)
	@for platform in $(PLATFORMS); do \
		os=$${platform%/*}; \
		arch=$${platform#*/}; \
		ext=""; \
		if [ "$$os" = "windows" ]; then ext=".exe"; fi; \
		dir="$(DIST)/dmcn-$$os-$$arch"; \
		mkdir -p "$$dir"; \
		for bin in $(BINARIES); do \
			echo "building $$bin $$os/$$arch"; \
			GOOS=$$os GOARCH=$$arch go build $(LDFLAGS) -o "$$dir/$$bin$$ext" ./cmd/$$bin || exit 1; \
		done; \
		if [ "$$os" = "windows" ] && command -v zip >/dev/null 2>&1; then \
			(cd $(DIST) && zip -qr "dmcn-$(VERSION)-$$os-$$arch.zip" "dmcn-$$os-$$arch"); \
		else \
			tar -czf "$(DIST)/dmcn-$(VERSION)-$$os-$$arch.tar.gz" -C $(DIST) "dmcn-$$os-$$arch"; \
		fi; \
		rm -rf "$$dir"; \
	done
	@cd $(DIST) && sha256sum *.tar.gz *.zip > checksums-sha256.txt 2>/dev/null; true
	@echo "release artifacts in $(DIST):"
	@ls -1 $(DIST)

clean:
	rm -f *.out bin/dmcn-*
	rm -rf dist cmd/dmcn-web/web/dist cmd/dmcn-web/web/node_modules

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
