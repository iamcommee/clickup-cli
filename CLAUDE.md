# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands

```bash
make build          # Build binary to bin/clickup
make build-all      # Build for all platforms (updates pre-built binaries)
make install        # Install pre-built binary to /usr/local/bin
make install-skill  # Install Claude Code skill only
make install-all    # Install binary + Claude Code skill
make uninstall      # Remove binary, skill, and config
make test           # Run all tests
make fmt            # Format code
make vet            # Run go vet
make check          # Run fmt, vet, and test
make run ARGS="..." # Run without building (e.g., make run ARGS="tasks")
```

## Installation

Pre-built binaries are included in `bin/` for all platforms. Users don't need Go installed.

```bash
git clone https://github.com/iamcommee/clickup-cli.git
cd clickup-cli
make install        # Installs pre-built binary
clickup install     # Configure API token
```

## Architecture

This is a Go CLI application for ClickUp task management built with Cobra.

### Package Structure

- `cmd/clickup/main.go` - Entry point, calls `commands.Execute()`
- `internal/commands/` - Cobra command definitions and CLI logic
  - `root.go` - Root command with global flags, shared config/client initialization
  - Commands register themselves in `init()` via `rootCmd.AddCommand()`
- `internal/api/` - ClickUp API client
  - `client.go` - HTTP client with auth header injection
  - Domain-specific files (`tasks.go`, `user.go`) add methods to Client
- `internal/config/` - Configuration management
  - Loads from: CLI flags → env vars → `~/.config/clickup/config.json`
- `internal/output/` - Output formatting (table, JSON, detail)
  - Uses `Formatter` interface pattern
- `pkg/models/` - Data structures for API responses

### Key Patterns

- Global state in `root.go`: `cfg *config.Config` and `apiClient *api.Client` are lazily initialized via `getConfig()` and `getAPIClient()`
- Version info injected via ldflags at build time (see Makefile)
- Config file uses 0600 permissions for security

### Security

- API tokens are encrypted using AES-GCM before storing
- Encryption key is derived from machine-specific data (hostname + username)
- Encrypted tokens have `enc:` prefix for identification
- Config file: `~/.config/clickup/config.json`
- Encryption implementation: `internal/config/crypto.go`
