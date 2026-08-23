#!/bin/bash
set -euxo pipefail

# 02-wireguard-binaries.sh
# Installs WireGuard UI, Secret Operator client, and Secret Operator server

echo "=== Installing WireGuard binaries ==="

# Versions are expected from versions.env
: "${WIREGUARD_UI_VERSION:?WIREGUARD_UI_VERSION is required}"
: "${SECRET_OPERATOR_VERSION:?SECRET_OPERATOR_VERSION is required}"
: "${TARGET_GOARCH:?TARGET_GOARCH is required}"

cd /tmp

# Install WireGuard UI (from ngoduykhanh/wireguard-ui upstream)
echo "=== Installing WireGuard UI v${WIREGUARD_UI_VERSION} (${TARGET_GOARCH}) ==="
curl -OL "https://github.com/ngoduykhanh/wireguard-ui/releases/download/v${WIREGUARD_UI_VERSION}/wireguard-ui-v${WIREGUARD_UI_VERSION}-linux-${TARGET_GOARCH}.tar.gz"
tar -xzf "wireguard-ui-v${WIREGUARD_UI_VERSION}-linux-${TARGET_GOARCH}.tar.gz"
mv wireguard-ui /usr/local/bin/wireguard-ui
chmod +x /usr/local/bin/wireguard-ui
rm -f "wireguard-ui-v${WIREGUARD_UI_VERSION}-linux-${TARGET_GOARCH}.tar.gz"

# Install Secret Operator Client
echo "=== Installing Secret Operator Client v${SECRET_OPERATOR_VERSION} (${TARGET_GOARCH}) ==="
curl -OL "https://github.com/containifyci/secret-operator/releases/download/v${SECRET_OPERATOR_VERSION}/secret-operator-client_linux_${TARGET_GOARCH}"
mv "secret-operator-client_linux_${TARGET_GOARCH}" /usr/bin/secret-operator-client
chmod +x /usr/bin/secret-operator-client

# Install Secret Operator Server
echo "=== Installing Secret Operator Server v${SECRET_OPERATOR_VERSION} (${TARGET_GOARCH}) ==="
curl -OL "https://github.com/containifyci/secret-operator/releases/download/v${SECRET_OPERATOR_VERSION}/secret-operator-server_linux_${TARGET_GOARCH}"
mv "secret-operator-server_linux_${TARGET_GOARCH}" /usr/bin/secret-operator-server
chmod +x /usr/bin/secret-operator-server

echo "=== WireGuard binaries installation complete ==="