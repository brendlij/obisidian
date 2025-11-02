#!/usr/bin/env bash
set -euo pipefail
cd mcs-manager && CGO_ENABLED=0 go build -o ../bin/mcs-manager ./cmd/mcs-manager
