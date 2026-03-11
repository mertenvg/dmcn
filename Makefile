.PHONY: proto test test-cover lint clean

proto:
	buf generate

test:
	go test ./internal/...

test-cover:
	go test ./internal/core/crypto/... -coverprofile=coverage-crypto.out
	go tool cover -func=coverage-crypto.out | tail -1
	go test ./internal/core/identity/... -coverprofile=coverage-identity.out
	go tool cover -func=coverage-identity.out | tail -1
	go test ./internal/core/message/... -coverprofile=coverage-message.out
	go tool cover -func=coverage-message.out | tail -1

lint:
	buf lint
	go vet ./internal/...

clean:
	rm -f coverage-*.out
