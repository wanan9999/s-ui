#!/bin/sh
set -eu
ROOT=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
cd "$ROOT"
[ "$(uname -s)" = Linux ] || { echo "仅支持 Linux 构建" >&2; exit 1; }
(cd frontend && npm ci && npm run build)
mkdir -p web/html
cp -R frontend/dist/. web/html/
. "$ROOT/build-tags.sh"
CGO_ENABLED=0 go build -trimpath -buildvcs=false -ldflags "$(ldflags_for dev)" -tags "$(tags_for dev)" -o sui .
