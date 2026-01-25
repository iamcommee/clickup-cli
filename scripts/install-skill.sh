#!/bin/bash
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

echo "Installing Claude Code skill..."
echo ""

# Get script directory (where the repo is)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_DIR="$(dirname "$SCRIPT_DIR")"

SKILL_DIR="$HOME/.claude/skills/clickup"
SKILL_SOURCE="$REPO_DIR/skills/SKILL.md"

if [ ! -f "$SKILL_SOURCE" ]; then
    echo -e "${RED}Error: Skill file not found at $SKILL_SOURCE${NC}"
    exit 1
fi

mkdir -p "$SKILL_DIR"
rm -f "$SKILL_DIR/SKILL.md"
ln -s "$SKILL_SOURCE" "$SKILL_DIR/SKILL.md"

echo -e "${GREEN}Claude Code skill installed!${NC}"
echo "  Location: $SKILL_DIR/SKILL.md"
echo "  Source: $SKILL_SOURCE"
echo ""
echo "You can now use '/clickup' in Claude Code"
