package verify

// PrebakedChecks returns the checks for verifying a pre-baked image
// ELF magic bytes check: first 4 bytes should be 7f 45 4c 46 (\x7fELF)
func PrebakedChecks() []Check {
	return []Check{
		{
			Name:     "Temporal CLI is valid ELF binary",
			Command:  "head -c 4 /usr/bin/temporal | od -An -tx1 | grep -q '7f 45 4c 46' && echo 'valid ELF'",
			Expected: "valid ELF",
		},
		{
			Name:     "Temporal server is valid ELF binary",
			Command:  "head -c 4 /usr/bin/temporal-server | od -An -tx1 | grep -q '7f 45 4c 46' && echo 'valid ELF'",
			Expected: "valid ELF",
		},
		{
			Name:     "Temporal UI server is valid ELF binary",
			Command:  "head -c 4 /usr/bin/temporal-ui-server | od -An -tx1 | grep -q '7f 45 4c 46' && echo 'valid ELF'",
			Expected: "valid ELF",
		},
		{
			Name:     "Temporal DuneBot Worker is valid ELF binary",
			Command:  "head -c 4 /usr/bin/temporal-dunebot-worker | od -An -tx1 | grep -q '7f 45 4c 46' && echo 'valid ELF'",
			Expected: "valid ELF",
		},
		{
			Name:     "Nginx installed",
			Command:  "nginx -v 2>&1",
			Expected: "exit:0",
		},
		{
			Name:     "Certbot installed",
			Command:  "certbot --version 2>&1",
			Expected: "exit:0",
		},
		{
			Name:     "Fail2ban installed",
			Command:  "fail2ban-client --version 2>&1",
			Expected: "exit:0",
		},
		{
			Name:     "Temporal config exists",
			Command:  "test -f /etc/temporal/temporal-server.yaml && echo exists",
			Expected: "exists",
		},
		{
			Name:     "Temporal UI config exists",
			Command:  "test -f /etc/temporal/temporal-ui-server.yaml && echo exists",
			Expected: "exists",
		},
		{
			Name:     "Temporal user exists",
			Command:  "id temporal",
			Expected: "exit:0",
		},
		{
			Name:     "Temporal systemd unit exists",
			Command:  "systemctl list-unit-files | grep temporal.service",
			Expected: "temporal.service",
		},
		{
			Name:     "Temporal UI systemd unit exists",
			Command:  "systemctl list-unit-files | grep temporal-ui.service",
			Expected: "temporal-ui.service",
		},
	}
}

// DeployedChecks returns all checks for a fully deployed server
// This includes all prebaked checks plus additional runtime checks
func DeployedChecks() []Check {
	checks := PrebakedChecks()

	deployedOnly := []Check{
		{
			Name:    "Temporal service running",
			Command: "systemctl is-active temporal",
			Equals:  "active",
		},
		{
			Name:    "Temporal UI service running",
			Command: "systemctl is-active temporal-ui",
			Equals:  "active",
		},
		{
			Name:    "Temporal Dunebot worker service running",
			Command: "systemctl is-active temporal-dunebot-worker",
			Equals:  "active",
		},
		{
			Name:    "DuneBot service running",
			Command: "systemctl is-active dunebot",
			Equals:  "active",
		},
		{
			Name:    "Nginx service running",
			Command: "systemctl is-active nginx",
			Equals:  "active",
		},
		{
			Name:    "Fail2ban service running",
			Command: "systemctl is-active fail2ban",
			Equals:  "active",
		},
		{
			Name:     "Port 7236 listening (Temporal gRPC)",
			Command:  "ss -tlnp | grep 7236",
			Expected: "LISTEN",
		},
		{
			Name:     "Port 8233 listening (Temporal UI)",
			Command:  "ss -tlnp | grep 8233",
			Expected: "LISTEN",
		},
		{
			Name:     "Port 80 listening (HTTP)",
			Command:  "ss -tlnp | grep :80",
			Expected: "LISTEN",
		},
		{
			Name:     "Port 443 listening (HTTPS)",
			Command:  "ss -tlnp | grep :443",
			Expected: "LISTEN",
		},
		{
			Name:     "SSL certificate exists",
			Command:  "find /etc/letsencrypt/live -name fullchain.pem 2>/dev/null | grep -q . && echo exists",
			Expected: "exists",
		},
		{
			Name:     "Temporal namespace available",
			Command:  "temporal --address localhost:7236 operator namespace list 2>&1",
			Expected: "default",
		},
	}

	return append(checks, deployedOnly...)
}

// WireguardPrebakedChecks returns the checks for verifying a pre-baked WireGuard image.
// These verify that all binaries, packages, and systemd units were installed
// during the Packer build without requiring runtime configuration.
func WireguardPrebakedChecks() []Check {
	return []Check{
		// Binaries
		{
			Name:     "WireGuard UI is valid ELF binary",
			Command:  "head -c 4 /usr/local/bin/wireguard-ui | od -An -tx1 | grep -q '7f 45 4c 46' && echo 'valid ELF'",
			Expected: "valid ELF",
		},
		{
			Name:     "Secret Operator Client is valid ELF binary",
			Command:  "head -c 4 /usr/bin/secret-operator-client | od -An -tx1 | grep -q '7f 45 4c 46' && echo 'valid ELF'",
			Expected: "valid ELF",
		},
		{
			Name:     "Secret Operator Server is valid ELF binary",
			Command:  "head -c 4 /usr/bin/secret-operator-server | od -An -tx1 | grep -q '7f 45 4c 46' && echo 'valid ELF'",
			Expected: "valid ELF",
		},

		// Packages
		{
			Name:     "WireGuard tools installed",
			Command:  "which wg 2>/dev/null && echo installed || echo missing",
			Expected: "installed",
		},
		{
			Name:     "wg-quick installed",
			Command:  "which wg-quick 2>/dev/null && echo installed || echo missing",
			Expected: "installed",
		},
		{
			Name:     "Nginx installed",
			Command:  "nginx -v 2>&1",
			Expected: "exit:0",
		},
		{
			Name:     "Certbot installed",
			Command:  "certbot --version 2>&1",
			Expected: "exit:0",
		},
		{
			Name:     "Fail2ban installed",
			Command:  "fail2ban-client --version 2>&1",
			Expected: "exit:0",
		},
		{
			Name:     "pwgen installed",
			Command:  "which pwgen 2>/dev/null && echo installed || echo missing",
			Expected: "installed",
		},

		// Configuration directories
		{
			Name:     "WireGuard config directory exists",
			Command:  "test -d /etc/wireguard && echo exists",
			Expected: "exists",
		},
		{
			Name:     "WireGuard UI data directory exists",
			Command:  "test -d /usr/local/bin/db/server && test -d /usr/local/bin/db/users && echo exists",
			Expected: "exists",
		},
		{
			Name:     "Letsencrypt live directory exists",
			Command:  "test -d /etc/letsencrypt/live && echo exists",
			Expected: "exists",
		},

		// Systemd units
		{
			Name:     "wgui systemd unit exists",
			Command:  "systemctl list-unit-files | grep wgui.service",
			Expected: "wgui.service",
		},
		{
			Name:     "wgui-restart-wg.service unit exists",
			Command:  "systemctl list-unit-files | grep wgui-restart-wg.service",
			Expected: "wgui-restart-wg.service",
		},
		{
			Name:     "wgui-restart-wg.path unit exists",
			Command:  "systemctl list-unit-files | grep wgui-restart-wg.path",
			Expected: "wgui-restart-wg.path",
		},
		{
			Name:     "secret-operator-server systemd unit exists",
			Command:  "systemctl list-unit-files | grep secret-operator-server.service",
			Expected: "secret-operator-server.service",
		},

		// AppArmor / sysctl / UFW pre-configuration
		{
			Name:     "IP forwarding enabled persistently",
			Command:  "grep -q 'net.ipv4.ip_forward=1' /etc/sysctl.d/99-wireguard.conf && echo enabled",
			Expected: "enabled",
		},
		{
			Name:     "UFW forward policy set to ACCEPT",
			Command:  "grep -q 'DEFAULT_FORWARD_POLICY=\"ACCEPT\"' /etc/default/ufw && echo configured",
			Expected: "configured",
		},
		{
			Name:     "wg-quick AppArmor local profile exists",
			Command:  "test -f /etc/apparmor.d/local/wg-quick && echo exists",
			Expected: "exists",
		},
	}
}

// WireguardDeployedChecks returns all checks for a fully deployed WireGuard server.
// This includes all prebaked checks plus additional runtime checks.
func WireguardDeployedChecks() []Check {
	checks := WireguardPrebakedChecks()

	deployedOnly := []Check{
		// Services running
		{
			Name:    "wg-quick@wg0 service running",
			Command: "systemctl is-active wg-quick@wg0",
			Equals:  "active",
		},
		{
			Name:    "wgui service running",
			Command: "systemctl is-active wgui",
			Equals:  "active",
		},
		{
			Name:    "secret-operator-server service running",
			Command: "systemctl is-active secret-operator-server",
			Equals:  "active",
		},
		{
			Name:    "Nginx service running",
			Command: "systemctl is-active nginx",
			Equals:  "active",
		},
		{
			Name:    "Fail2ban service running",
			Command: "systemctl is-active fail2ban",
			Equals:  "active",
		},

		// Listening ports
		{
			Name:     "Port 80 listening (HTTP)",
			Command:  "ss -tlnp | grep :80",
			Expected: "LISTEN",
		},
		{
			Name:     "Port 443 listening (HTTPS)",
			Command:  "ss -tlnp | grep :443",
			Expected: "LISTEN",
		},
		{
			Name:     "Port 8080 listening (WireGuard UI)",
			Command:  "ss -tlnp | grep :8080",
			Expected: "LISTEN",
		},
		{
			Name:     "Port 51820 listening (WireGuard UDP)",
			Command:  "ss -ulnp | grep :51820",
			Expected: "51820",
		},

		// WireGuard configuration
		{
			Name:     "wg0.conf exists",
			Command:  "test -f /etc/wireguard/wg0.conf && echo exists",
			Expected: "exists",
		},
		{
			Name:     "WireGuard interface up",
			Command:  "ip link show wg0 2>/dev/null | grep -E -q 'state (UNKNOWN|UP)' && echo up || echo down",
			Expected: "up",
		},
		{
			Name:     "SSL certificate exists",
			Command:  "find /etc/letsencrypt/live -name fullchain.pem 2>/dev/null | grep -q . && echo exists",
			Expected: "exists",
		},
		{
			Name:     "Nginx router site enabled",
			Command:  "test -L /etc/nginx/conf.d/router.conf && echo enabled",
			Expected: "enabled",
		},
		{
			Name:     "Nginx router.secret site enabled",
			Command:  "test -L /etc/nginx/conf.d/router-secret.conf && echo enabled",
			Expected: "enabled",
		},
		{
			Name:     "UFW active",
			Command:  "ufw status 2>/dev/null | grep -q 'Status: active' && echo active || echo inactive",
			Expected: "active",
		},
	}

	return append(checks, deployedOnly...)
}
