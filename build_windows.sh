#!/bin/bash
# Builds the Windows GUI binary from the current tree with the version injected
# at link time. Kept as a script because interactive shells on this host load a
# mise hook that mangles inline compound commands.
set -euo pipefail
cd /Users/sepyankristanto/Documents/3.Data_Lainnya/whatsapp-web.view

VERSION="${1:-1.5.9.4}"
OUT="dist_win/WhatsAppDesk.exe"
mkdir -p dist_win

export CGO_ENABLED=0
export GOOS=windows
export GOARCH="${2:-amd64}"

go build -ldflags="-H windowsgui -s -w -buildid= -X main.appVersion=${VERSION}" \
  -trimpath -o "$OUT" .

echo "built: $OUT"
ls -la "$OUT"
file "$OUT"
