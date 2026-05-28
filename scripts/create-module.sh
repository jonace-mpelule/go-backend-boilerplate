#!/bin/bash

MODULE=$1

mkdir -p internal/modules/$MODULE

touch internal/modules/$MODULE/module.go
touch internal/modules/$MODULE/routes.go
touch internal/modules/$MODULE/handler.go
touch internal/modules/$MODULE/service.go
touch internal/modules/$MODULE/repository.go

echo "module $MODULE created"
