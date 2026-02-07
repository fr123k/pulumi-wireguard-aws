#!/bin/bash
set -euxo pipefail

# 03-dunebot-binaries.sh
# Installs OAuth2 Storage and DuneBot binaries

echo "=== Installing DuneBot binaries ==="

# Versions are expected from versions.env
: "${OAUTH2_STORAGE_VERSION:?OAUTH2_STORAGE_VERSION is required}"
: "${DUNEBOT_VERSION:?DUNEBOT_VERSION is required}"

cd /tmp

# Install OAuth2 Storage
echo "=== Installing OAuth2 Storage v${OAUTH2_STORAGE_VERSION} ==="
curl -OL "https://github.com/containifyci/oauth2-storage/releases/download/v${OAUTH2_STORAGE_VERSION}/oauth2-storage_linux_amd64"
mv oauth2-storage_linux_amd64 /usr/bin/oauth2-storage
chmod +x /usr/bin/oauth2-storage

# Install DuneBot
echo "=== Installing DuneBot v${DUNEBOT_VERSION} ==="
curl -OL "https://github.com/containifyci/dunebot/releases/download/v${DUNEBOT_VERSION}/dunebot_linux_amd64"
mv dunebot_linux_amd64 /usr/bin/dunebot
chmod +x /usr/bin/dunebot

echo "=== DuneBot binaries installation complete ==="
