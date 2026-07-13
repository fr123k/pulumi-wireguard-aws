package verify

// MiniPCChecks returns the verification checks for a mini PC server.
// These check that the physical server was provisioned correctly with
// franky, nginx, WireGuard client, and security hardening.
func MiniPCChecks() []Check {
	return []Check{
		// SSH and user
		{
			Name:     "SSH service running",
			Command:  "systemctl is-active sshd || systemctl is-active ssh",
			Expected: "active",
		},
		{
			Name:    "franky user exists",
			Command: "id frank.ittermann 2>/dev/null && echo exists || echo missing",
			Equals:  "exists",
		},
		{
			Name:    "Password authentication disabled",
			Command: "grep -q 'PasswordAuthentication no' /etc/ssh/sshd_config && echo disabled || echo enabled",
			Equals:  "disabled",
		},

		// Firewall
		{
			Name:     "nftables service enabled",
			Command:  "systemctl is-enabled nftables",
			Expected: "enabled",
		},
		{
			Name:     "nftables service running",
			Command:  "systemctl is-active nftables",
			Expected: "active",
		},

		// fail2ban
		{
			Name:     "fail2ban installed",
			Command:  "fail2ban-client --version 2>&1",
			Expected: "exit:0",
		},
		{
			Name:     "fail2ban service running",
			Command:  "systemctl is-active fail2ban",
			Expected: "active",
		},

		// WireGuard client tools
		{
			Name:     "WireGuard tools installed",
			Command:  "which wg 2>/dev/null && echo installed || echo missing",
			Expected: "installed",
		},

		// Franky
		{
			Name:     "franky binary installed",
			Command:  "head -c 4 /usr/bin/franky | od -An -tx1 | grep -q '7f 45 4c 46' && echo 'valid ELF'",
			Expected: "valid ELF",
		},
		{
			Name:     "franky service running",
			Command:  "systemctl is-active franky",
			Expected: "active",
		},

		// Nginx
		{
			Name:     "Nginx installed",
			Command:  "nginx -v 2>&1",
			Expected: "exit:0",
		},
		{
			Name:     "Nginx service running",
			Command:  "systemctl is-active nginx",
			Expected: "active",
		},

		// Development tools
		{
			Name:     "Go installed",
			Command:  "go version 2>&1",
			Expected: "exit:0",
		},
		{
			Name:     "Zig installed",
			Command:  "zig version 2>&1",
			Expected: "exit:0",
		},

		// Docker (optional)
		{
			Name:     "Docker installed",
			Command:  "docker --version 2>/dev/null && echo installed || echo missing",
			Expected: "installed",
		},

		// Security
		{
			Name:    "Root SSH login disabled",
			Command: "grep -E '^PermitRootLogin' /etc/ssh/sshd_config",
			Equals:  "PermitRootLogin no",
		},
		{
			Name:     "Unattended upgrades configured",
			Command:  "test -f /etc/apt/apt.conf.d/20auto-upgrades && echo configured || echo missing",
			Expected: "configured",
		},
	}
}