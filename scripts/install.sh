#!/bin/bash
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Parse arguments
INSTALL_SKILL=false
for arg in "$@"; do
    case $arg in
        --with-skill)
            INSTALL_SKILL=true
            shift
            ;;
    esac
done

echo "ClickUp CLI Installer"
echo "====================="
echo ""

# Get script directory (where the repo is)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_DIR="$(dirname "$SCRIPT_DIR")"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *) echo -e "${RED}Unsupported architecture: $ARCH${NC}"; exit 1 ;;
esac

case $OS in
    darwin|linux) ;;
    *) echo -e "${RED}Unsupported OS: $OS${NC}"; exit 1 ;;
esac

BINARY_NAME="clickup-${OS}-${ARCH}"
BINARY_PATH="$REPO_DIR/bin/$BINARY_NAME"
INSTALL_DIR="/usr/local/bin"

echo "Detected: $OS/$ARCH"
echo ""

# Check pre-built binary exists
if [ ! -f "$BINARY_PATH" ]; then
    echo -e "${RED}Error: Pre-built binary not found at $BINARY_PATH${NC}"
    echo "Please ensure you have the complete repository."
    exit 1
fi

# Install binary
echo "Installing clickup to $INSTALL_DIR..."
if [ -w "$INSTALL_DIR" ]; then
    cp "$BINARY_PATH" "$INSTALL_DIR/clickup"
    chmod +x "$INSTALL_DIR/clickup"
else
    echo "Need sudo to install to $INSTALL_DIR"
    sudo cp "$BINARY_PATH" "$INSTALL_DIR/clickup"
    sudo chmod +x "$INSTALL_DIR/clickup"
fi
echo -e "${GREEN}Binary installed to $INSTALL_DIR/clickup${NC}"

# Install Claude Code skill (optional)
if [ "$INSTALL_SKILL" = true ]; then
    SKILL_DIR="$HOME/.claude/skills/clickup"
    SKILL_SOURCE="$REPO_DIR/skills/SKILL.md"

    if [ -f "$SKILL_SOURCE" ]; then
        echo ""
        echo "Installing Claude Code skill..."
        mkdir -p "$SKILL_DIR"
        rm -f "$SKILL_DIR/SKILL.md"
        ln -s "$SKILL_SOURCE" "$SKILL_DIR/SKILL.md"
        echo -e "${GREEN}Claude Code skill linked to $SKILL_DIR${NC}"
    else
        echo -e "${YELLOW}Skill file not found, skipping...${NC}"
    fi
fi

# Verify installation
echo ""
if command -v clickup &> /dev/null; then
    echo -e "${GREEN}Success!${NC}"
    echo ""
    clickup version
    echo ""
    echo "Next steps:"
    echo "  1. Run 'clickup install' to configure your API token"
    echo "  2. Run 'clickup tasks' to list your tasks"
    if [ "$INSTALL_SKILL" = true ]; then
        echo "  3. Use '/clickup' in Claude Code to access tasks"
    else
        echo ""
        echo "To add Claude Code integration, run:"
        echo "  make install-skill"
    fi
else
    echo -e "${RED}Installation may have failed. Please check your PATH.${NC}"
    exit 1
fi
