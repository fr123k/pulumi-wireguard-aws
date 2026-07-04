#!/bin/bash
set -euxo pipefail

echo "=== Installing franky ==="

# Version is expected from versions.env
: "${FRANKY_VERSION:?FRANKY_VERSION is required}"

cd /tmp

# Install franky
echo "=== Installing franky v${FRANKY_VERSION} ==="
curl -OL "https://github.com/fr12k/franky/releases/download/v${FRANKY_VERSION}/franky_linux_amd64"
mv franky_linux_amd64 /usr/bin/franky
chmod +x /usr/bin/franky

# Verify it's a valid ELF binary
echo "=== Verifying franky binary ==="
head -c 4 /usr/bin/franky | od -An -tx1 | grep -q '7f 45 4c 46' && echo "valid ELF binary" || echo "WARNING: not an ELF binary"

echo "=== franky installation complete ==="