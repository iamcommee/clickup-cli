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

2. Initialize configuration:
   ```bash
   clickup config init
   ```

3. List your assigned tasks:
   ```bash
   clickup tasks
   ```

## Commands

### Tasks

```bash
# List all assigned tasks (sorted by ID)
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
clickup tasks HGAI-1217

# Open task in browser
clickup tasks HGAI-1217 --open

# Get task as JSON
clickup tasks HGAI-1217 --output json
```

### Configuration

```bash
# Initialize (interactive)
clickup config init

# Initialize with token
clickup config init --token pk_xxxxx

# Show current config
clickup config show
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

The CLI looks for configuration in this order:

1. Command-line flags
2. Environment variables: `CLICKUP_API_TOKEN`, `CLICKUP_WORKSPACE_ID`, `CLICKUP_USER_ID`
3. `.clickup.json` in current directory
4. `~/.clickup.json` in home directory

### Config File

```json
{
  "api_token": "pk_xxxxx",
  "workspace_id": "12345678",
  "user_id": "87654321"
}
```

## Development

```bash
make build      # Build binary
make install    # Install to GOPATH/bin
make fmt        # Format code
make vet        # Run go vet
make build-all  # Build for all platforms
```

## License

MIT
