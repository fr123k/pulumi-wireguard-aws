#!/bin/bash
set -euxo pipefail

# 01-base-packages.sh
# Installs base OS packages for mini PC server

export DEBIAN_FRONTEND=noninteractive

echo "=== Installing base packages ==="

# Update and upgrade system
apt-get update -y
apt-get upgrade -y

# Install base packages
apt-get install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    fail2ban \
    gnupg2 \
    jq \
    lsb-release \
    nftables \
    pwgen \
    software-properties-common \
    ubuntu-keyring \
    wireguard-tools

echo "=== Base packages installation complete ==="