#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")"

export CGO_ENABLED=1

go test ./...
go build -o kaizen ./cmd/kaizen

echo "Built ./kaizen"