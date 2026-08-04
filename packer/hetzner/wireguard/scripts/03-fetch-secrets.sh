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
secret-operator-client fetch -token="${SECRET_OPERATOR_TOKEN}" --output=/tmp/secrets

# Install wildcard SSL certificate
mkdir -p /etc/letsencrypt/live/${DOMAIN}
mv /tmp/secrets/certs/fullchain.pem /etc/letsencrypt/live/${DOMAIN}/fullchain.pem
mv /tmp/secrets/certs/privkey.pem /etc/letsencrypt/live/${DOMAIN}/privkey.pem
chmod 600 /etc/letsencrypt/live/${DOMAIN}/privkey.pem
chmod 644 /etc/letsencrypt/live/${DOMAIN}/fullchain.pem

echo "=== SSL certificate installed for ${DOMAIN} ==="