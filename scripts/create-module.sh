#!/bin/bash

set -euo pipefail

MODULE=${1:-}

if [ -z "$MODULE" ]; then
  echo "usage: ./scripts/create-module.sh <module-name>"
  exit 1
fi

mkdir -p "internal/modules/$MODULE"

touch "internal/modules/$MODULE/module.go"
touch "internal/modules/$MODULE/routes.go"
touch "internal/modules/$MODULE/handler.go"
touch "internal/modules/$MODULE/service.go"
touch "internal/modules/$MODULE/repository.go"
touch "internal/modules/$MODULE/dto.go"

echo "module $MODULE created"
