#!/bin/bash

set -euo pipefail

NEW_MODULE=${1:-}
OLD_MODULE="github.com/username/project-name"

if [ -z "$NEW_MODULE" ]; then
  echo "Error: Please provide your new module name."
  exit 1
fi

go mod edit -module "$NEW_MODULE"

find . -type f \( -name "*.go" -o -name "*.md" -o -name "*.yml" -o -name "*.yaml" -o -name ".env.example" \) \
  -exec sed -i "" "s|$OLD_MODULE|$NEW_MODULE|g" {} +

echo "Module successfully renamed to $NEW_MODULE"
