export PULUMI_CONFIG_PASSPHRASE ?= test
#STACK_SUFFIX ?="-$(shell pwgen -s 8 1)"
STACK_NAME ?= wireguard-ec2${STACK_SUFFIX}
AWS_REGION ?= eu-west-1
WIREGUARD_SERVER_IP=$(shell pulumi stack output publicIp)

PRIVATE_KEY_FILE?=./keys/wireguard.pem
TMP_FOLDER?="./test/tmp"

go-init:
	go mod init main
	go mod vendor

pulumi-init: build
	pulumi plugin install resource aws 3.19.3
	pulumi plugin ls
	pulumi login --local
	# pulumi login --cloud-url s3://s3-pulumi-state-d12f2f1
	# pulumi stack rm -f ${STACK_NAME}
	# pulumi stack select ${STACK_NAME}
	pulumi stack select -c ${STACK_NAME}
	pulumi config set aws:region eu-west-1

init: pulumi-init

ssh-keygen:
	echo -e 'n\n' | ssh-keygen -t rsa -b 4096 -q -N "" -f ${PRIVATE_KEY_FILE} || true
	echo "No"

build: ssh-keygen
	go build -o $(shell basename $(shell pwd))

create: pulumi-init
	pulumi up --yes
	#verbose logging
	#pulumi up --yes --verbose 9 --logtostderr

clean:
	pulumi destroy --yes -s ${STACK_NAME} || true
	pulumi stack rm -f --yes ${STACK_NAME} || true

# local-cleanup:
# 	echo "ADMIN_PASSWORD = ${ADMIN_PASSWORD}"
# 	pulumi destroy --yes -s ${STACK_NAME} || true
# 	pulumi stack rm -f --yes ${STACK_NAME} || true

recreate: clean create output

deploy: init create output

travis: deploy
	sleep 120

local: local-cleanup deploy

# pre-shell: #check if the wireguard virtual machine exists
# 	terraform state show -state=terraform.tfstate module.wireguard.data.aws_instances.wireguards

shell:
	pulumi stack output publicDns
	ssh -i "${PRIVATE_KEY_FILE}" -v ubuntu@${WIREGUARD_SERVER_IP}

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
	@ssh -i "${PRIVATE_KEY_FILE}" -o "StrictHostKeyChecking no" ubuntu@${WIREGUARD_SERVER_IP} 'sudo cat /var/log/cloud-init-output.log'
	@ssh -i "${PRIVATE_KEY_FILE}" -o "StrictHostKeyChecking no" ubuntu@${WIREGUARD_SERVER_IP} 'sudo cat /tmp/server_publickey' > ${TMP_FOLDER}/server_publickey

validate: wireguard-public-key
	$(MAKE) -C test -e WIREGUARD_SERVER_IP=${WIREGUARD_SERVER_IP} -e TMP_FOLDER=${TMP_FOLDER} wireguard-client
