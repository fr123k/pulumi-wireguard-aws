#!/bin/bash
set -euxo pipefail

# 05-systemd-services.sh
# Creates all systemd service files and static configuration files

echo "=== Creating Temporal configuration files ==="

# Create temporal configuration directory
mkdir -p /etc/temporal

# Temporal Server configuration
cat > /etc/temporal/temporal-server.yaml <<'EOF'
log:
  stdout: true
  level: info

persistence:
  defaultStore: sqlite-default
  visibilityStore: sqlite-visibility
  numHistoryShards: 4
  datastores:
    sqlite-default:
      sql:
        pluginName: "sqlite"
        databaseName: "/etc/temporal/default.db"
        connectAddr: "localhost"
        connectProtocol: "tcp"
        connectAttributes:
          cache: "private"
          setup: true

    sqlite-visibility:
      sql:
        pluginName: "sqlite"
        databaseName: "/etc/temporal/visibility.db"
        connectAddr: "localhost"
        connectProtocol: "tcp"
        connectAttributes:
          cache: "private"
          setup: true

global:
  membership:
    maxJoinDuration: 30s
    broadcastAddress: "127.0.0.1"
  pprof:
    port: 7936

services:
  frontend:
    rpc:
      grpcPort: 7236
      membershipPort: 6933
      bindOnLocalHost: true
      httpPort: 7243

  matching:
    rpc:
      grpcPort: 7235
      membershipPort: 6935
      bindOnLocalHost: true

  history:
    rpc:
      grpcPort: 7234
      membershipPort: 6934
      bindOnLocalHost: true

  worker:
    rpc:
      membershipPort: 6939

clusterMetadata:
  enableGlobalNamespace: false
  failoverVersionIncrement: 10
  masterClusterName: "active"
  currentClusterName: "active"
  clusterInformation:
    active:
      enabled: true
      initialFailoverVersion: 1
      rpcName: "frontend"
      rpcAddress: "localhost:7236"
      httpAddress: "localhost:7243"

dcRedirectionPolicy:
  policy: "noop"
EOF

# Temporal UI Server configuration
cat > /etc/temporal/temporal-ui-server.yaml <<'EOF'
temporalGrpcAddress: 127.0.0.1:7236
host: 127.0.0.1
port: 8233
enableUi: true
cors:
  allowOrigins:
    - http://localhost:8233
defaultNamespace: default
EOF

chown -R temporal:temporal /etc/temporal

echo "=== Creating systemd service files ==="

# Temporal Server service
cat > /etc/systemd/system/temporal.service <<'EOF'
[Unit]
Description=Temporal Service
After=network.target

[Service]
User=temporal
Group=temporal
ExecStart=temporal-server -r / -c etc/temporal/ -e temporal-server start

[Install]
WantedBy=multi-user.target
EOF

# Temporal UI service
cat > /etc/systemd/system/temporal-ui.service <<'EOF'
[Unit]
Description=Temporal UI Server
After=network.target

[Service]
User=temporal
Group=temporal
ExecStart=temporal-ui-server -r / -c etc/temporal/ -e temporal-ui-server start

[Install]
WantedBy=multi-user.target
EOF

# Temporal Worker Update service
cat > /etc/systemd/system/temporal-dunebot-worker-update.service <<'EOF'
[Unit]
Description=Temporal Worker Update

[Service]
Type=oneshot
ExecStart=temporal-dunebot-worker update

[Install]
WantedBy=multi-user.target
EOF

# Temporal Worker Update timer
cat > /etc/systemd/system/temporal-dunebot-worker-update.timer <<'EOF'
[Unit]
Description=Run self update for temporal-dunebot-worker
Requires=temporal-dunebot-worker.service

[Timer]
Unit=temporal-dunebot-worker-update.service
OnCalendar=*:00/30

[Install]
WantedBy=timers.target
EOF

# Temporal Worker service
cat > /etc/systemd/system/temporal-dunebot-worker.service <<'EOF'
[Unit]
Description=Temporal Worker
After=network.target
StartLimitIntervalSec=120s
StartLimitBurst=10

[Service]
User=temporal
Group=temporal
Environment="ADDRESS=127.0.0.1"
Environment="ADDRESS=127.0.0.1"
Environment="DUNEBOT_JWT_SERVER_ADDRESS=localhost:50051"
Environment="DUNEBOT_GITHUB_OAUTH_SCOPES=repo"
Environment="DUNEBOT_GITHUB_V3_API_URL=https://api.github.com/"
Environment="DUNEBOT_GITHUB_WEB_URL=https://github.com/"
Environment="DUNEBOT_APP_CONFIGURATION_REVIEWER_TYPE=direct"
EnvironmentFile=/etc/systemd/system/temporal-worker.env
ExecStart=-temporal-dunebot-worker
ExecStartPre=-temporal-dunebot-worker update
Restart=on-failure
RestartSec=10s

[Install]
WantedBy=multi-user.target
EOF

# OAuth2 Storage Update service
cat > /etc/systemd/system/oauth2-storage-update.service <<'EOF'
[Unit]
Description=OAuth2 Storage Update

[Service]
Type=oneshot
ExecStart=oauth2-storage update

[Install]
WantedBy=multi-user.target
EOF

# OAuth2 Storage Update timer
cat > /etc/systemd/system/oauth2-storage-update.timer <<'EOF'
[Unit]
Description=Run self update for oauth2-storage
Requires=oauth2-storage-update.service

[Timer]
Unit=oauth2-storage-update.service
OnCalendar=*:00/30

[Install]
WantedBy=timers.target
EOF

# OAuth2 Storage service
cat > /etc/systemd/system/oauth2-storage.service <<'EOF'
[Unit]
Description=OAuth2 Storage
After=network.target
StartLimitIntervalSec=120s
StartLimitBurst=10

[Service]
EnvironmentFile=/etc/systemd/system/temporal-worker.env
ExecStart=oauth2-storage
ExecStartPre=-oauth2-storage update
Restart=on-failure
RestartSec=10s

[Install]
WantedBy=multi-user.target
EOF

# DuneBot Update service
cat > /etc/systemd/system/dunebot-update.service <<'EOF'
[Unit]
Description=DuneBot Update Service

[Service]
Type=oneshot
ExecStart=dunebot update

[Install]
WantedBy=multi-user.target
EOF

# DuneBot Update timer
cat > /etc/systemd/system/dunebot-update.timer <<'EOF'
[Unit]
Description=Run self update for DuneBot
Requires=dunebot-update.service

[Timer]
Unit=dunebot-update.service
OnCalendar=*:00/30

[Install]
WantedBy=timers.target
EOF

# DuneBot service
cat > /etc/systemd/system/dunebot.service <<'EOF'
[Unit]
Description=DuneBot Service
After=network.target
StartLimitIntervalSec=120s
StartLimitBurst=10

[Service]
Environment="PORT=8080"
Environment="ADDRESS=127.0.0.1"
Environment="DUNEBOT_JWT_SERVER_ADDRESS=localhost:50051"
Environment="DUNEBOT_GITHUB_OAUTH_SCOPES=repo"
Environment="DUNEBOT_GITHUB_V3_API_URL=https://api.github.com/"
Environment="DUNEBOT_GITHUB_WEB_URL=https://github.com/"
Environment="DUNEBOT_APP_CONFIGURATION_REVIEWER_TYPE=temporal"
Environment="DUNEBOT_APP_CONFIGURATION_REVIEWER_SERVER_ADDRESS=127.0.0.1:7233"
EnvironmentFile=/etc/systemd/system/temporal-worker.env
ExecStart=dunebot app
Restart=on-failure
RestartSec=10s

[Install]
WantedBy=multi-user.target
EOF

# DuneBot Dispatch Job service
cat > /etc/systemd/system/dunebot-dispatch-job.service <<'EOF'
[Unit]
Description=DuneBot Dispatch Job
After=network.target

[Service]
EnvironmentFile=/etc/systemd/system/temporal-worker.env
ExecStart=dunebot dispatch

[Install]
WantedBy=multi-user.target
EOF

# DuneBot Dispatch Job timer
cat > /etc/systemd/system/dunebot-dispatch-job.timer <<'EOF'
[Unit]
Description=Run dunebot dispatch job

[Timer]
Unit=dunebot-dispatch-job.service
OnCalendar=*-*-* 10:00:00 UTC
Persistent=true
AccuracySec=1m

[Install]
WantedBy=timers.target
EOF

# Create placeholder for temporal-worker.env (will be populated at runtime)
touch /etc/systemd/system/temporal-worker.env
chmod 600 /etc/systemd/system/temporal-worker.env

# Reload systemd
systemctl daemon-reload

echo "=== Systemd services setup complete ==="
