.PHONY: build
export PULUMI_CONFIG_PASSPHRASE ?= test
#STACK_SUFFIX ?="-$(shell pwgen -s 8 1)"
PROJECT ?= wireguard
CLOUD ?= aws
STACK_NAME ?= ${PROJECT}-${CLOUD}${STACK_SUFFIX}
AWS_REGION ?= eu-west-1
WIREGUARD_SERVER_IP=$(shell pulumi stack output publicIp)
WIREGUARD_SERVER_PUBLIC_KEY=$(shell pulumi stack output wireguard.publicKey)
SSH_USER ?= ubuntu

PRIVATE_KEY_FILE ?= ./keys/wireguard.pem
TMP_FOLDER ?= ./test/tmp
BUILD_FOLDER ?= $(PWD)

# Pulumi Configuration
export VPN_ENABLED_SSH ?= true
export CLIENT_IP_ADDRESS ?= 10.8.0.3
export CLIENT_PUBLICKEY ?= 872SDXKUNDyF7iE9qrfvi96rXgkPVN0b+MOHMAqcNFg=
export METADATA_URL ?= http://169.254.169.254/latest/meta-data/public-ipv4

go-init:
	go mod init github.com/fr123k/pulumi-wireguard-aws
	go mod vendor

pulumi-init: build
	pulumi plugin install resource aws 4.14.0
	pulumi plugin ls
	pulumi login --local
	# pulumi login --cloud-url s3://s3-pulumi-state-d12f2f1
	# pulumi stack rm -f ${STACK_NAME}
	# pulumi stack select ${STACK_NAME}
	pulumi stack select -c ${STACK_NAME}
	pulumi config set aws:region eu-west-1
	pulumi config set vpn_enabled_ssh ${VPN_ENABLED_SSH}

init: pulumi-init

build:
	go build -o ${BUILD_FOLDER}/build/${PROJECT}-${CLOUD} cmd/${PROJECT}/${CLOUD}/${PROJECT}.go
	go test -v --cover ./...
	mkdir -p ./build
	ln -fs ${BUILD_FOLDER}/build/${PROJECT}-${CLOUD} ./build/wireguard

create: pulumi-init
	pulumi up --yes
	# verbose logging
	# pulumi up --yes --logtostderr -v=9 2> out.txt

clean:
	pulumi destroy --yes -s ${STACK_NAME}
	pulumi stack rm -f --yes ${STACK_NAME} || true

recreate: clean create output

deploy: init create output

local: local-cleanup deploy

shell:
	pulumi stack output publicDns
	ssh -i "${PRIVATE_KEY_FILE}" -v ${SSH_USER}@${WIREGUARD_SERVER_IP}

browse:
	pulumi stack output publicDns
	open http://$(shell pulumi stack output publicDns)

output:
	mkdir -p ./output
	pulumi stack output --json > ./output/wireguard-ec2.json

## wireguard

prepare:
	mkdir -p ${TMP_FOLDER}

wireguard-client-keys: prepare
	wg genkey | tee ${TMP_FOLDER}/client_privatekey | wg pubkey > ${TMP_FOLDER}/client_publickey

wireguard-public-key: prepare
	echo "${WIREGUARD_SERVER_PUBLIC_KEY}" > ${TMP_FOLDER}/server_publickey

validate: wireguard-public-key
	$(MAKE) -C test -e WIREGUARD_SERVER_IP=${WIREGUARD_SERVER_IP} -e TMP_FOLDER=${TMP_FOLDER} wireguard-client
