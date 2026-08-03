# Makefile Redundancy Analysis & Improvement Plan

## Executive Summary

The 382-line Makefile at `/tmp/pulumi-wireguard-aws/Makefile` contains **massive structural redundancy**. Four distinct "stacks" (`temporal`, `franky`, `wireguard`, `minipc`) each replicate nearly identical target families (`*-init`, `*-set-snapshot`, `*-deploy-prebaked`, `*-deploy-base`, `*-create`, `*-preview`, `*-clean`, `*-output`). Additionally, the Packer toolchain is duplicated for the default path vs. `minipc`. With parameterized recipes and `$(eval)`-generated targets, this file can be reduced from **382 lines to approximately 140–170 lines** (~55% reduction) while preserving all behavior.

---

## 1. Identified Redundant Target Groups

### 1.1 The `*-set-snapshot` Quadruplet (Lines 168–173, 186–191, 258–263, 369–374)
**Targets:** `temporal-set-snapshot`, `franky-set-snapshot`, `wireguard-set-snapshot`, `minipc-set-snapshot`

All four targets follow this identical template, differing only in:
- The config key name (`temporal_snapshot_id`, `franky_snapshot_id`, `wireguard_snapshot_id`, `minipc_image_id`)
- The manifest file used (`PACKER_MANIFEST` vs `MINIPC_PACKER_MANIFEST` for minipc)

### 1.2 The `*-deploy-prebaked` Quadruplet (Lines 175–178, 193–196, 265–268, 376–378)
**Targets:** `temporal-deploy-prebaked`, `franky-deploy-prebaked`, `wireguard-deploy-prebaked`, `minipc-deploy-prebaked`

Template: `*-set-snapshot` + `*-init` prerequisites, then `pulumi refresh` + `pulumi up --yes`.
**Subtle differences:**
- `temporal-deploy-prebaked` exports `TEMPORAL_DOMAIN` and `DUNEBOT_DOMAIN`
- `wireguard-deploy-prebaked` exports `WIREGUARD_DOMAIN`
- `franky-deploy-prebaked` and `minipc-deploy-prebaked` export nothing

### 1.3 The `*-deploy-base` Quadruplet (Lines 180–182, 198–200, 270–272, 380–382)
**Targets:** `temporal-deploy-base`, `franky-deploy-base`, `wireguard-deploy-base`, `minipc-deploy-base`

Template: `pulumi config rm <service>_snapshot_id || true`, then `pulumi up --yes`.
**Subtle difference:**
- `wireguard-deploy-base` exports `WIREGUARD_DOMAIN`; others do not.

### 1.4 The `*-init` Triplet (Lines 28–40, 222–230, 312–320)
**Targets:** `pulumi-init`, `wireguard-init`, `minipc-init`

All log in to GCS backend, init/select stack, and set config keys. **Not identical:**
- `pulumi-init` (default/AWS): sets `aws:region`, `vpn_enabled_ssh`, `ssh_key_file`
- `wireguard-init`: sets same as default plus `wireguard_domain`
- `minipc-init`: sets `server_ip`, `ssh_key_file`, `username`, `ssh_port`

### 1.5 Lifecycle Targets: `*-create`, `*-preview`, `*-clean`/`*-destroy`, `*-output` (Lines 55–56, 60–61, 63–65, 81–83, 236–237, 239–240, 242–244, 246–248, 322–323, 325–326, 328–330, 332–334)
These are duplicated across the default stack, `wireguard`, and `minipc`.
- `create` / `wireguard-create` / `minipc-create`
- `preview` / `wireguard-preview` / `minipc-preview`
- `clean` / `wireguard-clean` / `minipc-destroy`
- `output` / `wireguard-output` / `minipc-output`

### 1.6 Packer Triplets (Lines 143–153, 350–359)
**Default Packer:** `packer-init`, `packer-validate`, `packer-build`
**Minipc Packer:** `minipc-packer-init`, `minipc-packer-validate`, `minipc-packer-build`

Identical commands; only `PACKER_DIR`/`MINIPC_PACKER_DIR` and the absence of `PACKER_VARS_FLAG` in minipc differ.

### 1.7 Full Pipeline Targets (Lines 204–207, 276–279, 366)
**Targets:** `temporal-full-deploy`, `temporal-recreate-prebaked`, `wireguard-full-deploy`, `wireguard-recreate-prebaked`, `minipc-full-deploy`

All follow: `[clean] [packer-build] <service>-deploy-prebaked`.

### 1.8 Domain Configuration Blocks (Lines 120–130, 212–220)
Two nearly identical `ifeq ($(ENV),test) ... else ... endif` blocks for domain defaults. The second block (lines 212–220) partially shadows the first (already sets `WIREGUARD_DOMAIN` in lines 124/129).

### 1.9 Sync-versions Pair (Lines 299–303)
`syncc-versions` and `sync-versions-wireguard` only differ in directory name (`temporal` vs `wireguard`).

---

## 2. Categorization of Duplication

| Category | Pattern | Services Affected | Lines |
|---|---|---|---|
| **A. Snapshot Config Setters** | `*-set-snapshot` | temporal, franky, wireguard, minipc | 168–173, 186–191, 258–263, 369–374 |
| **B. Pre-baked Deployments** | `*-deploy-prebaked` | temporal, franky, wireguard, minipc | 175–178, 193–196, 265–268, 376–378 |
| **C. Base Deployments** | `*-deploy-base` | temporal, franky, wireguard, minipc | 180–182, 198–200, 270–272, 380–382 |
| **D. Stack Init** | `*-init` | default, wireguard, minipc | 28–40, 222–230, 312–320 |
| **E. Lifecycle (CRUD)** | `*-create`, `*-preview`, `*-clean`, `*-output` | default, wireguard, minipc | scattered |
| **F. Packer Toolchain** | `packer-init/validate/build` | default, minipc | 143–153, 350–359 |
| **G. Full Pipeline** | `*-full-deploy`, `*-recreate-prebaked` | temporal, wireguard, minipc | 204–207, 276–279, 366 |
| **H. ENV Domain Blocks** | `ifeq ENV` domains | temporal/dunebot/franky/wireguard | 120–130, 212–220 |
| **I. Version Sync** | `sync-versions*` | temporal, wireguard | 299–303 |

---

## 3. Concrete Reduction Strategies & Before/After Sketches

### Strategy 1: Parameterized Snapshot Setters via `$(eval)` (Categories A, B, C)

**Before (20 lines):**
```makefile
temporal-set-snapshot:
	@if [ -z "$(SNAPSHOT_ID)" ]; then \
		SNAPSHOT_ID=$$(jq -r '.builds[-1].artifact_id' $(PACKER_MANIFEST)); \
	fi; \
	echo "Setting temporal_snapshot_id to $$SNAPSHOT_ID"; \
	pulumi config set temporal_snapshot_id $$SNAPSHOT_ID

temporal-deploy-prebaked: temporal-set-snapshot init
	TEMPORAL_DOMAIN=$(TEMPORAL_DOMAIN) DUNEBOT_DOMAIN=$(DUNEBOT_DOMAIN) pulumi refresh
	TEMPORAL_DOMAIN=$(TEMPORAL_DOMAIN) DUNEBOT_DOMAIN=$(DUNEBOT_DOMAIN) pulumi up --yes

temporal-deploy-base:
	pulumi config rm temporal_snapshot_id || true
	pulumi up --yes

# ... repeated for franky, wireguard, minipc
```

**After (~12 lines + eval):**
```makefile
# Configurable per service: name, config key, manifest, env exports
define SNAPSHOT_RULES
$(1)-set-snapshot:
	@if [ -z "$$$$(SNAPSHOT_ID)" ]; then \
		SNAPSHOT_ID=$$$$(jq -r '.builds[-1].artifact_id' $(3)); \
	fi; \
	echo "Setting $(2) to $$$$SNAPSHOT_ID"; \
	pulumi config set $(2) $$$$SNAPSHOT_ID

$(1)-deploy-prebaked: $(1)-set-snapshot $(4)
	$(5)
	pulumi refresh
	$(5) pulumi up --yes

$(1)-deploy-base:
	pulumi config rm $(2) || true
	$(5) pulumi up --yes
endef

# Usage:
$(eval $(call SNAPSHOT_RULES,temporal,temporal_snapshot_id,$(PACKER_MANIFEST),init,TEMPORAL_DOMAIN=$(TEMPORAL_DOMAIN) DUNEBOT_DOMAIN=$(DUNEBOT_DOMAIN)))
$(eval $(call SNAPSHOT_RULES,franky,franky_snapshot_id,$(PACKER_MANIFEST),init,))
$(eval $(call SNAPSHOT_RULES,wireguard,wireguard_snapshot_id,$(PACKER_MANIFEST),wireguard-init,WIREGUARD_DOMAIN=$(WIREGUARD_DOMAIN)))
$(eval $(call SNAPSHOT_RULES,minipc,minipc_image_id,$(MINIPC_PACKER_MANIFEST),minipc-init,))
```

**Lines saved:** ~24 lines (from ~32 to ~8, net ~24 after accounting for the macro definition).
**Risk:** **Medium**. Requires careful handling of `$$$$` escaping in recipes and verifying that the env-prefix string `$(5)` is empty or ends with a trailing space before `pulumi up`.

---

### Strategy 2: Unified Stack Lifecycle Targets via Pattern Rules (Categories D, E)

Many stack operations share the same structure: `<service>-<action>`. We can use **static pattern rules** or **parameterized `$(eval)`**.

**Before (~40 lines across default, wireguard, minipc):**
```makefile
create: pulumi-init
	pulumi up --yes

preview: pulumi-init
	pulumi preview --diff

clean:
	pulumi destroy --yes -s ${STACK_NAME}
	pulumi stack rm -f --yes ${STACK_NAME} || true

output:
	mkdir -p ./output
	pulumi stack output --json > ./output/wireguard-ec2.json

wireguard-create: wireguard-init
	WIREGUARD_DOMAIN=$(WIREGUARD_DOMAIN) pulumi up --yes

wireguard-preview: wireguard-init
	WIREGUARD_DOMAIN=$(WIREGUARD_DOMAIN) pulumi preview --diff

# ... etc
```

**After (~20 lines):**
```makefile
# Define stacks as a list of (name, stack_name, init_target, extra_env)
STACKS := \
  default:$(STACK_NAME):pulumi-init: \
  wireguard:$(WIREGUARD_STACK_NAME):wireguard-init:WIREGUARD_DOMAIN=$(WIREGUARD_DOMAIN) \
  minipc:$(MINIPC_STACK_NAME):minipc-init:

define STACK_RULES
$(1)-create: $(3)
	$(4) pulumi up --yes

$(1)-preview: $(3)
	$(4) pulumi preview --diff

$(1)-clean:
	pulumi destroy --yes -s $(2)
	pulumi stack rm -f --yes $(2) || true

$(1)-output:
	mkdir -p ./output
	pulumi stack output --json > ./output/$(1).json
endef

$(foreach s,$(STACKS),$(eval $(call STACK_RULES,$(word 1,$(subst :, ,$s)),$(word 2,$(subst :, ,$s)),$(word 3,$(subst :, ,$s)),$(word 4,$(subst :, ,$s)))))
```

**Gotcha:** The default stack currently uses bare target names (`create`, `preview`, `clean`, `output`) without a prefix. You must either:
1. Rename them to `default-create`, `default-preview`, etc., **OR**
2. Keep the bare names as aliases: `create: default-create`.

**Lines saved:** ~20 lines.
**Risk:** **Low to Medium**. Requires testing the `foreach/eval` syntax and ensuring stack name expansion happens correctly.

---

### Strategy 3: Unified Packer Targets via Variables (Category F)

**Before (~18 lines):**
```makefile
packer-init:
	cd $(PACKER_DIR) && packer init .

packer-validate: packer-init
	cd $(PACKER_DIR) && packer validate $(PACKER_VARS_FLAG) .

packer-build: packer-validate
	cd $(PACKER_DIR) && packer build $(PACKER_VARS_FLAG) .
	@echo "Build complete. Snapshot ID:"
	@jq -r '.builds[-1].artifact_id' $(PACKER_MANIFEST)

minipc-packer-init:
	cd $(MINIPC_PACKER_DIR) && packer init .

minipc-packer-validate: minipc-packer-init
	cd $(MINIPC_PACKER_DIR) && packer validate .

minipc-packer-build: minipc-packer-validate
	cd $(MINIPC_PACKER_DIR) && packer build .
	@echo "Build complete. Image info:"
	@jq -r '.builds[-1].artifact_id' $(MINIPC_PACKER_MANIFEST) 2>/dev/null || echo "No manifest found"
```

**After (~10 lines):**
```makefile
# Usage: make packer-build PACKER_DIR=packer/hetzner/wireguard PACKER_MANIFEST=...
#        make packer-build PACKER_DIR=packer/local/minipc  PACKER_MANIFEST=...

packer-init:
	cd $(PACKER_DIR) && packer init .

packer-validate: packer-init
	cd $(PACKER_DIR) && packer validate $(PACKER_VARS_FLAG) .

packer-build: packer-validate
	cd $(PACKER_DIR) && packer build $(PACKER_VARS_FLAG) .
	@echo "Build complete. Artifact:"
	@jq -r '.builds[-1].artifact_id' $(PACKER_MANIFEST) 2>/dev/null || echo "No manifest found"
```

You then call `make packer-build PACKER_DIR=$(MINIPC_PACKER_DIR) ...` instead of having dedicated `minipc-packer-*` targets. Alternatively, keep thin aliases:

```makefile
minipc-packer-init: ; $(MAKE) packer-init PACKER_DIR=$(MINIPC_PACKER_DIR) PACKER_VARS_FLAG=
```

**Lines saved:** ~10 lines.
**Risk:** **Low**. The `PACKER_VARS_FLAG` for minipc is already empty; overriding `PACKER_DIR` at invocation is standard Makefile practice.

---

### Strategy 4: Single Domain Configuration Block (Category H)

**Before (22 lines):**
Two separate `ifeq` blocks at lines 120–130 and 212–220, with `WIREGUARD_DOMAIN` duplicated.

**After (11 lines):**
```makefile
ifeq ($(ENV),test)
  TEMPORAL_DOMAIN  ?= temporal-test.dunebot.io
  DUNEBOT_DOMAIN   ?= githubapp-test.dunebot.io
  FRANKY_DOMAIN    ?= franky-test.dunebot.io
  WIREGUARD_DOMAIN ?= wg-test.fr123k.uk
  WIREGUARD_STACK_NAME ?= wireguard-test-hetzner
  PRIVATE_KEY_FILE  = ./keys/wireguard-test
else
  TEMPORAL_DOMAIN  ?= temporal.dunebot.io
  DUNEBOT_DOMAIN   ?= githubapp.dunebot.io
  FRANKY_DOMAIN    ?= franky.dunebot.io
  WIREGUARD_DOMAIN ?= wg.fr123k.uk
  WIREGUARD_STACK_NAME ?= wireguard-hetzner
  PRIVATE_KEY_FILE  = ./keys/wireguard
endif
```

**Lines saved:** ~11 lines (plus elimination of duplicate `WIREGUARD_DOMAIN ?= ...`).
**Risk:** **Low**. Purely moving existing lines. Ensure `WIREGUARD_STACK_NAME` and `PRIVATE_KEY_FILE` assignment semantics (`=` vs `?=`) are preserved.

---

### Strategy 5: Unified Sync Versions (Category I)

**Before (4 lines):**
```makefile
sync-versions:
	bash packer/hetzner/temporal/scripts/sync-versions.sh

sync-versions-wireguard:
	bash packer/hetzner/wireguard/scripts/sync-versions.sh
```

**After (2 lines):**
```makefile
SYNC_SERVICE ?= temporal
sync-versions:
	bash packer/hetzner/$(SYNC_SERVICE)/scripts/sync-versions.sh
```

**Lines saved:** 2 lines.
**Risk:** **Low**.

---

### Strategy 6: `.PHONY` Declaration Fix

**Before (line 1):**
```makefile
.PHONY: build
```

**After (~5 lines):**
```makefile
.PHONY: build init create preview clean recreate deploy local shell browse output \
        prepare wireguard-client-keys wireguard-public-key validate-wireguard validate-jenkins \
        packer-init packer-validate packer-build packer-build-debug packer-cleanup packer-list \
        temporal-set-snapshot temporal-deploy-prebaked temporal-deploy-base \
        franky-set-snapshot franky-deploy-prebaked franky-deploy-base \
        wireguard-init wireguard-set-domain wireguard-create wireguard-preview wireguard-clean \
        wireguard-output wireguard-deploy wireguard-deploy-test wireguard-set-snapshot \
        wireguard-deploy-prebaked wireguard-deploy-base wireguard-full-deploy \
        wireguard-recreate-prebaked cert-generate-wildcard cert-check-expiry sync-versions \
        sync-versions-wireguard minipc-init minipc-create minipc-preview minipc-destroy \
        minipc-output minipc-verify minipc-shell minipc-keys minipc-packer-init \
        minipc-packer-validate minipc-packer-build minipc-packer-build-debug \
        minipc-full-deploy minipc-set-snapshot minipc-deploy-prebaked minipc-deploy-base
```

**Risk:** **None** — correctness fix. Currently many phony targets are missing from `.PHONY`, which can cause stale-file bugs if artifacts with matching names exist.

---

### Strategy 7: Shared Full-Pipeline Macro (Category G)

**Before (~8 lines):**
```makefile
temporal-full-deploy: packer-build temporal-deploy-prebaked
	@echo "Full deployment complete!"

temporal-recreate-prebaked: clean packer-build temporal-deploy-prebaked

wireguard-full-deploy: packer-build wireguard-deploy-prebaked
	@echo "Full deployment complete!"

wireguard-recreate-prebaked: clean packer-build wireguard-deploy-prebaked

minipc-full-deploy: minipc-packer-build minipc-deploy-prebaked
	@echo "Full mini PC deployment complete!"
```

**After (~6 lines):**
```makefile
define FULL_DEPLOY
$(1)-full-deploy: $(2) $(1)-deploy-prebaked
	@echo "Full $(1) deployment complete!"

$(1)-recreate-prebaked: clean $(2) $(1)-deploy-prebaked
endef

$(eval $(call FULL_DEPLOY,temporal,packer-build))
$(eval $(call FULL_DEPLOY,wireguard,packer-build))
$(eval $(call FULL_DEPLOY,minipc,minipc-packer-build))
```

Or, if Strategy 3 is adopted (unified packer), then `minipc` also uses `packer-build` with overridden vars:

```makefile
$(eval $(call FULL_DEPLOY,temporal,packer-build))
$(eval $(call FULL_DEPLOY,wireguard,packer-build))
$(eval $(call FULL_DEPLOY,minipc,packer-build))
```

**Lines saved:** ~4 lines.
**Risk:** **Low**.

---

## 4. Consolidated Improvement Plan (Recommended Order)

### Phase 1: Low-Risk Cleanup (~60 lines saved)
1. **Merge domain `ifeq` blocks** (Strategy 4) — eliminates duplicate `WIREGUARD_DOMAIN`.
2. **Unify `sync-versions`** (Strategy 5).
3. **Fix `.PHONY`** (Strategy 6).
4. **Unify Packer targets** (Strategy 3) and replace `minipc-packer-*` with aliases/overrides.

### Phase 2: Medium-Risk Parameterization (~120 lines saved)
5. **Generate `*-set-snapshot`, `*-deploy-prebaked`, `*-deploy-base`** via `$(eval)` macro (Strategy 1). Handle the env-var prefix differences as a macro parameter.
6. **Generate stack lifecycle targets** (`*-create`, `*-preview`, `*-clean`, `*-output`) via `$(eval)` macro (Strategy 2). Keep bare aliases (`create`, `preview`, etc.) pointing to `default-*` variants for backward compatibility.
7. **Generate full-pipeline targets** via `$(eval)` macro (Strategy 7).

---

## 5. Line-Count Estimates & Risk Matrix

| Change | Lines Before | Lines After | Saved | Risk | Notes |
|---|---|---|---|---|---|
| Merge domain blocks | 22 | 11 | **11** | Low | Watch `=` vs `?=` |
| Fix `.PHONY` | 1 | 5 | −4 | None | Net *adds* lines but prevents bugs |
| Unify `sync-versions` | 4 | 2 | **2** | Low | |
| Unify Packer + aliases | 18 | 8 | **10** | Low | |
| Snapshot/deploy macros (Strat 1) | 32 | 8 | **24** | Medium | `$$$$` escaping, env prefixes |
| Lifecycle macros (Strat 2) | 40 | 20 | **20** | Medium | `foreach/eval` syntax |
| Full-deploy macros (Strat 7) | 8 | 6 | **2** | Low | |
| Misc whitespace/alignment | — | — | **~10** | Low | |
| **Total Estimated** | **~382** | **~160** | **~222** | | **~58% reduction** |

---

## 6. Correctness Gotchas & Edge Cases

### Gotcha A: Env-Var Prefix Differences in Deployments
- `temporal-deploy-prebaked` exports **two** domain variables; `wireguard-deploy-prebaked` exports **one**; `franky` and `minipc` export **none**.
- **Mitigation:** Pass an optional "environment prefix string" as the 5th argument to the deploy macro. When empty, no prefix is emitted.

### Gotcha B: `minipc-set-snapshot` Uses `minipc_image_id` (not `minipc_snapshot_id`)
- The config key naming is inconsistent (`image_id` vs `snapshot_id`).
- **Mitigation:** The macro must accept the config key as a free parameter (already designed in Strategy 1).

### Gotcha C: `clean` vs `wireguard-clean` vs `minipc-destroy` Naming Inconsistency
- The default uses `clean`; `wireguard` uses `wireguard-clean`; `minipc` uses `minipc-destroy`.
- **Mitigation:** Standardize on `<service>-destroy` or `<service>-clean` in generated targets, and add `clean: default-destroy` as a backward-compatible alias.

### Gotcha D: `minipc-packer-build-debug` Depends on `minipc-packer-validate`, But Default `packer-build-debug` Depends on `packer-init`
- The dependency graph differs slightly.
- **Mitigation:** Either align them (both depend on `packer-init` or both on `packer-validate`) or parameterize the dependency in a macro.

### Gotcha E: `pulumi-init` Installs AWS **and** Hetzner Plugins
- `pulumi-init` (default) installs plugins for **both** clouds. This may be intentional (shared backend state), but if services only need one plugin, do not split unless certain.
- **Mitigation:** Keep `pulumi-init` as a shared helper; do not fold it per-service.

### Gotcha F: `WIREGUARD_SERVER_IP` is Evaluated at Parse Time
```makefile
WIREGUARD_SERVER_IP=$(shell pulumi stack output publicIp)
```
This runs during Makefile parsing, which may fail if no stack is selected. This is **not** redundancy, but a latent bug.
- **Mitigation:** Use `=` (deferred) instead of `:=` (immediate), or make it a target-level shell command.

### Gotcha G: Commented-Out Code
Lines 10–11, 33–34, 176 contain commented-out commands and old IP addresses. These are harmless but add noise. Remove during refactor.

---

## 7. Final Recommended Structure (Sketch)

```makefile
.PHONY: build init create preview clean ...
export PULUMI_CONFIG_PASSPHRASE ?= test
PROJECT      ?= wireguard
VM           ?= $(PROJECT)
CLOUD        ?= hetzner
STACK_NAME   ?= $(VM)-$(CLOUD)$(STACK_SUFFIX)
AWS_REGION   ?= eu-west-1
SSH_USER     ?= ubuntu
PRIVATE_KEY_FILE ?= ./keys/id_rsa_wireguard
TMP_FOLDER   ?= ./test/tmp
BUILD_FOLDER ?= $(PWD)

# --- Merged domain config ---
ifeq ($(ENV),test)
  TEMPORAL_DOMAIN ?= temporal-test.dunebot.io
  DUNEBOT_DOMAIN  ?= githubapp-test.dunebot.io
  FRANKY_DOMAIN   ?= franky-test.dunebot.io
  WIREGUARD_DOMAIN ?= wg-test.fr123k.uk
  WIREGUARD_STACK_NAME ?= wireguard-test-hetzner
  PRIVATE_KEY_FILE = ./keys/wireguard-test
else
  TEMPORAL_DOMAIN ?= temporal.dunebot.io
  DUNEBOT_DOMAIN  ?= githubapp.dunebot.io
  FRANKY_DOMAIN   ?= franky.dunebot.io
  WIREGUARD_DOMAIN ?= wg.fr123k.uk
  WIREGUARD_STACK_NAME ?= wireguard-hetzner
  PRIVATE_KEY_FILE = ./keys/wireguard
endif

# --- Core build targets ---
build: ...
init: pulumi-init
create: default-create
preview: default-preview
clean: default-clean
output: default-output
deploy: init create output
recreate: clean create output

# --- Generated service targets via eval ---
# (insert macros and eval calls here)
```

This refactoring will make adding a **fifth service** (e.g., `grafana`, `monitoring`) a one-line `$(eval ...)` call instead of another 15-line copy-paste block.