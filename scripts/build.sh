#!/usr/bin/env bash

set -e

for lambda in create resolve processor; do
  mkdir -p dist/$lambda
  GOOS=linux GOARCH=arm64 go build -o dist/$lambda/bootstrap cmd/lambda/$lambda/main.go
  (cd dist/$lambda && zip -q function.zip bootstrap)
  echo "Built $lambda"
done
