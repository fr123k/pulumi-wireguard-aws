#!/bin/bash
set -euxo pipefail

# 04-systemd-services.sh
# Creates systemd service files for WireGuard UI, the wg-quick restart path,
# and the Secret Operator server. Static config only — runtime values
# (keys, passwords, env files) are populated by the cloud-init script.

echo "=== Creating WireGuard configuration directory ==="
mkdir -p /etc/wireguard
chown -R root:root /etc/wireguard

echo "=== Creating WireGuard UI data directory ==="
mkdir -p /usr/local/bin/db/server

echo "=== Creating systemd service files ==="

# wgui-restart-wg.service — restarts wg-quick when wg0.conf changes
cat > /etc/systemd/system/wgui-restart-wg.service <<'EOF'
[Unit]
Description=Restart WireGuard
After=network.target

[Service]
Type=oneshot
ExecStart=/usr/bin/systemctl restart wg-quick@wg0.service

[Install]
RequiredBy=wgui-restart-wg.path
EOF

# wgui-restart-wg.path — watches /etc/wireguard/wg0.conf for changes
cat > /etc/systemd/system/wgui-restart-wg.path <<'EOF'
[Unit]
Description=Watch /etc/wireguard/wg0.conf for changes

[Path]
PathModified=/etc/wireguard/wg0.conf

[Install]
WantedBy=multi-user.target
EOF

# wgui.service — runs the WireGuard UI web app
# NOTE: WGUI_PASSWORD and SESSION_SECRET are injected at runtime via
# /etc/systemd/system/wgui.env (created by the cloud-init script).
cat > /etc/systemd/system/wgui.service <<'EOF'
[Unit]
Description=Start wireguard-ui
After=network.target

[Service]
TimeoutStartSec=0
Restart=always
WorkingDirectory=/usr/local/bin
EnvironmentFile=/etc/systemd/system/wgui.env
ExecStart=/usr/local/bin/wireguard-ui --bind-address 127.0.0.1:8080

[Install]
WantedBy=multi-user.target
EOF

# Placeholder wgui.env (populated at runtime with WGUI_PASSWORD and SESSION_SECRET)
cat > /etc/systemd/system/wgui.env <<'EOF'
# Populated at runtime by cloud-init
WGUI_PASSWORD=
SESSION_SECRET=
EOF
chmod 600 /etc/systemd/system/wgui.env

# secret-operator-server.service — serves secrets over localhost:8088
# NOTE: GCP_PROJECT_ID and GOOGLE_APPLICATION_CREDENTIALS are runtime values
# and are injected via /etc/systemd/system/secret-operator-server.env
cat > /etc/systemd/system/secret-operator-server.service <<'EOF'
[Unit]
Description=Secret Operator Server
After=network.target
StartLimitIntervalSec=120s
StartLimitBurst=10

[Service]
EnvironmentFile=/etc/systemd/system/secret-operator-server.env
Environment="PORT=8088"
ExecStart=secret-operator-server
ExecStartPre=-secret-operator-server update
Restart=on-failure
RestartSec=10s

[Install]
WantedBy=multi-user.target
EOF

# Placeholder secret-operator-server.env (populated at runtime)
cat > /etc/systemd/system/secret-operator-server.env <<'EOF'
# Populated at runtime by cloud-init
GCP_PROJECT_ID=
GOOGLE_APPLICATION_CREDENTIALS=
EOF
chmod 600 /etc/systemd/system/secret-operator-server.env

# Reload systemd so the new units are picked up
systemctl daemon-reload

echo "=== Systemd services setup complete ==="