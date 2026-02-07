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
