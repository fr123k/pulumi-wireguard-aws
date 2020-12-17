#!/bin/bash -e

EXTERNAL_IP=$(curl ipinfo.io/ip)
if [ "${WIREGUARD_SERVER_IP}" == "${EXTERNAL_IP}" ]; then
    echo "Extern ip equal wireguard server ip so VPN connection successful established."
    exit 0
fi
echo "External ip different from wireguard server so VPN connection failed."
exit -1
