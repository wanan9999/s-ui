#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"
docker buildx build --platform linux/amd64,linux/arm64 --progress=plain . 2>&1 | tee docker-build-test.log
