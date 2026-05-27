#!/bin/bash
# Usage: ./module_name.sh ://github.com

NEW_MODULE=$1
OLD_MODULE="github.com/username/project-name"

if [ -z "$NEW_MODULE" ]; then
    echo "Error: Please provide your new module name."
    exit 1
fi

# Update go.mod module declaration
go mod edit -module "$NEW_MODULE"

# Replace all internal imports in .go files (Works on Linux/macOS)
find . -type f -name "*.go" -exec sed -i "" "s|$OLD_MODULE|$NEW_MODULE|g" {} +

echo "Module successfully renamed to $NEW_MODULE!"
