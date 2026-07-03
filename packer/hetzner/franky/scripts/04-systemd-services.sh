#!/bin/bash
set -euxo pipefail

# 04-systemd-services.sh
# Creates systemd service files for sbx-web-ui and dandbox

# echo "=== Creating sbx user ==="

# # Create sbx user with home directory for runtime config
# if ! id sbx &>/dev/null; then
#     useradd -m -s /bin/bash -G docker sbx
#     echo "Created sbx user"
# else
#     echo "sbx user already exists"
# fi

# echo "=== Creating dandbox directories ==="

# # Create dandbox runtime directories
# SBX_HOME=$(eval echo ~sbx)
# mkdir -p "${SBX_HOME}/.config/dandbox/policy"
# mkdir -p "${SBX_HOME}/.config/dandbox/ca"
# mkdir -p "${SBX_HOME}/.local/state/dandbox"
# chown -R sbx:sbx "${SBX_HOME}/.config" "${SBX_HOME}/.local"

echo "=== Creating systemd service files ==="

# sbx-web-ui service
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

# # dandbox daemon service
# cat > /etc/systemd/system/dandbox.service <<'EOF'
# [Unit]
# Description=dandbox — sandboxed code execution daemon
# After=network.target docker.service
# Requires=docker.service
# StartLimitIntervalSec=120s
# StartLimitBurst=10

# [Service]
# User=sbx
# Group=sbx
# Environment="PROXY_SIDECAR_BIN=/usr/bin/dandbox"
# ExecStart=/usr/bin/dandbox \
#     -socket /home/sbx/.config/dandbox/dandbox.sock \
#     -docker-socket /var/run/docker.sock \
#     -state-dir /home/sbx/.local/state/dandbox \
#     -policy-dir /home/sbx/.config/dandbox/policy \
#     -ca-dir /home/sbx/.config/dandbox/ca
# Restart=on-failure
# RestartSec=10s

# [Install]
# WantedBy=multi-user.target
# EOF

# echo "=== Systemd service files created ==="
