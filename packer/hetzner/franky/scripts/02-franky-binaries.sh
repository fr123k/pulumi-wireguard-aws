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

echo "=== install tools (zig, golang) ==="

cd /tmp
wget https://ziglang.org/builds/zig-x86_64-linux-0.17.0-dev.1267+300116b02.tar.xz
tar -xf zig-x86_64-linux-0.17.0-dev.1267+300116b02.tar.xz -C /usr/local/bin --strip-components=1
rm zig-x86_64-linux-0.17.0-dev.1267+300116b02.tar.xz

wget https://go.dev/dl/go1.26.4.linux-amd64.tar.gz
rm -rf /usr/local/go
tar -C /usr/local -xzf go1.26.4.linux-amd64.tar.gz
ln -sf /usr/local/go/bin/go /usr/local/bin/go
ln -sf /usr/local/go/bin/gofmt /usr/local/bin/gofmt

echo "=== tools installation complete ==="