#!/bin/bash
set -euxo pipefail

# 01-base-packages.sh
# Installs base OS packages, docker-sbx CLI, and creates sbx user

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
    pwgen \
    software-properties-common \
    ubuntu-keyring

# ============================================
# Install Docker apt repo and docker
# ============================================
# echo "=== Installing docker-sbx CLI ==="

# # Add Docker's official GPG key and repository (repo only mode)
# curl -fsSL https://get.docker.com | sudo REPO_ONLY=1 sh

# # Install docker-sbx CLI (Docker Engine)
# apt-get update -y
# apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin 

echo "=== Base packages installation complete ==="
