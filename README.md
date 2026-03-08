# ClickUp CLI

A command-line interface for ClickUp. List your assigned tasks and view task details directly from the terminal.

## Installation

```bash
git clone https://github.com/iamcommee/clickup-cli.git
cd clickup-cli
make install
```

Then configure your API token:
```bash
clickup install
```
Get your token from: https://app.clickup.com/settings/apps

### With Claude Code Integration

```bash
make install-all    # Install binary + Claude Code skill
# or
make install-skill  # Add Claude Code skill later
```

The skill is symlinked, so `git pull` will automatically update it.

### Uninstall

```bash
make uninstall      # Remove binary, skill, and config
```

## Claude Code Integration

After installing the skill, you can use `/clickup` in Claude Code:

```
/clickup tasks           # List your tasks
/clickup                 # Claude will help with ClickUp tasks
```

Or just ask Claude naturally: "Show my ClickUp tasks"

## Commands

### Tasks

```bash
clickup tasks                          # List all assigned tasks
clickup tasks --status "in progress"   # Filter by status
clickup tasks --closed                 # Include closed tasks
clickup tasks --limit 10               # Limit results
clickup tasks --output json            # Output as JSON
clickup tasks TASK-123                 # Get task details
clickup tasks TASK-123 --open          # Open in browser
clickup tasks TASK-123 --comments      # Show task comments
```

### Configuration

```bash
clickup install                        # Configure CLI (interactive)
clickup install --token pk_xxxxx       # Configure with token
clickup config                         # Show current config
clickup uninstall                      # Remove configuration
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

Requires Go 1.21+.

```bash
make build         # Build binary for current platform
make build-all     # Build for all platforms (updates pre-built binaries)
make install       # Build from source and install to /usr/local/bin
make install-skill # Install Claude Code skill
make install-all   # Install binary + skill
make uninstall     # Remove binary, skill, and config
make test          # Run tests
make check         # Run fmt, vet, and test
make help          # Show all targets
```

## License

MIT
