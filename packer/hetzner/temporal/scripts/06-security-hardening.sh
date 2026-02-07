#!/bin/bash
set -euxo pipefail

# 06-security-hardening.sh
# Configures SSH hardening and fail2ban

echo "=== Configuring SSH hardening ==="

# Harden SSH configuration
# Note: PermitRootLogin and AllowUsers will be finalized at runtime
# Here we just prepare the base hardening
sed -i 's|[#]*PasswordAuthentication yes|PasswordAuthentication no|g' /etc/ssh/sshd_config

# Add algorithm support if not already present
if ! grep -q "PubkeyAcceptedAlgorithms" /etc/ssh/sshd_config; then
    echo "PubkeyAcceptedAlgorithms +ssh-rsa" >> /etc/ssh/sshd_config
fi

echo "=== Configuring fail2ban ==="

# Create nginx-4xx filter
cat > /etc/fail2ban/filter.d/nginx-4xx.conf <<'EOF'
[Definition]
failregex = ^<HOST> - \S+ \[\] "[^"]*" (404|444|403|400)

datepattern = {^LN-BEG}%%ExY(?P<_sep>[-/.])%%m(?P=_sep)%%d[T ]%%H:%%M:%%S(?:[.,]%%f)?(?:\s*%%z)?
              ^[^\[]*\[({DATE})
              {^LN-BEG}
EOF

# Configure fail2ban jails
cat > /etc/fail2ban/jail.d/defaults-debian.conf <<'EOF'
[DEFAULT]
banaction = nftables
banaction_allports = nftables[type=allports]
backend = systemd

[sshd]
enabled = true

[nginx-4xx]
enabled  = true
port     = http,https
filter   = nginx-4xx
logpath  = %(nginx_access_log)s
backend  = polling
maxretry = 3
EOF

# Enable fail2ban but don't start (will start at runtime)
systemctl enable fail2ban

echo "=== Security hardening complete ==="
