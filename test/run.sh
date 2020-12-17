#!/bin/sh -e

## The below is modified from https://github.com/activeeos/wireguard-docker
cat /etc/wireguard/wg0.conf
# Find a Wireguard interface
interfaces=`find /etc/wireguard -type f`
if [ -z $interfaces ]; then
    echo "$(date): Interface not found in /etc/wireguard" >&2
    exit 1
fi

start_interfaces() {
    for interface in $interfaces; do
        echo "$(date): Starting Wireguard $interface"
        wg-quick up $interface
    done
}

stop_interfaces() {
    for interface in $interfaces; do
        wg-quick down $interface
    done
}


# Add masquerade rule for NAT'ing VPN traffic bound for the Internet

if [ $IPTABLES_MASQ -eq 1 ]; then
    echo "Adding iptables NAT rule"
    iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
fi

# Handle shutdown behavior
finish () {
    echo "$(date): Shutting down Wireguard"
    stop_interfaces
    if [ $IPTABLES_MASQ -eq 1 ]; then
        iptables -t nat -D POSTROUTING -o eth0 -j MASQUERADE
    fi
    exit 0
}

check_network() {
    ping 1.1.1.1 -c 4
    dig +short google.de
    curl ipinfo.io/ip
}

check_network
start_interfaces
sleep 5
wg show
check_network
./wireguard-connection-validation.sh
finish
exit 0
