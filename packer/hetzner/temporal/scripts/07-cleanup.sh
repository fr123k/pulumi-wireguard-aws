#!/bin/bash
set -euxo pipefail

# 07-cleanup.sh
# Cleans up apt cache, temporary files, and logs for a smaller image

echo "=== Cleaning up for snapshot ==="

# Clean apt cache
apt-get clean
apt-get autoremove -y
rm -rf /var/lib/apt/lists/*

# Remove temporary files
rm -rf /tmp/*
rm -rf /var/tmp/*

# Clean up cloud-init artifacts (will run fresh on boot)
cloud-init clean --logs

# Truncate log files
truncate -s 0 /var/log/*.log 2>/dev/null || true
truncate -s 0 /var/log/**/*.log 2>/dev/null || true

# Remove SSH host keys (will be regenerated on first boot)
rm -f /etc/ssh/ssh_host_*

# Remove machine-id (will be regenerated)
truncate -s 0 /etc/machine-id

# Clear bash history
rm -f /root/.bash_history
history -c || true

echo "=== Cleanup complete - image ready for snapshot ==="
