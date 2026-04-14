BINARY    := goto
PKG       := github.com/aaangelmartin/goto
VERSION   := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT    := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE      := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS   := -s -w \
  -X $(PKG)/internal/buildinfo.Version=$(VERSION) \
  -X $(PKG)/internal/buildinfo.Commit=$(COMMIT) \
  -X $(PKG)/internal/buildinfo.Date=$(DATE)

BIN_DIR := bin

.PHONY: help
help: ## Show this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: build
build: ## Build the binary to ./bin/goto (ad-hoc signs on macOS)
	@mkdir -p $(BIN_DIR)
	go build -ldflags '$(LDFLAGS)' -o $(BIN_DIR)/$(BINARY) ./cmd/goto
	@if [ "$$(uname)" = "Darwin" ]; then \
		codesign --force --sign - $(BIN_DIR)/$(BINARY) 2>/dev/null || true; \
	fi

.PHONY: install
install: ## Install to $GOBIN
	go install -ldflags '$(LDFLAGS)' ./cmd/goto

.PHONY: run
run: ## Run the binary (pass args with ARGS="...")
	go run ./cmd/goto $(ARGS)

.PHONY: test
test: ## Run tests with race detector
	go test -race -cover ./...

.PHONY: test-verbose
test-verbose: ## Run tests verbose
	go test -race -v ./...

.PHONY: lint
lint: ## Run golangci-lint
	golangci-lint run ./...

.PHONY: fmt
fmt: ## Format code
	gofmt -s -w .
	go mod tidy

.PHONY: demo
demo: ## Generate demo GIF with vhs
	@command -v vhs >/dev/null || { echo "vhs not installed (brew install vhs)"; exit 1; }
	vhs demo/demo.tape

.PHONY: snapshot
snapshot: ## Local snapshot release with goreleaser
	goreleaser release --snapshot --clean

.PHONY: clean
clean: ## Clean build artifacts
	rm -rf $(BIN_DIR) dist

.PHONY: tidy
tidy: ## Tidy go.mod/go.sum
	go mod tidy
