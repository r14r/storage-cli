set shell := ["bash", "-euo", "pipefail", "-c"]
default:
    @just --list
bootstrap:
    go mod download
    just fmt
    chmod +x .githooks/pre-commit
    git config --local core.hooksPath .githooks
fmt:
    gofmt -w cmd internal
fmt-check:
    @out="$(gofmt -l cmd internal)"; test -z "$out" || { printf '%s\n' "$out"; exit 1; }
vet:
    go vet ./...
test:
    go test -count=1 ./...
test-race:
    go test -race -count=1 ./...
precommit: fmt-check vet test
pre-commit: precommit
check: precommit
ci: check
build: fmt check
    mkdir -p dist
    go build -trimpath -ldflags "-s -w" -o dist/storage-cli ./cmd/storage-cli
install: build
    @destination="${INSTALL_DIR:-/usr/local/bin}"; mkdir -p "$destination"; install -m 0755 dist/storage-cli "$destination/storage-cli"
deploy: install
run *args:
    go run ./cmd/storage-cli {{args}}
release: check
    mkdir -p dist
    GOOS=darwin GOARCH=amd64 go build -trimpath -o dist/storage-cli-darwin-amd64 ./cmd/storage-cli
    GOOS=darwin GOARCH=arm64 go build -trimpath -o dist/storage-cli-darwin-arm64 ./cmd/storage-cli
    GOOS=linux GOARCH=amd64 go build -trimpath -o dist/storage-cli-linux-amd64 ./cmd/storage-cli
clean:
    rm -rf dist
