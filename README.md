# ClickUp CLI

A command-line interface for ClickUp. List your assigned tasks and view task details directly from the terminal.

## Installation

```bash
git clone https://github.com/iamcommee/clickup-cli.git
cd clickup-cli
make install
```

### With Claude Code Integration

```bash
make install-all    # Install binary + Claude Code skill
# or
make install        # Install binary only
make install-skill  # Add Claude Code skill later
```

The skill is symlinked, so `git pull` will automatically update it.

### Uninstall

```bash
make uninstall      # Remove binary, skill, and config
```

## Quick Start

1. Get your API token from https://app.clickup.com/settings/apps

2. Configure the CLI:
   ```bash
   clickup install
   ```

3. List your assigned tasks:
   ```bash
   clickup tasks
   ```

## Claude Code Integration

After installation, you can use `/clickup` in Claude Code:

```
/clickup tasks           # List your tasks
/clickup                 # Claude will help with ClickUp tasks
```

Or just ask Claude naturally: "Show my ClickUp tasks"

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
# Configure CLI (interactive)
clickup install

# Configure with token (non-interactive)
clickup install --token pk_xxxxx

# Show current config
clickup config

# Remove configuration
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

Config file: `~/.config/clickup/config.json`

The CLI looks for configuration in this order:
1. Command-line flags
2. Environment variables: `CLICKUP_API_TOKEN`, `CLICKUP_WORKSPACE_ID`, `CLICKUP_USER_ID`
3. `~/.config/clickup/config.json`

### Security

- API tokens are encrypted using AES-GCM before storing
- Config file permissions: 0600 (owner only)
- Config directory permissions: 0700 (owner only)

## Development

```bash
make build         # Build binary to bin/clickup
make install       # Install binary to /usr/local/bin
make install-skill # Install Claude Code skill
make install-all   # Install binary + skill
make uninstall     # Remove binary, skill, and config
make test          # Run tests
make check         # Run fmt, vet, and test
make help          # Show all targets
```

## License

MIT
