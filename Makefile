.PHONY: build
export PULUMI_CONFIG_PASSPHRASE ?= test
#STACK_SUFFIX ?="-$(shell pwgen -s 8 1)"
PROJECT ?= wireguard
VM ?= ${PROJECT}
CLOUD ?= aws
STACK_NAME ?= ${VM}-${CLOUD}${STACK_SUFFIX}
AWS_REGION ?= eu-west-1
WIREGUARD_SERVER_IP=$(shell pulumi stack output publicIp)
# WIREGUARD_SERVER_IP=78.47.97.138
# WIREGUARD_SERVER_PUBLIC_KEY=$(shell pulumi stack output wireguard.publicKey)
SSH_USER ?= ubuntu

PRIVATE_KEY_FILE ?= ./keys/id_rsa_wireguard
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
	pulumi plugin install resource aws 4.38.1
	pulumi plugin install resource hcloud 1.20.5
	pulumi plugin ls
	pulumi login --local
	# pulumi login --cloud-url s3://s3-pulumi-state-d12f2f1
	# pulumi stack rm -f ${STACK_NAME}
	pulumi stack init ${STACK_NAME} || echo ignore if stack ${STACK_NAME} already exists
	pulumi stack select -c ${STACK_NAME} 
	pulumi config set aws:region eu-west-1
	pulumi config set vpn_enabled_ssh ${VPN_ENABLED_SSH}
	pulumi config set ssh_key_file ${PRIVATE_KEY_FILE}

init: pulumi-init

build:
	go build -o ${BUILD_FOLDER}/build/${PROJECT}-${CLOUD} cmd/${PROJECT}/${CLOUD}/${PROJECT}.go
	go test -v -timeout 60s --cover ./...
	mkdir -p ./build
	ln -fs ${BUILD_FOLDER}/build/${PROJECT}-${CLOUD} ./build/wireguard

verify:
	go build -o ./build/verify ./cmd/verify/

verify-linux:
	GOOS=linux GOARCH=amd64 go build -o ./build/verify-linux ./cmd/verify/

create: pulumi-init
	pulumi up --yes
	# verbose logging
	# pulumi up --yes --logtostderr -v=9 2> out.txt

preview: pulumi-init
	pulumi preview --diff

clean:
	pulumi destroy --yes -s ${STACK_NAME}
	pulumi stack rm -f --yes ${STACK_NAME} || true

recreate: clean create output

deploy: init create output

local: local-cleanup deploy

shell:
	# pulumi stack output publicDns
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

validate-wireguard: wireguard-public-key
	$(MAKE) -C test -e WIREGUARD_SERVER_IP=${WIREGUARD_SERVER_IP} -e TMP_FOLDER=${TMP_FOLDER} wireguard-client

validate-jenkins:
	echo "valid"

## Packer targets for pre-baked Temporal images

PACKER_DIR ?= packer/hetzner/temporal
PACKER_MANIFEST ?= $(PACKER_DIR)/manifest.json
SNAPSHOT_KEEP_COUNT ?= 3

# Domain configuration (ENV=test for test domains)
ifeq ($(ENV),test)
  TEMPORAL_DOMAIN ?= temporal-test.dunebot.io
  DUNEBOT_DOMAIN ?= githubapp-test.dunebot.io
else
  TEMPORAL_DOMAIN ?= temporal.dunebot.io
  DUNEBOT_DOMAIN ?= githubapp.dunebot.io
endif

packer-init:
	cd $(PACKER_DIR) && packer init .

packer-validate: packer-init
	cd $(PACKER_DIR) && packer validate .

packer-build: packer-validate
	cd $(PACKER_DIR) && packer build .
	@echo "Build complete. Snapshot ID:"
	@jq -r '.builds[-1].artifact_id' $(PACKER_MANIFEST)

packer-build-debug: packer-validate
	cd $(PACKER_DIR) && PACKER_LOG=1 packer build -debug .

packer-cleanup:
	@echo "Cleaning up old snapshots (keeping last $(SNAPSHOT_KEEP_COUNT))..."
	@hcloud image list --selector packer_build=true --output json | \
		jq -r 'sort_by(.created) | reverse | .[$(SNAPSHOT_KEEP_COUNT):][] | .id' | \
		xargs -I {} hcloud image delete {}

packer-list:
	@hcloud image list --selector packer_build=true --output columns=id,description,created

## Temporal deployment with pre-baked image

temporal-set-snapshot:
	@if [ -z "$(SNAPSHOT_ID)" ]; then \
		SNAPSHOT_ID=$$(jq -r '.builds[-1].artifact_id' $(PACKER_MANIFEST)); \
	fi; \
	echo "Setting temporal_snapshot_id to $$SNAPSHOT_ID"; \
	pulumi config set temporal_snapshot_id $$SNAPSHOT_ID

temporal-deploy-prebaked: temporal-set-snapshot init
	# TEMPORAL_DOMAIN=$(TEMPORAL_DOMAIN) DUNEBOT_DOMAIN=$(DUNEBOT_DOMAIN) pulumi destroy
	TEMPORAL_DOMAIN=$(TEMPORAL_DOMAIN) DUNEBOT_DOMAIN=$(DUNEBOT_DOMAIN) pulumi refresh
	TEMPORAL_DOMAIN=$(TEMPORAL_DOMAIN) DUNEBOT_DOMAIN=$(DUNEBOT_DOMAIN) pulumi up --yes

temporal-deploy-base:
	pulumi config rm temporal_snapshot_id || true
	pulumi up --yes

## Full pipeline: build image and deploy

temporal-full-deploy: packer-build temporal-deploy-prebaked
	@echo "Full deployment complete!"

temporal-recreate-prebaked: clean packer-build temporal-deploy-prebaked

## Certificate Management

cert-generate-wildcard:
	@echo "Run certbot manually with DNS-01 challenge:"
	@echo "  sudo certbot certonly --manual --preferred-challenges dns -d '*.dunebot.io' -d 'dunebot.io'"
	@echo "Then store in GCP Secret Manager:"
	@echo "  gcloud secrets versions add dunebot-wildcard-cert --data-file=/etc/letsencrypt/live/dunebot.io/fullchain.pem"
	@echo "  gcloud secrets versions add dunebot-wildcard-key --data-file=/etc/letsencrypt/live/dunebot.io/privkey.pem"

cert-check-expiry:
	openssl s_client -connect temporal.dunebot.io:443 2>/dev/null | \
		openssl x509 -noout -dates

sync-versions:
	bash packer/hetzner/temporal/scripts/sync-versions.sh
