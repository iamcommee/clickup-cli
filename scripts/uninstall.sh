#!/bin/bash
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "ClickUp CLI Uninstaller"
echo "======================="
echo ""

INSTALL_DIR="/usr/local/bin"
SKILL_DIR="$HOME/.claude/skills/clickup"
CONFIG_DIR="$HOME/.config/clickup"

# Confirm
read -p "This will remove clickup CLI, Claude skill, and config. Continue? [y/N] " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Aborted."
    exit 0
fi

echo ""

# Remove binary
if [ -f "$INSTALL_DIR/clickup" ]; then
    echo "Removing binary..."
    if [ -w "$INSTALL_DIR/clickup" ]; then
        rm "$INSTALL_DIR/clickup"
    else
        sudo rm "$INSTALL_DIR/clickup"
    fi
    echo -e "${GREEN}Removed $INSTALL_DIR/clickup${NC}"
else
    echo -e "${YELLOW}Binary not found at $INSTALL_DIR/clickup${NC}"
fi

# Remove Claude skill
if [ -d "$SKILL_DIR" ]; then
    echo "Removing Claude Code skill..."
    rm -rf "$SKILL_DIR"
    echo -e "${GREEN}Removed $SKILL_DIR${NC}"
else
    echo -e "${YELLOW}Skill not found at $SKILL_DIR${NC}"
fi

# Remove config
if [ -d "$CONFIG_DIR" ]; then
    echo "Removing configuration..."
    rm -rf "$CONFIG_DIR"
    echo -e "${GREEN}Removed $CONFIG_DIR${NC}"
else
    echo -e "${YELLOW}Config not found at $CONFIG_DIR${NC}"
fi

echo ""
echo -e "${GREEN}Uninstall complete.${NC}"
