# Build variables
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

BINARY_NAME = clickup
BUILD_DIR = bin

LDFLAGS = -ldflags "\
	-X clickup-cli/internal/commands.Version=$(VERSION) \
	-X clickup-cli/internal/commands.Commit=$(COMMIT) \
	-X clickup-cli/internal/commands.BuildDate=$(BUILD_DATE)"

.PHONY: all build install install-skill uninstall clean test fmt vet check build-all run help

all: build

# Build the binary
build:
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/clickup

# Install binary only
install: build
	@./scripts/install.sh

# Install Claude Code skill only
install-skill:
	@./scripts/install-skill.sh

# Install binary + Claude Code skill
install-all: build
	@./scripts/install.sh --with-skill

# Uninstall everything (binary, skill, config)
uninstall:
	@./scripts/uninstall.sh

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)

# Run tests
test:
	go test -v ./...

# Format code
fmt:
	go fmt ./...

# Run go vet
vet:
	go vet ./...

# Run all checks
check: fmt vet test

# Build for multiple platforms
build-all: build-linux build-darwin build-windows

build-linux:
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/clickup
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/clickup

build-darwin:
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/clickup
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/clickup

build-windows:
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/clickup

# Development run
run:
	go run ./cmd/clickup $(ARGS)

# Show help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build         Build binary to bin/clickup"
	@echo "  install       Install binary to /usr/local/bin"
	@echo "  install-skill Install Claude Code skill only"
	@echo "  install-all   Install binary + Claude Code skill"
	@echo "  uninstall     Remove binary, skill, and config"
	@echo "  clean         Remove build artifacts"
	@echo "  test          Run tests"
	@echo "  fmt           Format code"
	@echo "  vet           Run go vet"
	@echo "  check         Run fmt, vet, and test"
	@echo "  build-all     Build for all platforms"
	@echo "  run           Run with ARGS (e.g., make run ARGS='tasks')"
	@echo "  help          Show this help"
