#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
TLS_DIR="$ROOT_DIR/deploy/overlays/local/tls"
CERT="$TLS_DIR/tls.crt"
KEY="$TLS_DIR/tls.key"

FORCE=false
if [[ ${1:-} == "--force" ]]; then
  FORCE=true
fi

if [[ -f "$CERT" || -f "$KEY" ]]; then
  if [[ "$FORCE" == "false" ]]; then
    echo "TLS files already exist. Use --force to regenerate." >&2
    exit 0
  fi
fi

mkdir -p "$TLS_DIR"

if command -v mkcert >/dev/null 2>&1; then
  mkcert -install >/dev/null 2>&1 || true
  mkcert -cert-file "$CERT" -key-file "$KEY" \
    localhost 127.0.0.1 resource-checker.local
else
  if ! command -v openssl >/dev/null 2>&1; then
    echo "openssl not found in PATH" >&2
    exit 1
  fi
  openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
    -keyout "$KEY" \
    -out "$CERT" \
    -subj "/CN=resource-checker.local" \
    -addext "subjectAltName=DNS:localhost,IP:127.0.0.1"
fi

echo "Generated $CERT and $KEY"
