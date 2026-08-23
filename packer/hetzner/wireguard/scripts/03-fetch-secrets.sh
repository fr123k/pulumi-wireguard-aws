#!/bin/bash
set -euxo pipefail

# 03-fetch-secrets.sh
# Fetches the wildcard SSL certificate from the Secret Operator at build time
# and installs it into /etc/letsencrypt/live/<domain>/ so it is baked into the
# snapshot. This avoids the runtime cloud-init needing to fetch certs on boot.
#
# If SECRET_OPERATOR_TOKEN is empty, this script is a no-op (certs will be
# fetched at runtime instead).
#
# Environment:
#   SECRET_OPERATOR_TOKEN  — Secret Operator authentication token (required)
#   DOMAIN                 — Domain for the SSL certificate (required)
#   SECRET_OPERATOR_HOST   — Secret Operator server URL (default: https://wg.fr123k.uk:8443/secrets)
#                           For VirtualBox NAT builds pointing at a local host service, use http://10.0.2.2:8888/secrets

echo "=== Fetching SSL certificate at build time ==="

if [ -z "${SECRET_OPERATOR_TOKEN:-}" ]; then
    echo "SECRET_OPERATOR_TOKEN not set — skipping build-time cert fetch (will fetch at runtime)"
    exit 0
fi

if [ -z "${DOMAIN:-}" ]; then
    echo "ERROR: DOMAIN is required when SECRET_OPERATOR_TOKEN is set"
    exit 1
fi

# Fetch secrets from the Secret Operator
secret-operator-client fetch \
  -token="${SECRET_OPERATOR_TOKEN}" \
  --host="${SECRET_OPERATOR_HOST:-https://wg.fr123k.uk:8443/secrets}" \
  --output=/tmp/secrets

# Install wildcard SSL certificate
mkdir -p /etc/letsencrypt/live/${DOMAIN}
mv /tmp/secrets/certs/fullchain.pem /etc/letsencrypt/live/${DOMAIN}/fullchain.pem
mv /tmp/secrets/certs/privkey.pem /etc/letsencrypt/live/${DOMAIN}/privkey.pem
chmod 600 /etc/letsencrypt/live/${DOMAIN}/privkey.pem
chmod 644 /etc/letsencrypt/live/${DOMAIN}/fullchain.pem

# Populate /etc/systemd/    system/service_account_adc.json from fetched secrets (if present)
if [ -f /tmp/secrets/service_account_adc-json ]; then
    mv /tmp/secrets/service_account_adc-json /etc/systemd/system/service_account_adc.json
    chmod 600 /etc/systemd/system/service_account_adc.json
fi

echo "=== SSL certificate installed for ${DOMAIN} ==="