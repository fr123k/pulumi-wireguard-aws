TMP_FOLDER?="./test/tmp"

build:
	docker build -f Dockerfile -t wireguard:local .

client-config:
	./scripts/wireguard-client-cfg.sh
	cat ./tmp/wg0.conf

wireguard-client: build client-config
	docker run \
		--privileged \
		--cap-add NET_ADMIN \
		--cap-add SYS_MODULE \
		-e WIREGUARD_SERVER_IP=${WIREGUARD_SERVER_IP} \
		-v $(PWD)/test/tmp/wg0.conf:/etc/wireguard/wg0.conf \
		wireguard:local
