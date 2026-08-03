Now I have a comprehensive view of the Makefile. Let me also check whether there's a `wfranky-deploy-base` or `franky` config key pattern in the project, just to confirm the snapshot config keys used.
I have all the information I need. Now let me write a comprehensive, structured refactoring plan.

# Makefile Refactoring Plan — `/tmp/pulumi-wireguard-aws/Makefile`

**File:** 382 lines, ~70 named targets, only **one** `.PHONY` declaration (line 1, just `build`).
**Bottom-line estimate:** the file can be reduced to **~170–200 lines** (≈45–50 % reduction) with **low-to-moderate** risk if the changes below are applied incrementally.

---

## 1. Inventory of Duplicated Patterns

The Makefile is essentially the same Pulumi / Packer workflow cloned four times (once per "service"). The duplicated *templates* are:

### A. The **service triplet** (init / create / output [/ destroy / clean / preview / set-domain])

| Service  | init        | create         | preview       | clean/destroy   | output        | set-domain |
|----------|-------------|----------------|---------------|-----------------|---------------|------------|
| default  | `pulumi-init` (28) | `create` (55) | `preview` (60) | `clean` (63) | `output` (81) | – |
| wireguard | `wireguard-init` (222) | `wireguard-create` (236) | `wireguard-preview` (239) | `wireguard-clean` (242) | `wireguard-output` (246) | `wireguard-set-domain` (232) |
| minipc   | `minipc-init` (312) | `minipc-create` (322) | `minipc-preview` (325) | `minipc-destroy` (328) | `minipc-output` (332) | – |
| temporal | – (uses `init`) | – (via `temporal-deploy-prebaked`) | – | – | – | – |
| franky   | – (uses `init`) | – (via `franky-deploy-prebaked`) | – | – | – | – |

Each row of the table is the same recipe with the service name stitched into:
- the stack name,
- the pulumi config key(s),
- the output filename,
- a `pulumi login gs://containifyci-pulumi-state-backend` line,
- a `pulumi stack init …` line.

### B. The **pre-baked deploy** triplet (`set-snapshot`, `deploy-prebaked`, `deploy-base`)

These three targets exist for **four** services and are byte-for-byte identical except for the config key:

```text
<SVC>-set-snapshot:        pulumi config set <SVC>_snapshot_id $SNAPSHOT_ID
<SVC>-deploy-prebaked:     pulumi refresh [--yes]; pulumi up --yes
<SVC>-deploy-base:         pulumi config rm <SVC>_snapshot_id || true; pulumi up --yes
```

| Service   | set-snapshot | deploy-prebaked        | deploy-base      | config key           |
|-----------|--------------|------------------------|------------------|----------------------|
| temporal  | 168          | 175                    | 180              | `temporal_snapshot_id` |
| franky    | 186          | 193                    | 198              | `franky_snapshot_id`   |
| wireguard | 258          | 265                    | 270              | `wireguard_snapshot_id` |
| minipc    | 369          | 376                    | 380              | `minipc_image_id`     |

The minipc variant uses a different config key (`minipc_image_id`) and a *different* manifest file (`MINIPC_PACKER_MANIFEST`) — that asymmetry is the only thing keeping it from being a perfect 4-way clone.

### C. The **packer triplet** (`packer-init` / `packer-validate` / `packer-build` / `packer-build-debug`)

Two near-identical clusters:

1. **wireguard/temporal/franky** — all use `PACKER_DIR` + `PACKER_VARS_FLAG` (lines 143–155).
2. **minipc** — same three targets, but using `MINIPC_PACKER_DIR`, no vars file, different path (lines 350–362).

### D. The **full-pipeline** targets (`*-full-deploy`, `*-recreate-prebaked`)

| Service   | full-deploy        | recreate-prebaked          |
|-----------|--------------------|----------------------------|
| temporal  | 204                | 207                        |
| wireguard | 276                | 279                        |
| minipc    | 366                | – (missing)                |

These are one-liners that just chain other targets; pure boilerplate.

### E. Pulumi **init** blocks

- `pulumi-init` (28) and `wireguard-init` (222) and `minipc-init` (312) are all the same `login / stack init / stack select / config set …` recipe. The minipc one drops the AWS plugin install but is otherwise the same shape.

### F. Minor duplicates

- `sync-versions` (300) and `sync-versions-wireguard` (303) — two near-identical targets calling a `sync-versions.sh` script in different packer dirs.
- `wireguard-client-keys` (90) and `minipc-keys` (342) — same `wg genkey | tee … | wg pubkey > …` recipe, different filenames.
- `pulumi login gs://containifyci-pulumi-state-backend` repeated 3×.
- The `echo "Full … deployment complete!"` message repeated 3×.

---

## 2. Categorisation — Targets that Follow the Same Template

| Template                                    | Variants                                          | Differing parameters                            |
|---------------------------------------------|---------------------------------------------------|-------------------------------------------------|
| `service-init` (login + stack init + config) | `pulumi-init`, `wireguard-init`, `minipc-init`   | stack name, config keys, optional plugin install |
| `service-set-snapshot`                      | `temporal-`, `franky-`, `wireguard-`, `minipc-`  | pulumi config key, manifest file                |
| `service-deploy-prebaked`                   | same 4                                            | config key, optional `WIREGUARD_DOMAIN` env     |
| `service-deploy-base`                       | same 4                                            | config key                                      |
| `service-create` / `-preview`               | `wireguard-create/-preview`, `minipc-create/-preview`, top-level `create`/`preview` | env vars forwarded |
| `service-output`                            | `output`, `wireguard-output`, `minipc-output`     | output filename                                 |
| `service-clean` / `-destroy`                | `clean`, `wireguard-clean`, `minipc-destroy`      | stack name                                      |
| `packer-init/-validate/-build/-build-debug` | packer cluster + minipc-packer cluster            | packer dir, vars flag                           |
| `service-full-deploy` / `-recreate-prebaked`| `temporal-`, `wireguard-`, `minipc-` (partial)    | which underlying targets                        |
| `service-keys`                             | `wireguard-client-keys`, `minipc-keys`            | key filename                                    |
| `service-sync-versions`                     | `sync-versions`, `sync-versions-wireguard`        | script path                                     |

---

## 3. Concrete Reduction Strategies

### Strategy 1 — **Service table + canned recipes** *(biggest win, ~120 lines)*

Declare each service once in a table, then `$(eval)` generate the pre-baked triplet and the init/output/clean/create/preview targets.

#### Before (lines 168–200, 258–272, 369–382 — ~60 lines)

```make
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

## franky …  (franky-set-snapshot, franky-deploy-prebaked, franky-deploy-base)  # same shape
## wireguard …  (wireguard-set-snapshot, wireguard-deploy-prebaked, wireguard-deploy-base)
## minipc …  (minipc-set-snapshot, minipc-deploy-prebaked, minipc-deploy-base)  # uses minipc_image_id
```

#### After

```make
# Service table: name | config_key | stack_name | manifest | extra_env
SERVICES := temporal,franky,wireguard,minipc

define svc-set-snapshot =
$(1)-set-snapshot:
	@if [ -z "$$(SNAPSHOT_ID_$(1))" ] && [ -z "$$(SNAPSHOT_ID)" ]; then \
		SNAPSHOT_ID=$$(jq -r '.builds[-1].artifact_id' $($(1)_MANIFEST)); \
	fi; \
	echo "Setting $(2) to $$SNAPSHOT_ID"; \
	pulumi config set $(2) $$SNAPSHOT_ID
endef

define svc-deploy-prebaked =
$(1)-deploy-prebaked: $(1)-set-snapshot $($(1)_INIT_TARGET)
	$($(1)_PREBAKED_ENV) pulumi refresh --yes
	$($(1)_PREBAKED_ENV) pulumi up --yes

$(1)-deploy-base:
	pulumi config rm $(2) || true
	$($(1)_PREBAKED_ENV) pulumi up --yes
endef

# Per-service config
temporal_CONFIG     = temporal_snapshot_id
temporal_INIT_TARGET = init
temporal_MANIFEST   = $(PACKER_MANIFEST)
temporal_PREBAKED_ENV= TEMPORAL_DOMAIN=$(TEMPORAL_DOMAIN) DUNEBOT_DOMAIN=$(DUNEBOT_DOMAIN)

franky_CONFIG       = franky_snapshot_id
franky_INIT_TARGET  = init
franky_MANIFEST     = $(PACKER_MANIFEST)

wireguard_CONFIG    = wireguard_snapshot_id
wireguard_INIT_TARGET = wireguard-init
wireguard_MANIFEST  = $(PACKER_MANIFEST)
wireguard_PREBAKED_ENV = WIREGUARD_DOMAIN=$(WIREGUARD_DOMAIN)

minipc_CONFIG       = minipc_image_id
minipc_INIT_TARGET  = minipc-init
minipc_MANIFEST     = $(MINIPC_PACKER_MANIFEST)

$(foreach svc,$(SERVICES),$(eval $(call svc-set-snapshot,$(svc),$($(svc)_CONFIG))))
$(foreach svc,$(SERVICES),$(eval $(call svc-deploy-prebaked,$(svc),$($(svc)_CONFIG))))
```

Lines saved: **~40**. Risk: **low** — behaviour identical, you can `make -n` to dry-run each generated target and compare against the originals.

### Strategy 2 — **Single `service-init` / `service-create` / `service-clean` / `service-output` recipes**

#### Before (lines 222–251, 312–334 — ~50 lines)

```make
wireguard-init: build
	pulumi login gs://containifyci-pulumi-state-backend
	#  pulumi login --local
	pulumi stack init $(WIREGUARD_STACK_NAME) || echo ignore if stack $(WIREGUARD_STACK_NAME) already exists
	pulumi stack select -c $(WIREGUARD_STACK_NAME)
	pulumi config set aws:region eu-west-1
	pulumi config set vpn_enabled_ssh ${VPN_ENABLED_SSH}
	pulumi config set ssh_key_file ${PRIVATE_KEY_FILE}
	pulumi config set wireguard_domain $(WIREGUARD_DOMAIN)

wireguard-set-domain:
	@echo "Setting wireguard_domain to $(WIREGUARD_DOMAIN) for stack $(WIREGUARD_STACK_NAME)"
	pulumi config set wireguard_domain $(WIREGUARD_DOMAIN)

wireguard-create: wireguard-init
	WIREGUARD_DOMAIN=$(WIREGUARD_DOMAIN) pulumi up --yes

wireguard-preview: wireguard-init
	WIREGUARD_DOMAIN=$(WIREGUARD_DOMAIN) pulumi preview --diff

wireguard-clean:
	pulumi destroy --yes -s $(WIREGUARD_STACK_NAME)
	pulumi stack rm -f --yes $(WIREGUARD_STACK_NAME) || true

wireguard-output:
	mkdir -p ./output
	pulumi stack output --json > ./output/wireguard-$(WIREGUARD_STACK_NAME).json

wireguard-deploy: wireguard-create wireguard-output
	@echo "Wireguard deployment complete! Stack: $(WIREGUARD_STACK_NAME) Domain: $(WIREGUARD_DOMAIN)"
```

#### After

```make
# Generic helpers (use as prerequisites in generated targets)
pulumi-login-gcs:
	pulumi login gs://containifyci-pulumi-state-backend

pulumi-stack-init = pulumi stack init $(1) || echo ignore if stack $(1) already exists; \
                    pulumi stack select -c $(1); \
                    pulumi config set aws:region eu-west-1; \
                    pulumi config set ssh_key_file ${PRIVATE_KEY_FILE}

define svc-lifecycle =
.PHONY: $(1)-init $(1)-create $(1)-preview $(1)-clean $(1)-output $(1)-deploy
$(1)-init: pulumi-login-gcs
	$(call pulumi-stack-init,$(2))
	$(foreach kv,$(3),pulumi config set $(kv);)
$(1)-create: $(1)-init
	$(4) pulumi up --yes
$(1)-preview: $(1)-init
	$(4) pulumi preview --diff
$(1)-clean:
	pulumi destroy --yes -s $(2)
	pulumi stack rm -f --yes $(2) || true
$(1)-output:
	mkdir -p ./output
	pulumi stack output --json > ./output/$(1).json
$(1)-deploy: $(1)-create $(1)-output
	@echo "$(1) deployment complete! Stack: $(2)"
endef

# Per-service instantiation
$(eval $(call svc-lifecycle,wireguard,$(WIREGUARD_STACK_NAME),vpn_enabled_ssh=${VPN_ENABLED_SSH} wireguard_domain=$(WIREGUARD_DOMAIN),WIREGUARD_DOMAIN=$(WIREGUARD_DOMAIN)))
$(eval $(call svc-lifecycle,minipc,$(MINIPC_STACK_NAME),server_ip=${MINIPC_SERVER_IP} username=${SSH_USER} ssh_port=${MINIPC_SSH_PORT},))
```

Lines saved: **~30**. Risk: **moderate** — `wireguard-clean` originally used `wireguard-clean` (named "clean") while minipc used "destroy". Pick one name; add a `clean-minipc` alias if you want backward compatibility.

### Strategy 3 — **Generic packer targets via `$(eval)`**

#### Before (lines 143–155 and 350–362 — ~25 lines)

```make
packer-init:
	cd $(PACKER_DIR) && packer init .

packer-validate: packer-init
	cd $(PACKER_DIR) && packer validate $(PACKER_VARS_FLAG) .

packer-build: packer-validate
	cd $(PACKER_DIR) && packer build $(PACKER_VARS_FLAG) .
	@echo "Build complete. Snapshot ID:"
	@jq -r '.builds[-1].artifact_id' $(PACKER_MANIFEST)

packer-build-debug: packer-init
	cd $(PACKER_DIR) && PACKER_LOG=1 packer build -debug $(PACKER_VARS_FLAG) .
```

and the minipc version with `MINIPC_PACKER_DIR` / no vars flag.

#### After

```make
# packer-cluster: name | dir | manifest | vars-flag
define packer-cluster =
.PHONY: $(1)-packer-init $(1)-packer-validate $(1)-packer-build $(1)-packer-build-debug
$(1)-packer-init:
	cd $(2) && packer init .
$(1)-packer-validate: $(1)-packer-init
	cd $(2) && packer validate $(3) .
$(1)-packer-build: $(1)-packer-validate
	cd $(2) && packer build $(3) .
	@echo "Build complete. Snapshot ID:"
	@jq -r '.builds[-1].artifact_id' $(4)
$(1)-packer-build-debug: $(1)-packer-init
	cd $(2) && PACKER_LOG=1 packer build -debug $(3) .
endef

# The default packer cluster (temporal/franky/wireguard all share PACKER_DIR)
$(eval $(call packer-cluster,packer,$(PACKER_DIR),$(PACKER_VARS_FLAG),$(PACKER_MANIFEST)))
# Minipc is its own cluster
$(eval $(call packer-cluster,minipc-packer,$(MINIPC_PACKER_DIR),,$(MINIPC_PACKER_MANIFEST)))
```

This lets you also generate the `temporal-packer-build`, `franky-packer-build`, `wireguard-packer-build` aliases by aliasing the default:

```make
temporal-packer-build: packer-build
franky-packer-build: packer-build
wireguard-packer-build: packer-build
```

Lines saved: **~12**. Risk: **low** — purely mechanical, the recipes are 1-for-1.

### Strategy 4 — **Collapse the `service-full-deploy` / `service-recreate-prebaked` one-liners**

#### Before (lines 204–207, 276–279, 366–367)

```make
temporal-full-deploy: packer-build temporal-deploy-prebaked
	@echo "Full deployment complete!"

temporal-recreate-prebaked: clean packer-build temporal-deploy-prebaked

wireguard-full-deploy: packer-build wireguard-deploy-prebaked
	@echo "Full deployment complete!"

wireguard-recreate-prebaked: clean packer-build wireguard-deploy-prebaked

minipc-full-deploy: minipc-packer-build minipc-deploy-prebaked
	@echo "Full mini PC deployment complete!"
```

#### After

```make
define svc-pipeline =
.PHONY: $(1)-full-deploy $(1)-recreate-prebaked
$(1)-full-deploy: $($(1)_PACKER_BUILD) $(1)-deploy-prebaked
	@echo "$(1) full deployment complete!"
$(1)-recreate-prebaked: clean $($(1)_PACKER_BUILD) $(1)-deploy-prebaked
endef

temporal_PACKER_BUILD  = packer-build
wireguard_PACKER_BUILD = packer-build
minipc_PACKER_BUILD    = minipc-packer-build

$(foreach svc,temporal wireguard minipc,$(eval $(call svc-pipeline,$(svc))))
```

Lines saved: **~10**. Risk: **low**.

### Strategy 5 — **One `keys` recipe with a stem parameter**

#### Before (lines 90–91, 342–343)

```make
wireguard-client-keys: prepare
	wg genkey | tee ${TMP_FOLDER}/client_privatekey | wg pubkey > ${TMP_FOLDER}/client_publickey
…
minipc-keys: prepare
	wg genkey | tee ${TMP_FOLDER}/minipc_client_privatekey | wg pubkey > ${TMP_FOLDER}/minipc_client_publickey
```

#### After — use a `pattern rule` (static-pattern):

```make
# stem = service name used in key filename
${TMP_FOLDER}/%_client_publickey: ${TMP_FOLDER}/%_client_privatekey
	wg pubkey < $< > $@

${TMP_FOLDER}/%_client_privatekey: prepare
	wg genkey | tee $@ | wg pubkey > $(@:_privatekey=_publickey)

# Convenience aliases
.PHONY: wireguard-client-keys minipc-keys
wireguard-client-keys: ${TMP_FOLDER}/client_privatekey
minipc-keys:           ${TMP_FOLDER}/minipc_client_privatekey
```

Lines saved: **~2**. Risk: **low** — adds a file timestamp dependency which is actually an *improvement* (re-runs only when key file is missing/corrupt, not every `make`).

### Strategy 6 — **Single `sync-versions` rule with the script path as the stem**

#### Before (lines 299–303)

```make
sync-versions:
	bash packer/hetzner/temporal/scripts/sync-versions.sh

sync-versions-wireguard:
	bash packer/hetzner/wireguard/scripts/sync-versions.sh
```

#### After

```make
# stem = service subdir under packer/hetzner (or packer/local for minipc)
.PHONY: sync-versions sync-versions-wireguard
sync-versions%:          # generic
	@test -n "$*" || { echo "Usage: make sync-versions-<service>"; exit 1; }
	bash packer/hetzner/$*/scripts/sync-versions.sh
```

Or, simpler, just call it via one parameterised target:

```make
SYNC_DIRS := temporal wireguard
.PHONY: $(addprefix sync-versions-,$(SYNC_DIRS))
$(addprefix sync-versions-,$(SYNC_DIRS)): sync-versions-%:
	bash packer/hetzner/$*/scripts/sync-versions.sh
```

Lines saved: **~2**. Risk: **low**.

### Strategy 7 — **Add a complete `.PHONY:` block** (style, not redundancy)

Currently line 1 only declares `build` as `.PHONY`. Every other target in the file is implicitly phony by virtue of always running, but it's a latent bug — if a file named `packer-build` or `wireguard-output` ever appears in the working dir, `make` will silently skip the recipe. Consolidate the `.PHONY` list at the top:

```make
.PHONY: build go-init pulumi-init init verify verify-linux create preview clean \
        recreate deploy local shell browse output prepare …  (full list)
```

Risk: **none** (only safety improvement).

### Strategy 8 — **Drop dead or shadowed code**

- Lines 33, 176, 194, 224, 266, 314: ` pulumi login --local` and `# pulumi destroy` are commented out; harmless but you can collect them in a single commented-out block.
- Line 10 hard-codes a public IP — clearly dead.
- Line 11 references `WIREGUARD_SERVER_PUBLIC_KEY` which is itself never defined (line 94 echoes it but it is empty).
- Line 71: `local: local-cleanup deploy` — `local-cleanup` is **never defined** anywhere; this target is broken. Either define it or delete the target.

Risk: **low** (removing dead code), **moderate** for `local` because fixing it requires knowing intent.

---

## 4. Summary of Estimated Savings

| Strategy                                       | Lines saved | Risk   | Notes                                           |
|------------------------------------------------|-------------|--------|-------------------------------------------------|
| 1. Canned `svc-set-snapshot` / `-deploy-prebaked` / `-deploy-base` | ~40         | Low    | Pure duplication                                |
| 2. Generic `svc-lifecycle` (init/create/preview/clean/output/deploy) | ~30         | Moderate | Touches wireguard-clean vs minipc-destroy naming; ensure all config keys carry over |
| 3. Generic `packer-cluster`                    | ~12         | Low    | Mechanical 1-for-1                              |
| 4. Generic `svc-pipeline`                      | ~10         | Low    |                                                 |
| 5. Static-pattern key files                    | ~2          | Low    | Plus better incremental behaviour               |
| 6. Parameterised `sync-versions`               | ~2          | Low    |                                                 |
| 7. Reorganise + add full `.PHONY` list         | +5 (added)  | None   |                                                 |
| 8. Dead-code removal                            | ~6          | Low    | `local` target needs owner sign-off             |
| **Total**                                      | **~90–100 lines removed** |        | From 382 → ~280–290 lines **after 1–4**, then the `$(eval)` machinery is ~30 lines itself |
| **Net file size**                              | **~200–220 lines** |        | ≈ 45 % reduction                                |

After strategies 1–4 the file would look roughly like this (rough sketch, ~200 lines):

```text
[header / vars / PACKER_VARS_FLAG / domain ENV conditional]    ~60 lines
[Pulumi login / stack init helper]                              ~10 lines
[svc-lifecycle macro]                                           ~20 lines
[packer-cluster macro]                                          ~15 lines
[svc-set-snapshot / svc-deploy-prebaked / svc-deploy-base macro]~25 lines
[svc-pipeline macro]                                            ~10 lines
[Per-service instantiations: 4 services × ~6 lines]            ~25 lines
[Special one-offs: cert, go-init, verify, shell, browse, prepare, sync-versions] ~30 lines
[.PHONY block]                                                  ~5 lines
```

---

## 5. Correctness Gotchas to Watch

1. **The `temporal` / `franky` `deploy-prebaked` targets do `pulumi refresh` and `pulumi up` *without* `--yes`**, while `wireguard-deploy-prebaked` uses `pulumi refresh --yes` (line 267) but no `--yes` on the up, and `minipc-deploy-prebaked` has neither. When you collapse them, decide which one is canonical (recommend *with* `--yes` for all, matching `create` at line 56) and document the change. This is a real behavioural difference.

2. **`temporal-deploy-prebaked` is missing the `init` prerequisite** that the others have. It currently relies on `temporal-set-snapshot` → no `init` chain. After refactoring, make sure the generated `temporal-deploy-prebaked` still does *not* require `init` (or add it — verify by reading `cmd/temporal/hetzner/temporal.go` if you have it handy).

3. **`WIREGUARD_DOMAIN` env-var forward** is present in `wireguard-create`/`-preview`/`-deploy-prebaked`/`-deploy-base` but **absent** in `temporal-`/`franky-`/`minipc-` equivalents. Easy to lose when collapsing. Pass it through a `$(svc)_PREBAKED_ENV` variable (as in Strategy 1 above) to make the asymmetry explicit.

4. **`minipc` uses `minipc_image_id`, the others use `<svc>_snapshot_id`**. Don't blindly normalise to one form — that would break `cmd/minipc/local/minipc.go:78`. Keep the config key as a per-service variable.

5. **`MINIPC_PACKER_BUILD` does not depend on `minipc-init`** like the wireguard one does. Check if that's intentional; if not, add it.

6. **`clean` (line 63) operates on `${STACK_NAME}` (the default `wireguard-hetzner`), while `wireguard-clean` (242) operates on `$(WIREGUARD_STACK_NAME)`** (which under `ENV=test` is `wireguard-test-hetzner`). After refactoring, the generic `clean` target will need to take the service name as a stem — don't accidentally collapse these into the same target.

7. **`minipc-verify` references `${MINIPC_SSH_KEY_FILE}` and `${MINIPC_SSH_USER}` which are never defined** (only `${PRIVATE_KEY_FILE}` and `${SSH_USER}` exist). The target is currently broken. This is a pre-existing bug, not a refactoring risk, but worth flagging.

8. **`local` depends on the undefined `local-cleanup` target** (line 71). Verify with the owner whether to fix or remove.

9. **`pulumi-init` installs AWS plugin v7.35.0 and hcloud plugin v1.39.0 unconditionally** (lines 29–30) but `wireguard-init` and `minipc-init` do not. If you generate a single `pulumi-init` helper, decide which behaviour to keep (or parametrise on `$(PROJECT)`).

10. **`.PHONY` is currently under-declared** — many targets that look like file names (`output`, `build`, `shell`) are not marked phony. The refactor should add a comprehensive `.PHONY:` block; otherwise running `make` from a directory that happens to contain a file named `output` (e.g. `./output/`) will silently no-op.

---

## 6. Recommended Order of Changes

1. **Strategy 8** (dead code) — quick win, low risk, do first to clean the field.
2. **Strategy 7** (`.PHONY` block) — safety net before mass edits.
3. **Strategy 3** (packer cluster) — most mechanical; verify with `make -n`.
4. **Strategy 1** (set-snapshot / deploy-prebaked / deploy-base) — biggest single win, but verify the `pulumi refresh --yes` gotcha (#1 above).
5. **Strategy 4** (full-pipeline targets).
6. **Strategy 2** (service-lifecycle) — largest *and* the riskiest because of gotchas #2, #3, #6. Tackle last and keep the old targets in commented-out form for one release cycle.
7. **Strategies 5 & 6** (keys, sync-versions) — small, safe, do whenever.

Each step should leave the file in a runnable state; run `make -n <target>` and compare against the pre-refactor `make -n <target>` output to catch silent regressions.