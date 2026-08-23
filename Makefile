.PHONY: build
export PULUMI_CONFIG_PASSPHRASE ?= test
#STACK_SUFFIX ?="-$(shell pwgen -s 8 1)"
PROJECT ?= wireguard
VM ?= ${PROJECT}
CLOUD ?= hetzner
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

go-init:
	go mod init github.com/fr123k/pulumi-wireguard-aws
	go mod vendor

pulumi-init-generic: build
	pulumi plugin install resource aws 7.35.0
	pulumi plugin install resource hcloud 1.39.0
	pulumi plugin ls
	pulumi login gs://containifyci-pulumi-state-backend
# 	pulumi login --local
	# pulumi stack rm -f ${STACK_NAME}
	pulumi stack init ${STACK_NAME} || echo ignore if stack ${STACK_NAME} already exists
	pulumi stack select -c ${STACK_NAME} 
	pulumi config set aws:region eu-west-1
	pulumi config set vpn_enabled_ssh ${VPN_ENABLED_SSH}
	pulumi config set ssh_key_file ${PRIVATE_KEY_FILE}
	pulumi config set domain $(DOMAIN)

wireguard-config: packer-set-snapshot

temporal-config: packer-set-snapshot

minipc-config:
	pulumi config set server_ip "${MINIPC_SERVER_IP}"
	pulumi config set username "${SSH_USER}"

pulumi-init: pulumi-init-generic ${PROJECT}-config

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
	pulumi up
	# verbose logging
	# pulumi up --yes --logtostderr -v=9 2> out.txt

preview: pulumi-init
	pulumi preview --diff

clean:
# 	echo "pulumi destroy ${STACK_NAME}"
	pulumi refresh --yes -s ${STACK_NAME}
	pulumi destroy --yes -s ${STACK_NAME}
	pulumi stack rm -f --yes ${STACK_NAME} || true

recreate: clean create output

deploy: pulumi-init create output

deploy-prebaked: pulumi-init packer-set-snapshot create output

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

## Packer targets for pre-baked images

# Both the hcloud (production) and VirtualBox (local testing) sources live in
# the same template. PACKER_PROVIDER selects which source is built via Packer's
# -only flag.
#
#   hetzner  -> -only=hcloud.<project>      (default)
#   vagrant  -> -only=vagrant.<project> (outputs .box directly)
#
# Usage:
#   make packer-build PROJECT=wireguard SECRET_OPERATOR_TOKEN=xxx
#   make packer-build PROJECT=wireguard PACKER_VARS_FILE=wireguard.pkrvars.hcl
#   make packer-build PROJECT=wireguard PACKER_PROVIDER=vagrant DOMAIN=wg.fr123k.uk SECRET_OPERATOR_TOKEN=xxx
PACKER_PROVIDER ?= hetzner
PACKER_ONLY ?= -only=hcloud.$(PROJECT)
PACKER_DIR ?= packer/hetzner/${PROJECT}
PACKER_MANIFEST ?= $(PACKER_DIR)/manifest.json
SNAPSHOT_KEEP_COUNT ?= 3

PACKER_VARS_FILE ?= 
ifneq ($(PACKER_VARS_FILE),)
  PACKER_VARS_FLAG = -var-file=$(PACKER_VARS_FILE)
else
  PACKER_VARS_FLAG = 
endif

# Secrets and config can be passed three ways (all equivalent, later ones win):
#   1. A .pkrvars.hcl file:  make packer-build ... PACKER_VARS_FILE=wireguard.pkrvars.hcl
#   2. Env vars (PKR_VAR_ prefix, read automatically by Packer):
#        PKR_VAR_secret_operator_token=xxx PKR_VAR_domain=wg.fr123k.uk make packer-build ...
#   3. Make args (forwarded as PKR_VAR_ env vars by the targets below):
#        make packer-build ... SECRET_OPERATOR_TOKEN=xxx DOMAIN=wg.fr123k.uk
ifdef SECRET_OPERATOR_TOKEN
  export PKR_VAR_secret_operator_token ?= $(SECRET_OPERATOR_TOKEN)
endif
ifdef DOMAIN
  export PKR_VAR_domain ?= $(DOMAIN)
endif

# Detect host arch once and map to Go's GOARCH for cross-compiling the
# verify binary. Used by both the virtualbox packer-build path and the
# cloudinit-test recipe. Packer reads PKR_VAR_<name> env vars as overrides.
HOST_ARCH := $(shell uname -m)
ifeq ($(HOST_ARCH),aarch64)
  PACKER_GOARCH = arm64
else ifeq ($(HOST_ARCH),arm64)
  PACKER_GOARCH = arm64
else ifeq ($(HOST_ARCH),x86_64)
  PACKER_GOARCH = amd64
else
  PACKER_GOARCH = amd64
endif
export PKR_VAR_target_goarch ?= $(PACKER_GOARCH)

# Map PACKER_PROVIDER to the Packer -only source selector.
ifeq ($(PACKER_PROVIDER),virtualbox)
  PACKER_ONLY = -only=virtualbox-ovf.$(PROJECT)
endif

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
	cd $(PACKER_DIR) && packer validate $(PACKER_ONLY) $(PACKER_VARS_FLAG) .

packer-build: packer-validate
	cd $(PACKER_DIR) && packer build -force $(PACKER_ONLY) $(PACKER_VARS_FLAG) .
	@echo "Build complete. Artifact ID:"
	@jq -r '.builds[-1].artifact_id' $(PACKER_MANIFEST)

packer-build-debug: packer-init
	cd $(PACKER_DIR) && PACKER_LOG=1 packer build -force -debug $(PACKER_ONLY) $(PACKER_VARS_FLAG) .

hetzner-packer-set-snapshot:
	@if [ -z "$(SNAPSHOT_ID)" ]; then \
		SNAPSHOT_ID=$$(jq -r '.builds[-1].artifact_id' $(PACKER_MANIFEST)); \
	fi; \
		echo "Setting snapshot_id to $$SNAPSHOT_ID"; \
	pulumi config set snapshot_id $$SNAPSHOT_ID

packer-set-snapshot: ${PACKER_PROVIDER}-packer-set-snapshot

vagrant-packer-set-snapshot:
	@echo "packer-set-snapshot is only available for PACKER_PROVIDER=hetzner (current: $(PACKER_PROVIDER))"

# virtualbox builds are local OVAs (not cloud snapshots), so there is no
# snapshot to register; mirror the vagrant stub to keep the dispatch valid.
virtualbox-packer-set-snapshot:
	@echo "packer-set-snapshot is a no-op for PACKER_PROVIDER=virtualbox (local OVA build)"

# Print the latest built OVA path from the project's manifest.json. Used to
# chain packer-build -> cloudinit-test without copy-pasting the timestamped
# OVA name. Requires jq.
packer-artifact:
	@jq -r '.builds[-1].files[0].name' $(PACKER_MANIFEST) | sed 's|^|$(PACKER_DIR)/output-$(PROJECT)/|'



## Cloud-Init Test (local VirtualBox testing of cloud-init scripts)
#
# Tests any cloud-init script by booting a base or pre-baked OVA and
# executing the rendered cloud-init on it. Requires the OVA and the
# cloud-init template to exist.
#
# Usage:
#   # Test pre-baked wireguard cloud-init
#   make cloudinit-test \
#     SOURCE_PATH=packer/hetzner/wireguard/output-wireguard/packer-wireguard-vagrant-20260822-113607.ova \
#     CLOUD_INIT_FILE=cloud-init/wireguard-prebaked.txt \
#     TARGET=wireguard MODE=deployed
#
#   # Test fresh wireguard (non-prebaked)
#   make cloudinit-test \
#     SOURCE_PATH=$(HOME)/.virtualbox/packer/packer_base_ubuntu_26.ova \
#     CLOUD_INIT_FILE=cloud-init/wireguard.txt \
#     TARGET=wireguard MODE=deployed
.PHONY: cloudinit-test
CLOUDINIT_DIR ?= packer/cloudinit-test
CLOUDINIT_SOURCE_PATH ?=
CLOUDINIT_CLOUD_INIT_FILE ?=
CLOUDINIT_DOMAIN ?= 
CLOUDINIT_CLIENT_PUBLICKEY ?= 872SDXKUNDyF7iE9qrfvi96rXgkPVN0b+MOHMAqcNFg=
CLOUDINIT_CLIENT_IP_ADDRESS ?= 172.16.16.2
CLOUDINIT_SECRET_OPERATOR_TOKEN ?=
CLOUDINIT_TARGET ?= wireguard
CLOUDINIT_MODE ?= deployed
CLOUDINIT_TEMPORAL_DOMAIN ?= temporal.dunebot.io
CLOUDINIT_DUNEBOT_DOMAIN ?= githubapp.dunebot.io

# If CLOUDINIT_SOURCE_PATH is empty, auto-resolve the latest OVA built by
# `make packer-build PROJECT=$(CLOUDINIT_TARGET) PACKER_PROVIDER=virtualbox`
# from that project's manifest.json. This lets you chain the two commands
# without copy-pasting the timestamped OVA name. Requires jq.
ifeq ($(CLOUDINIT_SOURCE_PATH),)
  CLOUDINIT_MANIFEST ?= packer/hetzner/$(CLOUDINIT_TARGET)/manifest.json
  CLOUDINIT_ARTIFACT := $(shell jq -r '.builds[-1].files[0].name // empty' $(CLOUDINIT_MANIFEST) 2>/dev/null)
  ifneq ($(CLOUDINIT_ARTIFACT),)
    CLOUDINIT_SOURCE_PATH := $(CURDIR)/packer/hetzner/$(CLOUDINIT_TARGET)/output-$(CLOUDINIT_TARGET)/$(CLOUDINIT_ARTIFACT)
  endif
endif

ifneq ($(CLOUDINIT_SOURCE_PATH),)
CLOUDINIT_VAR_SOURCE_PATH = -var source_path=$(CLOUDINIT_SOURCE_PATH)
else
CLOUDINIT_VAR_SOURCE_PATH =
endif

ifneq ($(CLOUDINIT_CLOUD_INIT_FILE),)
CLOUDINIT_VAR_CLOUD_INIT_FILE = -var cloud_init_file=$(CLOUDINIT_CLOUD_INIT_FILE)
else
CLOUDINIT_VAR_CLOUD_INIT_FILE =
endif

cloudinit-test:
	cd $(CLOUDINIT_DIR) && packer init . && packer build -force \
		$(CLOUDINIT_VAR_SOURCE_PATH) \
		$(CLOUDINIT_VAR_CLOUD_INIT_FILE) \
		-var target=$(CLOUDINIT_TARGET) \
		-var mode=$(CLOUDINIT_MODE) \
		-var domain=$(CLOUDINIT_DOMAIN) \
		-var target_goarch=$(PACKER_GOARCH) \
		-var dunebot_domain=$(CLOUDINIT_DUNEBOT_DOMAIN) \
		-var temporal_domain=$(CLOUDINIT_TEMPORAL_DOMAIN) \
		-var client_publickey=$(CLOUDINIT_CLIENT_PUBLICKEY) \
		-var client_ip_address=$(CLOUDINIT_CLIENT_IP_ADDRESS) \
		$(CLOUDINIT_SECRET_OPERATOR_TOKEN:%= -var secret_operator_token=%) \
		.

## Certificate Management

cert-generate-wildcard:
	@echo "Run certbot manually with DNS-01 challenge:"
	@echo "  For dunebot.io:"
	@echo "    sudo certbot certonly --manual --preferred-challenges dns -d '*.dunebot.io' -d 'dunebot.io'"
	@echo "  For fr123k.uk:"
	@echo "    sudo certbot certonly --manual --preferred-challenges dns -d '*.fr123k.uk' -d 'fr123k.uk'"
	@echo "Then store in GCP Secret Manager:"
	@echo "  gcloud secrets versions add dunebot-wildcard-cert --data-file=/etc/letsencrypt/live/dunebot.io/fullchain.pem"
	@echo "  gcloud secrets versions add dunebot-wildcard-key --data-file=/etc/letsencrypt/live/dunebot.io/privkey.pem"
	@echo "  gcloud secrets versions add fr123k-wildcard-cert --data-file=/etc/letsencrypt/live/fr123k.uk/fullchain.pem"
	@echo "  gcloud secrets versions add fr123k-wildcard-key --data-file=/etc/letsencrypt/live/fr123k.uk/privkey.pem"

cert-check-expiry:
	openssl s_client -connect temporal.dunebot.io:443 2>/dev/null | \
		openssl x509 -noout -dates

sync-versions:
	bash packer/hetzner/temporal/scripts/sync-versions.sh

sync-versions-wireguard:
	bash packer/hetzner/wireguard/scripts/sync-versions.sh

## Mini PC (physical server) targets

MINIPC_SERVER_IP ?=

minipc-verify: verify
	./build/verify --host "${MINIPC_SERVER_IP}" --key "${MINIPC_SSH_KEY_FILE}" --user "${MINIPC_SSH_USER}" --port ${MINIPC_SSH_PORT} --deployed --target minipc

minipc-shell:
	ssh -i "${MINIPC_SSH_KEY_FILE}" -p ${MINIPC_SSH_PORT} ${MINIPC_SSH_USER}@${MINIPC_SERVER_IP}

minipc-deploy-base:
	pulumi up --yes