# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands

```bash
make build          # Build binary to bin/clickup
make install        # Install to $GOPATH/bin
make test           # Run all tests
make fmt            # Format code
make vet            # Run go vet
make check          # Run fmt, vet, and test
make run ARGS="..." # Run without building (e.g., make run ARGS="tasks")
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
  - Loads from: CLI flags → env vars → `.clickup.json` (cwd) → `~/.clickup.json`
- `internal/output/` - Output formatting (table, JSON, detail)
  - Uses `Formatter` interface pattern
- `pkg/models/` - Data structures for API responses

### Key Patterns

- Global state in `root.go`: `cfg *config.Config` and `apiClient *api.Client` are lazily initialized via `getConfig()` and `getAPIClient()`
- Version info injected via ldflags at build time (see Makefile)
- Config file uses 0600 permissions for security
