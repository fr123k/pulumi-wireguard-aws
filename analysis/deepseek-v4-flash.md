

<details>
<summary>Let me read the Makefile first to understand its full structure.</summary>
</details>
Now let me also check for any additional patterns I might have missed:
Now I have the full picture. Let me produce the detailed analysis.

---

# Makefile Redundancy Analysis & Improvement Plan

**File:** `/tmp/pulumi-wireguard-aws/Makefile` (382 lines)  
**Services represented:** `temporal`, `franky`, `wireguard`, `minipc` (plus a generic `pulumi` prefix)

---

## 1. Inventory of All Repeated Patterns

| Pattern Family | # of Copies | Lines (approx) | Services |
|---|---|---|---|
| `*-init` (pulumi login + stack init/select + config sets) | 3 | 28–40, 222–230, 312–320 | `pulumi`, `wireguard`, `minipc` |
| `*-create` (pulumi up) | 3 | 55–58, 236–237, 322–323 | `pulumi`, `wireguard`, `minipc` |
| `*-preview` (pulumi preview) | 3 | 60–61, 239–240, 325–326 | `pulumi`, `wireguard`, `minipc` |
| `*-clean`/`*-destroy` (destroy + rm stack) | 3 | 63–65, 242–244, 328–330 | `pulumi`, `wireguard`, `minipc` |
| `*-output` (mkdir + stack output --json) | 3 | 81–83, 246–248, 332–334 | `pulumi`, `wireguard`, `minipc` |
| `*-set-snapshot` (read manifest + config set) | 4 | 168–173, 186–191, 258–263, 369–374 | `temporal`, `franky`, `wireguard`, `minipc` |
| `*-deploy-prebaked` (set-snapshot + init + refresh + up) | 4 | 175–178, 193–196, 265–268, 376–378 | `temporal`, `franky`, `wireguard`, `minipc` |
| `*-deploy-base` (config rm snapshot_id + up) | 4 | 180–182, 198–200, 270–272, 380–382 | `temporal`, `franky`, `wireguard`, `minipc` |
| Packer init/validate/build/build-debug triplet | 2 | 143–155, 350–362 | `packer`, `minipc-packer` |
| `*-full-deploy` (packer-build + deploy-prebaked) | 3 | 204–205, 276–277, 366–367 | `temporal`, `wireguard`, `minipc` |
| `*-recreate-prebaked` (clean + packer-build + deploy-prebaked) | 2 | 207, 279 | `temporal`, `wireguard` |
| `*-deploy` composite (init + create + output) | 2 | 69, 250–251 | `pulumi`, `wireguard` |
| `sync-versions` (bash script call) | 2 | 299–300, 302–303 | `temporal`, `wireguard` |

**Total duplicated lines (conservative): ~180 lines** out of 382 — nearly **47%** is boilerplate.

---

## 2. Categorization of Duplication

### Category A: Identical template, only service name + config key differ
These are **pure template** targets. The recipe body is identical except for a variable substitution.

| Template | Varying elements | Copies |
|---|---|---|
| `*-create` | `pulumi up` (some pass extra env vars) | 3 |
| `*-preview` | `pulumi preview --diff` (some pass extra env vars) | 3 |
| `*-clean`/`*-destroy` | `pulumi destroy --yes -s ${STACK_NAME}` + `pulumi stack rm` | 3 |
| `*-output` | `mkdir -p ./output` + `pulumi stack output --json > ./output/...` | 3 |
| `*-set-snapshot` | `PACKER_MANIFEST` path, config key name | 4 |
| `*-deploy-base` | config key name for `rm` | 4 |
| `sync-versions` | script path | 2 |

### Category B: Nearly identical, minor recipe differences
These share ~80%+ of their body but have one or two divergent lines.

| Template | Differences | Copies |
|---|---|---|
| `*-init` | Different config keys set; `wireguard-init` adds `wireguard_domain` | 3 |
| `*-deploy-prebaked` | `temporal` passes 2 env vars; `wireguard` passes 1; `franky`/`minipc` pass none; dependency order varies | 4 |
| Packer triplet | `minipc` uses `MINIPC_PACKER_DIR` (no vars flag); `minipc-packer-build` has extra `2>/dev/null` fallback; `minipc-packer-build-debug` depends on validate not init | 2 |
| `*-full-deploy` | `minipc` uses `minipc-packer-build` + `minipc-deploy-prebaked` | 3 |
| `*-deploy` composite | `deploy` depends on `init create output`; `wireguard-deploy` depends on `wireguard-create wireguard-output` | 2 |

### Category C: Unique / one-off targets (no duplication)
- `go-init` (line 24)
- `build` (line 43)
- `verify` (line 49)
- `verify-linux` (line 52)
- `recreate` (line 67)
- `local` (line 71)
- `shell` (line 73)
- `browse` (line 78)
- `prepare` (line 87)
- `wireguard-client-keys` (line 90)
- `wireguard-public-key` (line 93)
- `validate-wireguard` (line 96)
- `validate-jenkins` (line 99)
- `packer-cleanup` (line 157)
- `packer-list` (line 163)
- `wireguard-set-domain` (line 232)
- `wireguard-deploy-test` (line 253)
- `cert-generate-wildcard` (line 283)
- `cert-check-expiry` (line 295)
- `minipc-verify` (line 336)
- `minipc-shell` (line 339)
- `minipc-keys` (line 342)

---

## 3. Proposed Reduction Strategies

### Strategy 1: Parameterized `*-init` via a canned recipe + per-service config variables

**Problem:** Three `*-init` targets (lines 28, 222, 312) repeat the same pulumi login/stack init/select pattern with different config keys.

**Solution:** Define a shared `_init` helper that reads config keys from a service-specific variable.

```makefile
# Define per-service init config
SERVICE ?= pulumi

# Service-specific init variables
pulumi_STACK_NAME = $(STACK_NAME)
pulumi_INIT_CONFIG = aws:region=$(AWS_REGION) vpn_enabled_ssh=$(VPN_ENABLED_SSH) ssh_key_file=$(PRIVATE_KEY_FILE)

wireguard_STACK_NAME = $(WIREGUARD_STACK_NAME)
wireguard_INIT_CONFIG = aws:region=$(AWS_REGION) vpn_enabled_ssh=$(VPN_ENABLED_SSH) ssh_key_file=$(PRIVATE_KEY_FILE) wireguard_domain=$(WIREGUARD_DOMAIN)

minipc_STACK_NAME = $(MINIPC_STACK_NAME)
minipc_INIT_CONFIG = server_ip="$(MINIPC_SERVER_IP)" ssh_key_file="$(PRIVATE_KEY_FILE)" username="$(SSH_USER)" ssh_port=$(MINIPC_SSH_PORT)

# Shared init recipe
define INIT_TEMPLATE
$1-init: build
	pulumi login gs://containifyci-pulumi-state-backend
	pulumi stack init $$($1_STACK_NAME) || echo ignore if stack $$($1_STACK_NAME) already exists
	pulumi stack select -c $$($1_STACK_NAME)
	for kv in $$($1_INIT_CONFIG); do \
		pulumi config set $$(echo $$kv | tr '=' ' '); \
	done
endef

$(eval $(call INIT_TEMPLATE,pulumi))
$(eval $(call INIT_TEMPLATE,wireguard))
$(eval $(call INIT_TEMPLATE,minipc))
```

**Before (3 targets, ~30 lines):**
```makefile
pulumi-init: build
	pulumi plugin install resource aws 7.35.0
	pulumi plugin install resource hcloud 1.39.0
	pulumi plugin ls
	pulumi login gs://...
	pulumi stack init ${STACK_NAME} || echo ...
	pulumi stack select -c ${STACK_NAME}
	pulumi config set aws:region eu-west-1
	pulumi config set vpn_enabled_ssh ${VPN_ENABLED_SSH}
	pulumi config set ssh_key_file ${PRIVATE_KEY_FILE}

wireguard-init: build
	pulumi login gs://...
	pulumi stack init $(WIREGUARD_STACK_NAME) || echo ...
	pulumi stack select -c $(WIREGUARD_STACK_NAME)
	pulumi config set aws:region eu-west-1
	pulumi config set vpn_enabled_ssh ${VPN_ENABLED_SSH}
	pulumi config set ssh_key_file ${PRIVATE_KEY_FILE}
	pulumi config set wireguard_domain $(WIREGUARD_DOMAIN)

minipc-init: build
	pulumi login gs://...
	pulumi stack init ${MINIPC_STACK_NAME} || echo ...
	pulumi stack select -c ${MINIPC_STACK_NAME}
	pulumi config set server_ip "${MINIPC_SERVER_IP}"
	pulumi config set ssh_key_file "${PRIVATE_KEY_FILE}"
	pulumi config set username "${SSH_USER}"
	pulumi config set ssh_port ${MINIPC_SSH_PORT}
```

**After (~15 lines):** The three targets generated from one template + 6 lines of variable definitions.

**Lines saved:** ~15  
**Risk:** Low. The `pulumi-init` target also runs `pulumi plugin install` and `pulumi plugin ls` — those are **not** in `wireguard-init` or `minipc-init`. If you fold them into the template, `wireguard-init` and `minipc-init` would gain those steps. That's likely harmless (idempotent), but verify. If you want to keep them separate, add a `PULUMI_PLUGIN_INSTALL` variable.

---

### Strategy 2: Parameterized `*-set-snapshot` via a single template

**Problem:** Four nearly identical targets (lines 168, 186, 258, 369) differ only by the manifest path variable and the config key name.

**Solution:**

```makefile
# Per-service snapshot config
temporal_PACKER_MANIFEST = $(PACKER_MANIFEST)
temporal_SNAPSHOT_CONFIG_KEY = temporal_snapshot_id

franky_PACKER_MANIFEST = $(PACKER_MANIFEST)
franky_SNAPSHOT_CONFIG_KEY = franky_snapshot_id

wireguard_PACKER_MANIFEST = $(PACKER_MANIFEST)
wireguard_SNAPSHOT_CONFIG_KEY = wireguard_snapshot_id

minipc_PACKER_MANIFEST = $(MINIPC_PACKER_MANIFEST)
minipc_SNAPSHOT_CONFIG_KEY = minipc_image_id

define SET_SNAPSHOT_TEMPLATE
$1-set-snapshot:
	@if [ -z "$$(SNAPSHOT_ID)" ]; then \
		SNAPSHOT_ID=$$$$(jq -r '.builds[-1].artifact_id' $$($1_PACKER_MANIFEST)); \
	fi; \
	echo "Setting $$($1_SNAPSHOT_CONFIG_KEY) to $$$$SNAPSHOT_ID"; \
	pulumi config set $$($1_SNAPSHOT_CONFIG_KEY) $$$$SNAPSHOT_ID
endef

$(eval $(call SET_SNAPSHOT_TEMPLATE,temporal))
$(eval $(call SET_SNAPSHOT_TEMPLATE,franky))
$(eval $(call SET_SNAPSHOT_TEMPLATE,wireguard))
$(eval $(call SET_SNAPSHOT_TEMPLATE,minipc))
```

**Before (4 targets, ~24 lines):** Four copies of the same 6-line recipe.  
**After (~12 lines):** One template + 8 lines of variable definitions.  
**Lines saved:** ~12  
**Risk:** Low. All four are structurally identical.

---

### Strategy 3: Parameterized `*-deploy-prebaked` with env-var overrides

**Problem:** Four targets (lines 175, 193, 265, 376) share the same `pulumi refresh; pulumi up --yes` core but differ in:
- Dependency order (`temporal-set-snapshot init` vs `wireguard-init wireguard-set-snapshot` vs `minipc-set-snapshot minipc-init`)
- Extra env vars passed to pulumi commands

**Solution:** Use a template with an `ENV_VARS` variable.

```makefile
temporal_DEPLOY_ENV = TEMPORAL_DOMAIN=$(TEMPORAL_DOMAIN) DUNEBOT_DOMAIN=$(DUNEBOT_DOMAIN)
franky_DEPLOY_ENV =
wireguard_DEPLOY_ENV = WIREGUARD_DOMAIN=$(WIREGUARD_DOMAIN)
minipc_DEPLOY_ENV =

define DEPLOY_PREBAKED_TEMPLATE
$1-deploy-prebaked: $1-set-snapshot $1-init
	# pulumi destroy
	$$($1_DEPLOY_ENV) pulumi refresh
	$$($1_DEPLOY_ENV) pulumi up --yes
endef

$(eval $(call DEPLOY_PREBAKED_TEMPLATE,temporal))
$(eval $(call DEPLOY_PREBAKED_TEMPLATE,franky))
$(eval $(call DEPLOY_PREBAKED_TEMPLATE,wireguard))
$(eval $(call DEPLOY_PREBAKED_TEMPLATE,minipc))
```

**Before (4 targets, ~16 lines):** Four copies.  
**After (~10 lines):** One template + 4 lines of variable definitions.  
**Lines saved:** ~6  
**Risk:** Medium. The dependency order differs: `temporal` and `franky` use `temporal-set-snapshot init` (set-snapshot first), while `wireguard` uses `wireguard-init wireguard-set-snapshot` (init first). The template above uses `$1-set-snapshot $1-init` (set-snapshot first). If `wireguard`'s order matters, you'd need a separate `DEPLOY_PREBAKED_DEPS` variable. Also, `temporal-deploy-prebaked` has a commented-out `pulumi destroy` line that the others don't — that's cosmetic.

---

### Strategy 4: Parameterized `*-deploy-base` via template

**Problem:** Four targets (lines 180, 198, 270, 380) differ only by config key name and optional env vars.

**Solution:**

```makefile
temporal_BASE_CONFIG_KEY = temporal_snapshot_id
temporal_BASE_ENV =
franky_BASE_CONFIG_KEY = franky_snapshot_id
franky_BASE_ENV =
wireguard_BASE_CONFIG_KEY = wireguard_snapshot_id
wireguard_BASE_ENV = WIREGUARD_DOMAIN=$(WIREGUARD_DOMAIN)
minipc_BASE_CONFIG_KEY = minipc_image_id
minipc_BASE_ENV =

define DEPLOY_BASE_TEMPLATE
$1-deploy-base:
	pulumi config rm $$($1_BASE_CONFIG_KEY) || true
	$$($1_BASE_ENV) pulumi up --yes
endef

$(eval $(call DEPLOY_BASE_TEMPLATE,temporal))
$(eval $(call DEPLOY_BASE_TEMPLATE,franky))
$(eval $(call DEPLOY_BASE_TEMPLATE,wireguard))
$(eval $(call DEPLOY_BASE_TEMPLATE,minipc))
```

**Before (4 targets, ~12 lines):** Four copies.  
**After (~10 lines):** One template + 8 lines of variable definitions.  
**Lines saved:** ~2 (modest, but eliminates repetition).  
**Risk:** Low.

---

### Strategy 5: Parameterized Packer init/validate/build/build-debug triplet

**Problem:** Two packer triplets (lines 143–155 and 350–362) differ only by `PACKER_DIR` vs `MINIPC_PACKER_DIR` and the presence of `PACKER_VARS_FLAG`.

**Solution:**

```makefile
# Per-packer-set variables
cloud_PACKER_DIR = $(PACKER_DIR)
cloud_PACKER_MANIFEST = $(PACKER_MANIFEST)
cloud_PACKER_VARS = $(PACKER_VARS_FLAG)

minipc_PACKER_DIR = $(MINIPC_PACKER_DIR)
minipc_PACKER_MANIFEST = $(MINIPC_PACKER_MANIFEST)
minipc_PACKER_VARS =

define PACKER_TEMPLATE
$1-packer-init:
	cd $$($1_PACKER_DIR) && packer init .

$1-packer-validate: $1-packer-init
	cd $$($1_PACKER_DIR) && packer validate $$($1_PACKER_VARS) .

$1-packer-build: $1-packer-validate
	cd $$($1_PACKER_DIR) && packer build $$($1_PACKER_VARS) .
	@echo "Build complete. Snapshot ID:"
	@jq -r '.builds[-1].artifact_id' $$($1_PACKER_MANIFEST)

$1-packer-build-debug: $1-packer-init
	cd $$($1_PACKER_DIR) && PACKER_LOG=1 packer build -debug $$($1_PACKER_VARS) .
endef

$(eval $(call PACKER_TEMPLATE,cloud))
$(eval $(call PACKER_TEMPLATE,minipc))
```

**Before (2 sets, ~20 lines):** Two copies.  
**After (~14 lines):** One template + 6 lines of variable definitions.  
**Lines saved:** ~6  
**Risk:** Medium. Note that `minipc-packer-build` has an extra `2>/dev/null || echo "No manifest found"` fallback that the main `packer-build` lacks. Also `minipc-packer-build-debug` depends on `minipc-packer-validate` while the main `packer-build-debug` depends on `packer-init`. These are subtle differences that must be preserved via conditional variables.

---

### Strategy 6: Parameterized `*-create`, `*-preview`, `*-clean`, `*-output` via template

**Problem:** Three services each have create/preview/clean/output targets that are nearly identical.

**Solution:**

```makefile
pulumi_STACK_NAME_VAR = $(STACK_NAME)
pulumi_CREATE_ENV =
pulumi_PREVIEW_ENV =
pulumi_OUTPUT_FILE = ./output/wireguard-ec2.json

wireguard_STACK_NAME_VAR = $(WIREGUARD_STACK_NAME)
wireguard_CREATE_ENV = WIREGUARD_DOMAIN=$(WIREGUARD_DOMAIN)
wireguard_PREVIEW_ENV = WIREGUARD_DOMAIN=$(WIREGUARD_DOMAIN)
wireguard_OUTPUT_FILE = ./output/wireguard-$(WIREGUARD_STACK_NAME).json

minipc_STACK_NAME_VAR = $(MINIPC_STACK_NAME)
minipc_CREATE_ENV =
minipc_PREVIEW_ENV =
minipc_OUTPUT_FILE = ./output/minipc.json

define SERVICE_TEMPLATE
$1-create: $1-init
	$$($1_CREATE_ENV) pulumi up --yes

$1-preview: $1-init
	$$($1_PREVIEW_ENV) pulumi preview --diff

$1-clean:
	pulumi destroy --yes -s $$($1_STACK_NAME_VAR)
	pulumi stack rm -f --yes $$($1_STACK_NAME_VAR) || true

$1-output:
	mkdir -p ./output
	pulumi stack output --json > $$($1_OUTPUT_FILE)
endef

$(eval $(call SERVICE_TEMPLATE,pulumi))
$(eval $(call SERVICE_TEMPLATE,wireguard))
$(eval $(call SERVICE_TEMPLATE,minipc))
```

**Before (12 targets, ~24 lines):** Three copies of four targets.  
**After (~18 lines):** One template + 9 lines of variable definitions.  
**Lines saved:** ~6  
**Risk:** Low. The `pulumi` prefix targets (`create`, `preview`, `clean`, `output`) are the generic ones. The `wireguard` variants pass `WIREGUARD_DOMAIN` env var. The `minipc` variant is named `minipc-destroy` not `minipc-clean` — the template above uses `$1-clean` which would produce `minipc-clean`. You'd need an alias or rename the template output.

---

### Strategy 7: Parameterized `*-full-deploy` and `*-recreate-prebaked`

**Problem:** Three `*-full-deploy` targets (lines 204, 276, 366) and two `*-recreate-prebaked` targets (lines 207, 279) follow a pattern.

**Solution:**

```makefile
define FULL_DEPLOY_TEMPLATE
$1-full-deploy: $1-packer-build $1-deploy-prebaked
	@echo "Full deployment complete!"
endef

$(eval $(call FULL_DEPLOY_TEMPLATE,temporal))
$(eval $(call FULL_DEPLOY_TEMPLATE,wireguard))
$(eval $(call FULL_DEPLOY_TEMPLATE,minipc))

define RECREATE_PREBAKED_TEMPLATE
$1-recreate-prebaked: clean $1-packer-build $1-deploy-prebaked
endef

$(eval $(call RECREATE_PREBAKED_TEMPLATE,temporal))
$(eval $(call RECREATE_PREBAKED_TEMPLATE,wireguard))
```

**Before (5 targets, ~5 lines):** Trivial but repeated.  
**After (~6 lines):** Two templates.  
**Lines saved:** ~0 (neutral, but eliminates repetition).  
**Risk:** Low.

---

### Strategy 8: Parameterized `sync-versions` targets

**Problem:** Two targets (lines 299, 302) differ only by script path.

**Solution:**

```makefile
temporal_SYNC_SCRIPT = packer/hetzner/temporal/scripts/sync-versions.sh
wireguard_SYNC_SCRIPT = packer/hetzner/wireguard/scripts/sync-versions.sh

define SYNC_TEMPLATE
sync-versions-$1:
	bash $$($1_SYNC_SCRIPT)
endef

$(eval $(call SYNC_TEMPLATE,temporal))
$(eval $(call SYNC_TEMPLATE,wireguard))
```

**Before (2 targets, ~4 lines):** Two copies.  
**After (~5 lines):** One template + 2 lines of variable definitions.  
**Lines saved:** ~0 (neutral).  
**Risk:** Low.

---

### Strategy 9: Composite `*-deploy` target via template

**Problem:** `deploy` (line 69) and `wireguard-deploy` (line 250) are similar composites.

**Solution:**

```makefile
define DEPLOY_TEMPLATE
$1-deploy: $1-create $1-output
	@echo "$1 deployment complete!"
endef

$(eval $(call DEPLOY_TEMPLATE,pulumi))
$(eval $(call DEPLOY_TEMPLATE,wireguard))
```

**Before (2 targets, ~4 lines):** Two copies.  
**After (~4 lines):** One template.  
**Lines saved:** ~0 (neutral).  
**Risk:** Low. Note: `deploy` currently depends on `init create output` (3 deps) while `wireguard-deploy` depends on `wireguard-create wireguard-output` (2 deps, since `wireguard-create` already depends on `wireguard-init`). The template above uses the `wireguard` pattern (2 deps). This is actually cleaner — `init` is already a prerequisite of `create`.

---

## 4. Summary of Lines Saved

| Strategy | Lines Before | Lines After | Lines Saved | Risk |
|---|---|---|---|---|
| 1. Parameterized `*-init` | ~30 | ~15 | **~15** | Low (plugin install diff) |
| 2. Parameterized `*-set-snapshot` | ~24 | ~12 | **~12** | Low |
| 3. Parameterized `*-deploy-prebaked` | ~16 | ~10 | **~6** | Medium (dep order) |
| 4. Parameterized `*-deploy-base` | ~12 | ~10 | **~2** | Low |
| 5. Parameterized Packer triplet | ~20 | ~14 | **~6** | Medium (minor recipe diffs) |
| 6. Parameterized create/preview/clean/output | ~24 | ~18 | **~6** | Low (naming diff) |
| 7. Parameterized full-deploy/recreate | ~5 | ~6 | **~0** | Low |
| 8. Parameterized sync-versions | ~4 | ~5 | **~0** | Low |
| 9. Composite deploy template | ~4 | ~4 | **~0** | Low |
| **Total** | **~139** | **~94** | **~47** | |

**Net reduction: ~47 lines** (from 382 to ~335, a ~12% reduction).  
If you also remove the now-redundant variable definitions that get consolidated, you could save another ~10–15 lines.

---

## 5. Correctness Gotchas & Subtle Differences

### ⚠️ `pulumi-init` has extra steps (lines 28–40)
Lines 29–31 run `pulumi plugin install resource aws 7.35.0`, `pulumi plugin install resource hcloud 1.39.0`, and `pulumi plugin ls`. These are **not** present in `wireguard-init` or `minipc-init`. If you fold all three into a template, the wireguard and minipc init targets will gain these steps. That's likely harmless (plugin install is idempotent), but it's a behavioral change.

### ⚠️ Dependency order in `*-deploy-prebaked`
- `temporal-deploy-prebaked`: `temporal-set-snapshot init` (set-snapshot **before** init)
- `franky-deploy-prebaked`: `franky-set-snapshot init` (set-snapshot **before** init)
- `wireguard-deploy-prebaked`: `wireguard-init wireguard-set-snapshot` (init **before** set-snapshot)
- `minipc-deploy-prebaked`: `minipc-set-snapshot minipc-init` (set-snapshot **before** init)

The order matters if `init` creates the stack and `set-snapshot` configures it. If the stack doesn't exist yet, `set-snapshot` would fail. The `wireguard` variant (init first) is arguably the correct order. The template should use `$1-init $1-set-snapshot` to match the safe order.

### ⚠️ `temporal-deploy-prebaked` passes two env vars
`TEMPORAL_DOMAIN=$(TEMPORAL_DOMAIN) DUNEBOT_DOMAIN=$(DUNEBOT_DOMAIN)` — the other services pass either one (`WIREGUARD_DOMAIN`) or none. This must be preserved via a per-service `DEPLOY_ENV` variable.

### ⚠️ `minipc-packer-build` has an extra fallback (line 359)
```makefile
@jq -r '.builds[-1].artifact_id' $(MINIPC_PACKER_MANIFEST) 2>/dev/null || echo "No manifest found"
```
The main `packer-build` (line 152) does not have the `2>/dev/null || echo` fallback. This must be preserved via a conditional or a separate variable.

### ⚠️ `minipc-packer-build-debug` depends on `minipc-packer-validate` (line 361)
The main `packer-build-debug` (line 154) depends on `packer-init`, not `packer-validate`. This is a real difference — the template must handle it.

### ⚠️ `minipc-destroy` vs `minipc-clean` naming
The minipc family uses `minipc-destroy` (line 328) while the others use `*-clean`. A template generating `$1-clean` would produce `minipc-clean`, breaking backward compatibility. Either rename the existing target or add an alias.

### ⚠️ Output file names differ
- `output` → `./output/wireguard-ec2.json`
- `wireguard-output` → `./output/wireguard-$(WIREGUARD_STACK_NAME).json`
- `minipc-output` → `./output/minipc.json`

These must be preserved via a per-service `OUTPUT_FILE` variable.

### ⚠️ `wireguard-deploy-test` (line 253) is a special wrapper
```makefile
wireguard-deploy-test:
	$(MAKE) wireguard-deploy ENV=test
```
This is a one-off convenience target — not worth templating.

### ⚠️ `.PHONY` declarations
The file has only `.PHONY: build` at line 1. All generated targets should be added to `.PHONY`. With `$(eval)`-generated targets, you can do:
```makefile
.PHONY: $(foreach s,pulumi wireguard minipc temporal franky cloud,$s-init $s-create $s-preview $s-clean $s-output $s-set-snapshot $s-deploy-prebaked $s-deploy-base $s-full-deploy $s-packer-init $s-packer-validate $s-packer-build $s-packer-build-debug)
```

### ⚠️ `minipc-verify` and `minipc-shell` use undefined variables
Lines 337 and 340 reference `MINIPC_SSH_KEY_FILE`, `MINIPC_SSH_USER` which are not defined anywhere in the Makefile. These are pre-existing bugs, not introduced by refactoring.

---

## 6. Recommended Implementation Order

1. **Start with Strategy 2** (`*-set-snapshot`) — lowest risk, highest savings.
2. **Strategy 4** (`*-deploy-base`) — trivial, safe.
3. **Strategy 6** (create/preview/clean/output) — careful with the `minipc-destroy` naming.
4. **Strategy 1** (`*-init`) — watch the plugin install difference.
5. **Strategy 5** (Packer triplet) — handle the `2>/dev/null` and dependency differences.
6. **Strategy 3** (`*-deploy-prebaked`) — handle env vars and dependency order.
7. **Strategies 7–9** (full-deploy, sync, composite) — trivial, do last.

Each strategy can be implemented independently and tested in isolation.