#!/bin/bash

# Entrypoint for Docker-based mini PC test container
# Starts systemd as PID 1 so that systemctl works normally.

# Render the cloud-init template with defaults
if [ -f /tmp/minipc-template.txt ]; then
    echo "=== Rendering cloud-init template ==="
    sed -e 's/{{ MINIPC_USER }}//g' \
        -e 's/{{ MINIPC_SSH_PORT }}//g' \
        -e 's/{{ MINIPC_DOCKER }}/true/g' \
        -e 's/{{ MINIPC_NIC }}//g' \
        -e 's/{{ FRANKY_VERSION }}//g' \
        /tmp/minipc-template.txt > /tmp/minipc-rendered.sh
    chmod +x /tmp/minipc-rendered.sh
    echo "Cloud-init template rendered to /tmp/minipc-rendered.sh"
fi

# Ensure SSH is ready
ssh-keygen -A 2>/dev/null || true
mkdir -p /run/sshd

echo ""
echo "=== Mini PC Docker Test Container Ready ==="
echo ""
echo "  SSH:       ssh root@localhost -p 2222  (password: root)"
echo "             ssh vagrant@localhost -p 2222  (password: vagrant)"
echo "  Franky UI: http://localhost:8788  (guest:8787)"
echo "  Nginx:     http://localhost:8080"
echo ""
echo "  To run the cloud-init setup:"
echo "    docker exec -it <container> bash /tmp/minipc-rendered.sh"
echo ""

# Start systemd as PID 1 (required for systemctl to work)
exec /sbin/init