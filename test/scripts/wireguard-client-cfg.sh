#!/bin/bash

cat > ./tmp/wg0.conf << EOF
[Interface]
Address = 10.8.0.2/32
PrivateKey = $(cat ./tmp/client_privatekey)
DNS = 1.1.1.1

[Peer]
PublicKey = $(cat ./tmp/server_publickey)
AllowedIPs = 0.0.0.0/0
Endpoint = ${WIREGUARD_SERVER_IP}:51820
PersistentKeepalive = 25
EOF
