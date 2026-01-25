---
name: clickup
description: Use this skill when the user wants to list ClickUp tasks, view task details, or manage ClickUp CLI configuration. Triggers on phrases like "show my tasks", "clickup tasks", "what are my assigned tasks", "task details".
user-invocable: true
argument-hint: "[tasks|config|install]"
---

# ClickUp CLI Skill

You have access to the `clickup` CLI tool for managing ClickUp tasks.

**IMPORTANT**: You MUST display the result exactly as you receive it.

## Available Commands

### List Tasks
```bash
clickup tasks                          # List all assigned tasks
clickup tasks --status "in progress"   # Filter by status
clickup tasks --limit 10               # Limit results
clickup tasks --output json            # Output as JSON
```

### View Task Details
```bash
clickup tasks TASK-ID                  # View task details
clickup tasks TASK-ID --open           # Open in browser
clickup tasks TASK-ID --output json    # Output as JSON
```

### Configuration
```bash
clickup config                         # Show current config
clickup install                        # Set up CLI with API token
clickup uninstall                      # Remove configuration
```

## Usage Guidelines

1. When user asks about their tasks, run `clickup tasks` first
2. Use `--output json` when you need to process the data programmatically
3. If CLI is not configured, guide user to run `clickup install`
4. Task IDs look like `PROJ-123` or just `abc123`

## Error Handling

If you see "config file not found", tell the user to run:
```bash
clickup install
```

This will prompt for their ClickUp API token from: https://app.clickup.com/settings/apps
