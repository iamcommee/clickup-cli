# ClickUp CLI

A command-line interface for ClickUp. List your assigned tasks and view task details directly from the terminal.

## Installation

```bash
git clone <repo-url>
cd clickup-cli
make build

# Or install to $GOPATH/bin
make install
```

## Quick Start

1. Get your API token from https://app.clickup.com/settings/apps

2. Install and configure:
   ```bash
   clickup install
   ```

3. List your assigned tasks:
   ```bash
   clickup tasks
   ```

## Commands

### Tasks

```bash
# List all assigned tasks
clickup tasks

# Filter by status
clickup tasks --status "in progress"

# Include closed tasks
clickup tasks --closed

# Limit results
clickup tasks --limit 10

# Output as JSON
clickup tasks --output json

# Get task details
clickup tasks TASK-123

# Open task in browser
clickup tasks TASK-123 --open
```

### Configuration

```bash
# Install (interactive setup)
clickup install

# Install with token (non-interactive)
clickup install --token pk_xxxxx

# Show current config
clickup config

# Remove all configuration
clickup uninstall
```

### Other

```bash
clickup version    # Show version
clickup --help     # Help
```

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--config` | `-c` | Custom config file path |
| `--workspace` | `-w` | Override workspace ID |
| `--debug` | | Enable debug output |

## Configuration

Config file location: `~/.config/clickup/config.json`

The CLI looks for configuration in this order:
1. Command-line flags
2. Environment variables: `CLICKUP_API_TOKEN`, `CLICKUP_WORKSPACE_ID`, `CLICKUP_USER_ID`
3. `~/.config/clickup/config.json`

### Security

- API tokens are encrypted using AES-GCM before storing
- Config file permissions are set to 0600 (owner only)
- Config directory permissions are set to 0700 (owner only)

## Development

```bash
make build      # Build binary
make install    # Install to GOPATH/bin
make test       # Run tests
make fmt        # Format code
make vet        # Run go vet
make check      # Run fmt, vet, and test
```

## License

MIT
