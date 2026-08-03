export PULUMI_CONFIG_PASSPHRASE ?= test
#STACK_SUFFIX ?="-$(shell pwgen -s 8 1)"
PROJECT ?= wireguard
VM ?= ${PROJECT}
CLOUD ?= hetzner
STACK_NAME ?= ${VM}-${CLOUD}${STACK_SUFFIX}
AWS_REGION ?= eu-west-1
WIREGUARD_SERVER_IP = $(shell pulumi stack output publicIp)
SSH_USER ?= ubuntu
PRIVATE_KEY_FILE ?= ./keys/id_rsa_wireguard
TMP_FOLDER ?= ./test/tmp
BUILD_FOLDER ?= $(PWD)

export VPN_ENABLED_SSH ?= true
export CLIENT_IP_ADDRESS ?= 10.8.0.3
export CLIENT_PUBLICKEY ?= 872SDXKUNDyF7iE9qrfvi96rXgkPVN0b+MOHMAqcNFg=
export METADATA_URL ?= http://169.254.169.254/latest/meta-data/public-ipv4

# Domain configuration (merged ENV=test blocks)
ifeq ($(ENV),test)
  TEMPORAL_DOMAIN      ?= temporal-test.dunebot.io
  DUNEBOT_DOMAIN       ?= githubapp-test.dunebot.io
  FRANKY_DOMAIN        ?= franky-test.dunebot.io
  WIREGUARD_DOMAIN     ?= wg-test.fr123k.uk
  WIREGUARD_STACK_NAME ?= wireguard-test-hetzner
  PRIVATE_KEY_FILE      = ./keys/wireguard-test
else
  TEMPORAL_DOMAIN      ?= temporal.dunebot.io
  DUNEBOT_DOMAIN       ?= githubapp.dunebot.io
  FRANKY_DOMAIN        ?= franky.dunebot.io
  WIREGUARD_DOMAIN     ?= wg.fr123k.uk
  WIREGUARD_STACK_NAME ?= wireguard-hetzner
  PRIVATE_KEY_FILE      = ./keys/wireguard
endif

# Packer config
PACKER_DIR ?= packer/hetzner/${PROJECT}
PACKER_MANIFEST ?= $(PACKER_DIR)/manifest.json
SNAPSHOT_KEEP_COUNT ?= 3
PACKER_VARS_FILE ?=
ifneq ($(PACKER_VARS_FILE),)
  PACKER_VARS_FLAG = -var-file=$(PACKER_VARS_FILE)
else
  PACKER_VARS_FLAG =
endif
SECRET_OPERATOR_TOKEN ?=
ifeq ($(PROJECT),wireguard)
  export PKR_VAR_secret_operator_token = $(SECRET_OPERATOR_TOKEN)
  export PKR_VAR_wireguard_domain = $(WIREGUARD_DOMAIN)
endif

# Mini PC config
MINIPC ?= minipc
MINIPC_STACK_NAME ?= ${MINIPC}-local
MINIPC_SERVER_IP ?=
MINIPC_SSH_PORT ?= 22
MINIPC_PACKER_DIR ?= packer/local/${MINIPC}
MINIPC_PACKER_MANIFEST ?= $(MINIPC_PACKER_DIR)/manifest.json

# Per-service table: name | config_key | manifest | init | deploy_env | stack | output | packer_build | init_body
SERVICES := temporal wireguard minipc
RECREATE_SERVICES := temporal wireguard

temporal_CONFIG       = temporal_snapshot_id
temporal_MANIFEST     = $(PACKER_MANIFEST)
temporal_INIT         = init
temporal_DEPLOY_ENV   = TEMPORAL_DOMAIN=$(TEMPORAL_DOMAIN) DUNEBOT_DOMAIN=$(DUNEBOT_DOMAIN)
temporal_PACKER_BUILD = packer-build

wireguard_CONFIG       = wireguard_snapshot_id
wireguard_MANIFEST     = $(PACKER_MANIFEST)
wireguard_INIT         = wireguard-init
wireguard_DEPLOY_ENV   = WIREGUARD_DOMAIN=$(WIREGUARD_DOMAIN)
wireguard_PACKER_BUILD = packer-build
wireguard_STACK        = $(WIREGUARD_STACK_NAME)
wireguard_OUTPUT       = ./output/wireguard-$(WIREGUARD_STACK_NAME).json
wireguard_INIT_BODY    = pulumi config set aws:region eu-west-1 && pulumi config set vpn_enabled_ssh ${VPN_ENABLED_SSH} && pulumi config set ssh_key_file ${PRIVATE_KEY_FILE} && pulumi config set wireguard_domain $(WIREGUARD_DOMAIN)

minipc_CONFIG       = minipc_image_id
minipc_MANIFEST     = $(MINIPC_PACKER_MANIFEST)
minipc_INIT         = minipc-init
minipc_DEPLOY_ENV   =
minipc_PACKER_BUILD = minipc-packer-build
minipc_STACK        = $(MINIPC_STACK_NAME)
minipc_OUTPUT       = ./output/minipc.json
minipc_INIT_BODY    = pulumi config set server_ip "${MINIPC_SERVER_IP}" && pulumi config set ssh_key_file "${PRIVATE_KEY_FILE}" && pulumi config set username "${SSH_USER}" && pulumi config set ssh_port ${MINIPC_SSH_PORT}

# Packer cluster config
cloud_PACKER_DIR      = $(PACKER_DIR)
cloud_PACKER_VARS     = $(PACKER_VARS_FLAG)
cloud_PACKER_MANIFEST = $(PACKER_MANIFEST)
cloud_PACKER_ECHO     = Snapshot ID
minipc_packer_PACKER_DIR      = $(MINIPC_PACKER_DIR)
minipc_packer_PACKER_VARS     =
minipc_packer_PACKER_MANIFEST = $(MINIPC_PACKER_MANIFEST)
minipc_packer_PACKER_ECHO     = Image info

.PHONY: build go-init pulumi-init pulumi-login init verify verify-linux create preview clean \
	recreate deploy local shell browse output prepare wireguard-client-keys wireguard-public-key \
	validate-wireguard validate-jenkins wireguard-set-domain wireguard-deploy wireguard-deploy-test \
	packer-cleanup packer-list cert-generate-wildcard cert-check-expiry sync-versions \
	sync-versions-wireguard minipc-verify minipc-shell minipc-keys minipc-destroy \
	$(foreach s,$(SERVICES),$(addprefix $s-,set-snapshot deploy-prebaked deploy-base full-deploy)) \
	$(foreach s,$(RECREATE_SERVICES),$s-recreate-prebaked) \
	$(foreach s,wireguard minipc,$(addprefix $s-,init create preview clean output)) \
	$(foreach c,packer minipc-packer,$(addprefix $c-,init validate build build-debug))

# Core build targets
go-init:
	go mod init github.com/fr123k/pulumi-wireguard-aws
	go mod vendor

build:
	go build -o ${BUILD_FOLDER}/build/${PROJECT}-${CLOUD} cmd/${PROJECT}/${CLOUD}/${PROJECT}.go
	go test -v -timeout 60s --cover ./...
	mkdir -p ./build
	ln -fs ${BUILD_FOLDER}/build/${PROJECT}-${CLOUD} ./build/wireguard

verify:
	go build -o ./build/verify ./cmd/verify/
verify-linux:
	GOOS=linux GOARCH=amd64 go build -o ./build/verify-linux ./cmd/verify/

# Pulumi login helper
pulumi-login:
	pulumi login gs://containifyci-pulumi-state-backend

# Default stack (pulumi) — kept explicit (unique plugin install steps)
pulumi-init: build pulumi-login
	pulumi plugin install resource aws 7.35.0
	pulumi plugin install resource hcloud 1.39.0
	pulumi plugin ls
	pulumi stack init ${STACK_NAME} || echo ignore if stack ${STACK_NAME} already exists
	pulumi stack select -c ${STACK_NAME}
	pulumi config set aws:region eu-west-1
	pulumi config set vpn_enabled_ssh ${VPN_ENABLED_SSH}
	pulumi config set ssh_key_file ${PRIVATE_KEY_FILE}

init: pulumi-init
create: pulumi-init
	pulumi up --yes
preview: pulumi-init
	pulumi preview --diff
clean:
	pulumi destroy --yes -s ${STACK_NAME}
	pulumi stack rm -f --yes ${STACK_NAME} || true
recreate: clean create output
deploy: init create output
local: local-cleanup deploy
shell:
	ssh -i "${PRIVATE_KEY_FILE}" -v ${SSH_USER}@${WIREGUARD_SERVER_IP}
browse:
	pulumi stack output publicDns
	open http://$(shell pulumi stack output publicDns)
output:
	mkdir -p ./output
	pulumi stack output --json > ./output/wireguard-ec2.json

# Wireguard helpers
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
wireguard-set-domain:
	@echo "Setting wireguard_domain to $(WIREGUARD_DOMAIN) for stack $(WIREGUARD_STACK_NAME)"
	pulumi config set wireguard_domain $(WIREGUARD_DOMAIN)
wireguard-deploy: wireguard-create wireguard-output
	@echo "Wireguard deployment complete! Stack: $(WIREGUARD_STACK_NAME) Domain: $(WIREGUARD_DOMAIN)"
wireguard-deploy-test:
	$(MAKE) wireguard-deploy ENV=test

# Packer snapshot management
packer-cleanup:
	@echo "Cleaning up old snapshots (keeping last $(SNAPSHOT_KEEP_COUNT))..."
	@hcloud image list --selector packer_build=true --output json | \
		jq -r 'sort_by(.created) | reverse | .[$(SNAPSHOT_KEEP_COUNT):][] | .id' | \
		xargs -I {} hcloud image delete {}
packer-list:
	@hcloud image list --selector packer_build=true --output columns=id,description,created

# Certificate management
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

# Version sync (unified)
SYNC_SERVICE ?= temporal
sync-versions:
	bash packer/hetzner/$(SYNC_SERVICE)/scripts/sync-versions.sh
sync-versions-wireguard:
	$(MAKE) sync-versions SYNC_SERVICE=wireguard

# Mini PC helpers
minipc-verify: verify
	./build/verify --host "${MINIPC_SERVER_IP}" --key "${MINIPC_SSH_KEY_FILE}" --user "${MINIPC_SSH_USER}" --port ${MINIPC_SSH_PORT} --deployed
minipc-shell:
	ssh -i "${MINIPC_SSH_KEY_FILE}" -p ${MINIPC_SSH_PORT} ${MINIPC_SSH_USER}@${MINIPC_SERVER_IP}
minipc-keys: prepare
	wg genkey | tee ${TMP_FOLDER}/minipc_client_privatekey | wg pubkey > ${TMP_FOLDER}/minipc_client_publickey

# ============================================================================
# Macros for generated targets
# ============================================================================

# Service init (wireguard, minipc — pulumi-init is explicit above)
define INIT_TEMPLATE
$(1)-init: build pulumi-login
	pulumi stack init $$($(1)_STACK) || echo ignore if stack $$($(1)_STACK) already exists
	pulumi stack select -c $$($(1)_STACK)
	$$($(1)_INIT_BODY)
endef

# Snapshot / deploy-prebaked / deploy-base (all services)
# Normalised: pulumi refresh --yes; DEPLOY_ENV on refresh+up; init before set-snapshot.
define SNAPSHOT_RULES
$(1)-set-snapshot:
	@if [ -z "$$(SNAPSHOT_ID)" ]; then \
		SNAPSHOT_ID=$$$$(jq -r '.builds[-1].artifact_id' $$($(1)_MANIFEST)); \
	fi; \
	echo "Setting $$($(1)_CONFIG) to $$$$SNAPSHOT_ID"; \
	pulumi config set $$($(1)_CONFIG) $$$$SNAPSHOT_ID

$(1)-deploy-prebaked: $$($(1)_INIT) $(1)-set-snapshot
	$$($(1)_DEPLOY_ENV) pulumi refresh --yes
	$$($(1)_DEPLOY_ENV) pulumi up --yes

$(1)-deploy-base:
	pulumi config rm $$($(1)_CONFIG) || true
	$$($(1)_DEPLOY_ENV) pulumi up --yes
endef

# Lifecycle: create / preview / clean / output (wireguard, minipc)
define LIFECYCLE_TEMPLATE
$(1)-create: $(1)-init
	$$($(1)_DEPLOY_ENV) pulumi up --yes
$(1)-preview: $(1)-init
	$$($(1)_DEPLOY_ENV) pulumi preview --diff
$(1)-clean:
	pulumi destroy --yes -s $$($(1)_STACK)
	pulumi stack rm -f --yes $$($(1)_STACK) || true
$(1)-output:
	mkdir -p ./output
	pulumi stack output --json > $$($(1)_OUTPUT)
endef

# Packer triplet (cloud + minipc). Normalised: manifest fallback; build-debug depends on init.
define PACKER_TEMPLATE
$(1)-init:
	cd $$($(2)_PACKER_DIR) && packer init .
$(1)-validate: $(1)-init
	cd $$($(2)_PACKER_DIR) && packer validate $$($(2)_PACKER_VARS) .
$(1)-build: $(1)-validate
	cd $$($(2)_PACKER_DIR) && packer build $$($(2)_PACKER_VARS) .
	@echo "Build complete. $$($(2)_PACKER_ECHO):"
	@jq -r '.builds[-1].artifact_id' $$($(2)_PACKER_MANIFEST) 2>/dev/null || echo "No manifest found"
$(1)-build-debug: $(1)-init
	cd $$($(2)_PACKER_DIR) && PACKER_LOG=1 packer build -debug $$($(2)_PACKER_VARS) .
endef

# Full pipeline (all services)
define FULL_DEPLOY_TEMPLATE
$(1)-full-deploy: $$($(1)_PACKER_BUILD) $(1)-deploy-prebaked
	@echo "Full $(1) deployment complete!"
endef

# Recreate prebaked (temporal, wireguard only)
define RECREATE_TEMPLATE
$(1)-recreate-prebaked: clean $$($(1)_PACKER_BUILD) $(1)-deploy-prebaked
endef

# ============================================================================
# Generated targets
# ============================================================================

$(eval $(call INIT_TEMPLATE,wireguard))
$(eval $(call INIT_TEMPLATE,minipc))
$(foreach s,$(SERVICES),$(eval $(call SNAPSHOT_RULES,$(s))))
$(eval $(call LIFECYCLE_TEMPLATE,wireguard))
$(eval $(call LIFECYCLE_TEMPLATE,minipc))
minipc-destroy: minipc-clean
$(eval $(call PACKER_TEMPLATE,packer,cloud))
$(eval $(call PACKER_TEMPLATE,minipc-packer,minipc_packer))
$(foreach s,$(SERVICES),$(eval $(call FULL_DEPLOY_TEMPLATE,$(s))))
$(foreach s,$(RECREATE_SERVICES),$(eval $(call RECREATE_TEMPLATE,$(s))))