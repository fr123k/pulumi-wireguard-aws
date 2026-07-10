#!/bin/bash
set -euxo pipefail

# 03-nginx-setup.sh
# Installs Nginx from official repo (no SSL certs — those are runtime)

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

# Create franky nginx config
cat > /etc/nginx/sites-available/franky <<'EOF'
server {
    listen 80;
    listen [::]:80;
    server_name _;

    access_log /var/log/nginx/franky.access.log;
    error_log /var/log/nginx/franky.error.log;

    location / {
        proxy_pass http://127.0.0.1:8787;
        proxy_http_version 1.1;
        proxy_read_timeout 300;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $http_host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Real-PORT $remote_port;
        error_page 403 = /custom_444;

        location = /custom_444 {
            return 444;
        }
    }
}
EOF

# Enable franky site
ln -sf /etc/nginx/sites-available/franky /etc/nginx/conf.d/franky.conf

echo "=== Nginx setup complete ==="