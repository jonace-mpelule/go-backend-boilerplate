#!/bin/bash

set -euo pipefail

if ! rg -q '^run:' Makefile; then
  echo "Makefile is missing the run target"
  exit 1
fi

if ! rg -q 'bin/server' README.md Makefile Dockerfile; then
  echo "server binary references are inconsistent"
  exit 1
fi

if rg -n 'cmd/api|bin/api|\./bin/api' README.md Makefile Dockerfile .github >/dev/null; then
  echo "stale api binary references found"
  exit 1
fi

if [ ! -f .env.example ]; then
  echo ".env.example is missing"
  exit 1
fi

echo "template consistency checks passed"
