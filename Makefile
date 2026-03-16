.PHONY: proto test test-cover lint clean build build-bridge

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

build:
	go build -o bin/dmcn-node ./cmd/dmcn-node
	go build -o bin/dmcn-bridge ./cmd/dmcn-bridge

build-bridge:
	go build -o bin/dmcn-bridge ./cmd/dmcn-bridge

clean:
	rm -f coverage-*.out bin/dmcn-node bin/dmcn-bridge
