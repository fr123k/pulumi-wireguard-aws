#!/bin/bash
set -euxo pipefail

# 04-nginx-setup.sh
# Installs Nginx from official repo and Certbot (no SSL certs - those are runtime)

export DEBIAN_FRONTEND=noninteractive

echo "=== Installing Nginx ==="

# Add Nginx official repository
curl -fsSL https://nginx.org/keys/nginx_signing.key | gpg --dearmor -o /usr/share/keyrings/nginx-archive-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/nginx-archive-keyring.gpg] http://nginx.org/packages/ubuntu $(lsb_release -cs) nginx" | tee /etc/apt/sources.list.d/nginx.list
apt-get update
apt-get install -y nginx

# Create sites-available directory
mkdir -p /etc/nginx/sites-available

# Remove default config
rm -f /etc/nginx/conf.d/default.conf

echo "=== Installing Certbot ==="

# Install Certbot via snap
snap install --classic certbot
ln -sf /snap/bin/certbot /usr/bin/certbot

# Create directory for certificates (fetched at runtime from Secret Manager)
mkdir -p /etc/letsencrypt/live

echo "=== Nginx setup complete ==="
