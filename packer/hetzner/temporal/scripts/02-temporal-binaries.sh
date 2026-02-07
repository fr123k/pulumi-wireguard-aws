#!/bin/bash
set -euxo pipefail

# 02-temporal-binaries.sh
# Installs Temporal CLI, Server, UI Server, Secret Operator Client, and Worker

echo "=== Installing Temporal binaries ==="

# Versions are expected from versions.env
: "${TEMPORAL_CLI_VERSION:?TEMPORAL_CLI_VERSION is required}"
: "${TEMPORAL_SERVER_VERSION:?TEMPORAL_SERVER_VERSION is required}"
: "${TEMPORAL_UI_VERSION:?TEMPORAL_UI_VERSION is required}"
: "${TEMPORAL_WORKER_VERSION:?TEMPORAL_WORKER_VERSION is required}"
: "${SECRET_OPERATOR_VERSION:?SECRET_OPERATOR_VERSION is required}"

cd /tmp

# Install Temporal CLI
echo "=== Installing Temporal CLI v${TEMPORAL_CLI_VERSION} ==="
curl -OL "https://github.com/temporalio/cli/releases/download/v${TEMPORAL_CLI_VERSION}/temporal_cli_${TEMPORAL_CLI_VERSION}_linux_amd64.tar.gz"
tar -xzf "temporal_cli_${TEMPORAL_CLI_VERSION}_linux_amd64.tar.gz"
mv temporal /usr/bin/temporal
chmod +x /usr/bin/temporal
rm -f "temporal_cli_${TEMPORAL_CLI_VERSION}_linux_amd64.tar.gz"

# Install Temporal Server
echo "=== Installing Temporal Server v${TEMPORAL_SERVER_VERSION} ==="
curl -OL "https://github.com/temporalio/temporal/releases/download/v${TEMPORAL_SERVER_VERSION}/temporal_${TEMPORAL_SERVER_VERSION}_linux_amd64.tar.gz"
tar -xzf "temporal_${TEMPORAL_SERVER_VERSION}_linux_amd64.tar.gz"
mv temporal-server /usr/bin/temporal-server
chmod +x /usr/bin/temporal-server
rm -f "temporal_${TEMPORAL_SERVER_VERSION}_linux_amd64.tar.gz"

# Install Temporal UI Server
echo "=== Installing Temporal UI Server v${TEMPORAL_UI_VERSION} ==="
curl -OL "https://github.com/temporalio/ui-server/releases/download/v${TEMPORAL_UI_VERSION}/ui-server_${TEMPORAL_UI_VERSION}_linux_amd64.tar.gz"
tar -xzf "ui-server_${TEMPORAL_UI_VERSION}_linux_amd64.tar.gz"
mv ui-server /usr/bin/temporal-ui-server
chmod +x /usr/bin/temporal-ui-server
rm -f "ui-server_${TEMPORAL_UI_VERSION}_linux_amd64.tar.gz"

# Install Secret Operator Client
echo "=== Installing Secret Operator Client v${SECRET_OPERATOR_VERSION} ==="
curl -OL "https://github.com/containifyci/secret-operator/releases/download/v${SECRET_OPERATOR_VERSION}/secret-operator-client_linux_amd64"
mv secret-operator-client_linux_amd64 /usr/bin/secret-operator-client
chmod +x /usr/bin/secret-operator-client

# Install Temporal Worker
echo "=== Installing Temporal Worker v${TEMPORAL_WORKER_VERSION} ==="
curl -OL "https://github.com/containifyci/temporal-worker/releases/download/v${TEMPORAL_WORKER_VERSION}//temporal-dunebot-worker_linux_amd64"
mv temporal-dunebot-worker_linux_amd64 /usr/bin/temporal-dunebot-worker
chmod +x /usr/bin/temporal-dunebot-worker

echo "=== Temporal binaries installation complete ==="
