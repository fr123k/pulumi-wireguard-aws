#!/bin/bash
set -euxo pipefail

# 04-systemd-services.sh
# Creates systemd service files for franky

echo "=== Creating systemd service files ==="

# franky service
cat > /etc/systemd/system/franky.service <<'EOF'
[Unit]
Description=franky
After=network.target
StartLimitIntervalSec=120s
StartLimitBurst=10

[Service]
Environment=""
WorkingDirectory=/home/frank.ittermann/github/
ExecStart=/usr/bin/franky --profile ollama-deepseek-flash --role full --yes --mode proxy
Restart=on-failure
RestartSec=10s

[Install]
WantedBy=multi-user.target
EOF

echo "=== Systemd service files created ==="