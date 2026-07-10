#!/bin/bash
set -euxo pipefail

# 05-security-hardening.sh
# Configures nftables firewall, SSH hardening, and fail2ban

echo "=== Setting up nftables firewall ==="

cat > /etc/nftables.conf <<'NFTABLES'
#!/usr/sbin/nft -f
# Mini PC firewall rules

flush ruleset

table inet filter {
    chain input {
        type filter hook input priority 0; policy drop;

        # Allow loopback
        iif "lo" accept

        # Allow established/related connections
        ct state established,related accept
        ct state invalid drop

        # Allow ICMP (ping)
        ip protocol icmp accept
        ip6 nexthdr icmpv6 accept

        # Allow SSH
        tcp dport { 22 } accept

        # Allow WireGuard client port (for outgoing connections)
        udp sport { 51820 } accept

        # Allow HTTP/HTTPS
        tcp dport { 80, 443 } accept

        # Allow DNS
        tcp dport { 53 } accept
        udp dport { 53 } accept

        # Allow mDNS for local network discovery
        udp dport { 5353 } accept

        # Allow franky web UI (proxied through nginx)
        tcp dport { 8787 } accept

        # Log dropped packets
        log prefix "nftables-drop: " drop
    }

    chain forward {
        type filter hook forward priority 0; policy drop;

        # Allow established/related
        ct state established,related accept

        # Log dropped packets
        log prefix "nftables-fwd-drop: " drop
    }

    chain output {
        type filter hook output priority 0; policy accept;
    }
}
NFTABLES

systemctl enable nftables

echo "=== Configuring SSH hardening ==="

# Harden SSH configuration
# Note: PermitRootLogin and AllowUsers will be finalized at runtime
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

# Enable automatic security updates
apt-get install -y unattended-upgrades
dpkg-reconfigure -plow unattended-upgrades

cat > /etc/apt/apt.conf.d/20auto-upgrades <<'EOF'
APT::Periodic::Update-Package-Lists "1";
APT::Periodic::Unattended-Upgrade "1";
APT::Periodic::AutocleanInterval "7";
EOF

echo "=== Security hardening complete ==="