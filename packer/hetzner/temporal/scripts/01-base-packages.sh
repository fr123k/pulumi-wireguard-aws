#!/bin/bash
set -euxo pipefail

# 01-base-packages.sh
# Installs base OS packages and creates temporal user

export DEBIAN_FRONTEND=noninteractive

echo "=== Installing base packages ==="

# Update and upgrade system
apt-get update -y
apt-get upgrade -y

# Install base packages
apt-get install -y \
    curl \
    gnupg2 \
    ca-certificates \
    lsb-release \
    ubuntu-keyring \
    jq \
    fail2ban

# Create temporal user for running services
echo "=== Creating temporal user ==="
useradd -r -s /bin/false temporal || true

# Create temporal configuration directory
mkdir -p /etc/temporal
chown temporal:temporal /etc/temporal

echo "=== Base packages installation complete ==="
