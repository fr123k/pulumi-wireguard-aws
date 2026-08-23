#!/bin/bash
set -euxo pipefail

# 01-base-packages.sh
# Installs base OS packages, WireGuard, and creates the wireguard runtime user

export DEBIAN_FRONTEND=noninteractive

echo "=== Installing base packages ==="

# Update and upgrade system
apt-get update -y
apt-get upgrade -y

# Install base packages
apt-get install -y \
    curl \
    gnupg2 \
    ca-certificates \
    lsb-release \
    ubuntu-keyring \
    jq \
    fail2ban \
    pwgen \
    ufw \
    wireguard \
    wireguard-tools

echo "=== Base packages installation complete ==="

# ============================================
# AppArmor profile fix for wg-quick on Ubuntu 26.04+
# The default profile denies access to /etc/nsswitch.conf, /etc/passwd,
# and coreutils locales.
# ============================================
echo "=== Patching wg-quick AppArmor profile ==="
mkdir -p /etc/apparmor.d/local
cat > /etc/apparmor.d/local/wg-quick <<'EOF'
  # Allow DNS resolution
  /etc/nsswitch.conf r,
  /etc/passwd r,
  # Allow coreutils locale files (Ubuntu 26.04+ gnu-coreutils)
  /usr/share/coreutils/locales/** r,
EOF
apparmor_parser -r /etc/apparmor.d/wg-quick || true

# ============================================
# Enable IP forwarding persistently
# (applied live at runtime by the cloud-init script)
# ============================================
echo "=== Enabling IP forwarding ==="
grep -q '^net.ipv4.ip_forward' /etc/sysctl.d/99-wireguard.conf 2>/dev/null || \
    echo 'net.ipv4.ip_forward=1' > /etc/sysctl.d/99-wireguard.conf

# Configure UFW default forward policy to ACCEPT (required for WireGuard NAT)
sed -i 's/DEFAULT_FORWARD_POLICY="DROP"/DEFAULT_FORWARD_POLICY="ACCEPT"/' /etc/default/ufw
grep -q 'DEFAULT_FORWARD_POLICY' /etc/default/ufw || \
    echo 'DEFAULT_FORWARD_POLICY="ACCEPT"' >> /etc/default/ufw

# Pre-configure UFW rules (enabled live at runtime to avoid locking out packer)
ufw allow http
ufw allow https
ufw allow 8443/tcp  # bypass wireguard firewall rule
ufw allow ssh
ufw allow 51820/udp

echo "=== Base setup complete ==="