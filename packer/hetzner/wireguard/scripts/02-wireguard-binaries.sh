#!/bin/bash
set -euxo pipefail

# 02-wireguard-binaries.sh
# Installs WireGuard UI, Secret Operator client, and Secret Operator server

echo "=== Installing WireGuard binaries ==="

# Versions are expected from versions.env
: "${WIREGUARD_UI_VERSION:?WIREGUARD_UI_VERSION is required}"
: "${SECRET_OPERATOR_VERSION:?SECRET_OPERATOR_VERSION is required}"

cd /tmp

# Install WireGuard UI (from fr123k/wireguard-ui fork)
echo "=== Installing WireGuard UI v${WIREGUARD_UI_VERSION} ==="
curl -OL "https://github.com/fr123k/wireguard-ui/releases/download/v${WIREGUARD_UI_VERSION}/wireguard-ui-v${WIREGUARD_UI_VERSION}-linux-amd64.tar.gz"
tar -xzf "wireguard-ui-v${WIREGUARD_UI_VERSION}-linux-amd64.tar.gz"
mv wireguard-ui /usr/local/bin/wireguard-ui
chmod +x /usr/local/bin/wireguard-ui
rm -f "wireguard-ui-v${WIREGUARD_UI_VERSION}-linux-amd64.tar.gz"

# Install Secret Operator Client
echo "=== Installing Secret Operator Client v${SECRET_OPERATOR_VERSION} ==="
curl -OL "https://github.com/containifyci/secret-operator/releases/download/v${SECRET_OPERATOR_VERSION}/secret-operator-client_linux_amd64"
mv secret-operator-client_linux_amd64 /usr/bin/secret-operator-client
chmod +x /usr/bin/secret-operator-client

# Install Secret Operator Server
echo "=== Installing Secret Operator Server v${SECRET_OPERATOR_VERSION} ==="
curl -OL "https://github.com/containifyci/secret-operator/releases/download/v${SECRET_OPERATOR_VERSION}/secret-operator-server_linux_amd64"
mv secret-operator-server_linux_amd64 /usr/bin/secret-operator-server
chmod +x /usr/bin/secret-operator-server

echo "=== WireGuard binaries installation complete ==="