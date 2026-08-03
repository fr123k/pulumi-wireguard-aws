# Makefile Redundancy Reduction Plan

**File:** `Makefile` (364 lines after franky removal, ~66 named targets)
**Repo:** `fr123k/pulumi-wireguard-aws` (local fork at `./pulumi-wireguard-aws`)
**Analysis by:** 4 models from the `ollama-cloud` provider, run in parallel:
- `deepseek-v4-flash`
- `kimi-k2.6`
- `gemma4:31b`
- `minimax-m3`

> NOTE: The `franky` service (Makefile targets `franky-set-snapshot`,
> `franky-deploy-prebaked`, `franky-deploy-base`, and the project at
> `cmd/franky/hetzner/`) was removed after the four-model analysis was
> conducted. The original model outputs in `analysis/` still reference franky
> as one of the four services; the consolidated strategies below should be
> applied to the **three remaining services**: `temporal`, `wireguard`,
> `minipc`. `FRANKY_DOMAIN` in the shared ENV domain block and franky
> references in `cloud-init/` and `packer/hetzner/franky/` were intentionally
> kept.
>
> This document consolidates the four independent analyses into a single
> actionable plan. Per-model outputs are kept verbatim in
> `analysis/<model>.md`.

## Executive Summary

All four models agree: the Makefile is one Pulumi/Packer workflow cloned for
the services (`temporal`, `wireguard`, `minipc` — `franky` has since been
removed) plus a generic default stack. Roughly **45–58% of the file is
boilerplate** — the same recipe bodies repeated with only a service name and
a config key swapped in.

Consensus targets across models:

| Model           | Estimated final size | Reduction | Primary mechanism |
|-----------------|----------------------|-----------|-------------------|
| deepseek-v4-flash | ~335 lines          | ~12%      | `define` + `$(eval)` per target family, incremental |
| kimi-k2.6         | ~160 lines          | ~58%      | Aggressive `$(eval)` macros + alias layer |
| gemma4:31b        | ~200–220 lines      | ~45%      | Service table + canned recipes + dead-code cleanup |
| minimax-m3        | ~210 lines          | ~45%      | `define` macros + generic helper targets + var overrides |

**Recommended target: ~180–210 lines (~45–50% reduction)** using the
macro-driven approach below, which all four models independently arrived at.

---

## 1. Inventory of Repeated Patterns

All four models identified the same duplicated families. Consolidated table:

| # | Pattern Family | Copies | Line ranges | Services | Template |
|---|----------------|--------|-------------|----------|----------|
| 1 | `*-init` (login + stack init/select + config set) | 3 | 28–40, 222–230, 312–320 | pulumi, wireguard, minipc | nearly identical |
| 2 | `*-create` (pulumi up) | 3 | 55–58, 236–237, 322–323 | pulumi, wireguard, minipc | identical (some pass env) |
| 3 | `*-preview` (pulumi preview --diff) | 3 | 60–61, 239–240, 325–326 | pulumi, wireguard, minipc | identical (some pass env) |
| 4 | `*-clean`/`*-destroy` (destroy + rm stack) | 3 | 63–65, 242–244, 328–330 | pulumi, wireguard, minipc | identical (stack name differs) |
| 5 | `*-output` (mkdir + stack output --json) | 3 | 81–83, 246–248, 332–334 | pulumi, wireguard, minipc | identical (filename differs) |
| 6 | `*-set-snapshot` (read manifest + config set) | 4 | 168–173, 186–191, 258–263, 369–374 | temporal, franky, wireguard, minipc | identical (config key + manifest differ) |
| 7 | `*-deploy-prebaked` (set-snapshot + init + refresh + up) | 4 | 175–178, 193–196, 265–268, 376–378 | temporal, franky, wireguard, minipc | same core (env vars + dep order differ) |
| 8 | `*-deploy-base` (config rm snapshot + up) | 4 | 180–182, 198–200, 270–272, 380–382 | temporal, franky, wireguard, minipc | identical (config key + env differ) |
| 9 | Packer init/validate/build/debug triplet | 2 | 143–155, 350–362 | packer, minipc-packer | same (dir + vars flag differ) |
| 10 | `*-full-deploy` (packer-build + deploy-prebaked) | 3 | 204–205, 276–277, 366–367 | temporal, wireguard, minipc | one-liners |
| 11 | `*-recreate-prebaked` (clean + packer-build + deploy-prebaked) | 2 | 207, 279 | temporal, wireguard | one-liners |
| 12 | `*-deploy` composite (init + create + output) | 2 | 69, 250–251 | pulumi, wireguard | one-liners |
| 13 | `sync-versions*` (bash script call) | 2 | 299–300, 302–303 | temporal, wireguard | identical (path differs) |
| 14 | `*-keys` (wg genkey | tee | wg pubkey) | 2 | 90–91, 342–343 | wireguard, minipc | identical (filename differs) |
| 15 | ENV domain `ifeq` blocks | 2 | 120–130, 212–220 | (shared) | partial shadow (WIREGUARD_DOMAIN duplicated) |

**Conservative estimate of duplicated lines: ~180 / 382 ≈ 47%** (deepseek, kimi).

---

## 2. Categorization

### Category A — Pure-template duplication
Recipe body is identical except for a variable substitution. Safe to fold into
a single `define` + `$(eval)`.

| Template | Varying elements | Copies |
|----------|------------------|--------|
| `*-create` | optional env prefix | 3 |
| `*-preview` | optional env prefix | 3 |
| `*-clean`/`*-destroy` | stack name | 3 |
| `*-output` | output filename | 3 |
| `*-set-snapshot` | manifest path, config key | 4 |
| `*-deploy-base` | config key, optional env | 4 |
| `sync-versions*` | script path | 2 |
| `*-keys` | key filename | 2 |

### Category B — Near-identical with one or two divergent lines
Share ~80%+ of the body. Foldable but require per-service variable knobs.

| Template | Differences | Copies |
|----------|-------------|--------|
| `*-init` | different config keys; `pulumi-init` adds plugin install; `wireguard-init` adds `wireguard_domain` | 3 |
| `*-deploy-prebaked` | temporal passes 2 env vars; wireguard passes 1; franky/minipc pass none; dep order varies; `--yes` on refresh varies | 4 |
| Packer triplet | minipc uses `MINIPC_PACKER_DIR` (no vars flag); `minipc-packer-build` has `2>/dev/null \|\| echo` fallback; `minipc-packer-build-debug` depends on validate not init | 2 |
| `*-full-deploy` | minipc uses `minipc-packer-build` | 3 |
| `*-deploy` composite | `deploy` deps `init create output`; `wireguard-deploy` deps `wireguard-create wireguard-output` | 2 |
| ENV domain blocks | second block shadows `WIREGUARD_DOMAIN` from first | 2 |

### Category C — Unique one-off targets (do not template)
`go-init`, `build`, `verify`, `verify-linux`, `recreate`, `local`, `shell`,
`browse`, `prepare`, `wireguard-client-keys`, `wireguard-public-key`,
`validate-wireguard`, `validate-jenkins`, `packer-cleanup`, `packer-list`,
`wireguard-set-domain`, `wireguard-deploy-test`, `cert-generate-wildcard`,
`cert-check-expiry`, `minipc-verify`, `minipc-shell`.

---

## 3. Reduction Strategies (Consolidated)

The four models converged on the same core technique: **`define` macros +
`$(eval)` instantiation + a per-service variable table**. kimi and gemma4
propose a single consolidated service table; deepseek proposes per-family
templates; minimax proposes generic helper targets with variable overrides.
The plan below combines the most robust elements of each.

### Strategy 1 — Per-service variable table

Declare each service once with all its knobs, then generate targets from it.

```makefile
# name | config_key        | stack_name           | manifest               | init_target | deploy_env
temporal_CONFIG     = temporal_snapshot_id
temporal_STACK      =
temporal_MANIFEST   = $(PACKER_MANIFEST)
temporal_INIT       = init
temporal_DEPLOY_ENV = TEMPORAL_DOMAIN=$(TEMPORAL_DOMAIN) DUNEBOT_DOMAIN=$(DUNEBOT_DOMAIN)

franky_CONFIG       = franky_snapshot_id
franky_STACK        =
franky_MANIFEST     = $(PACKER_MANIFEST)
franky_INIT         = init
franky_DEPLOY_ENV   =

wireguard_CONFIG     = wireguard_snapshot_id
wireguard_STACK      = $(WIREGUARD_STACK_NAME)
wireguard_MANIFEST   = $(PACKER_MANIFEST)
wireguard_INIT       = wireguard-init
wireguard_DEPLOY_ENV = WIREGUARD_DOMAIN=$(WIREGUARD_DOMAIN)

minipc_CONFIG       = minipc_image_id
minipc_STACK        = $(MINIPC_STACK_NAME)
minipc_MANIFEST     = $(MINIPC_PACKER_MANIFEST)
minipc_INIT         = minipc-init
minipc_DEPLOY_ENV   =

SERVICES := temporal franky wireguard minipc
```

### Strategy 2 — `*-set-snapshot` / `*-deploy-prebaked` / `*-deploy-base` macro

**Before:** ~32 lines (4 copies × 3 targets). **After:** ~12 lines.

```makefile
define SNAPSHOT_RULES
$(1)-set-snapshot:
	@if [ -z "$$(SNAPSHOT_ID)" ]; then \
		SNAPSHOT_ID=$$(jq -r '.builds[-1].artifact_id' $$($(1)_MANIFEST)); \
	fi; \
	echo "Setting $$($(1)_CONFIG) to $$$$SNAPSHOT_ID"; \
	pulumi config set $$($(1)_CONFIG) $$$$SNAPSHOT_ID

$(1)-deploy-prebaked: $(1)-set-snapshot $$($(1)_INIT)
	$$($(1)_DEPLOY_ENV) pulumi refresh
	$$($(1)_DEPLOY_ENV) pulumi up --yes

$(1)-deploy-base:
	pulumi config rm $$($(1)_CONFIG) || true
	$$($(1)_DEPLOY_ENV) pulumi up --yes
endef

$(foreach svc,$(SERVICES),$(eval $(call SNAPSHOT_RULES,$(svc))))
```

> ⚠️ **Gotcha (all models):** `wireguard-deploy-prebaked` currently uses
> `pulumi refresh --yes`; the others use `pulumi refresh` (no `--yes`). Decide
> canonical behavior (recommend `--yes` everywhere) and document. Also
> dependency order differs (`wireguard` does init-before-set-snapshot; others
> do set-snapshot-before-init). The macro above uses set-snapshot-first; if
> order matters, add an `$(svc)_DEPLOY_DEPS` knob.

### Strategy 3 — `*-init` macro

**Before:** ~30 lines (3 copies). **After:** ~15 lines.

```makefile
define INIT_TEMPLATE
$(1)-init: build
	pulumi login gs://containifyci-pulumi-state-backend
	pulumi stack init $$($(1)_STACK_NAME) || echo ignore if stack $$($(1)_STACK_NAME) already exists
	pulumi stack select -c $$($(1)_STACK_NAME)
	$(foreach kv,$($(1)_INIT_CONFIG),pulumi config set $(kv);)
endef

pulumi_INIT_STACK   = $(STACK_NAME)
pulumi_INIT_CONFIG  = aws:region=$(AWS_REGION) vpn_enabled_ssh=$(VPN_ENABLED_SSH) ssh_key_file=$(PRIVATE_KEY_FILE)
wireguard_INIT_STACK = $(WIREGUARD_STACK_NAME)
wireguard_INIT_CONFIG = aws:region=$(AWS_REGION) vpn_enabled_ssh=$(VPN_ENABLED_SSH) ssh_key_file=$(PRIVATE_KEY_FILE) wireguard_domain=$(WIREGUARD_DOMAIN)
minipc_INIT_STACK   = $(MINIPC_STACK_NAME)
minipc_INIT_CONFIG  = server_ip="$(MINIPC_SERVER_IP)" ssh_key_file="$(PRIVATE_KEY_FILE)" username="$(SSH_USER)" ssh_port=$(MINIPC_SSH_PORT)
```

> ⚠️ **Gotcha (gemma, deepseek):** `pulumi-init` runs `pulumi plugin install
> resource aws 7.35.0`, `pulumi plugin install resource hcloud 1.39.0`, and
> `pulumi plugin ls` — the others do not. Folding it into the template adds
> those steps to wireguard/minipc (idempotent but a behavior change). Either
> keep a `PULUMI_PLUGIN_INSTALL` knob or leave `pulumi-init` as a one-off.

### Strategy 4 — `*-create` / `*-preview` / `*-clean` / `*-output` macro

**Before:** ~24 lines (3 copies × 4 targets). **After:** ~18 lines.

```makefile
define LIFECYCLE_TEMPLATE
$(1)-create: $(1)-init
	$$($(1)_DEPLOY_ENV) pulumi up --yes

$(1)-preview: $(1)-init
	$$($(1)_DEPLOY_ENV) pulumi preview --diff

$(1)-clean:
	pulumi destroy --yes -s $$($(1)_STACK_NAME)
	pulumi stack rm -f --yes $$($(1)_STACK_NAME) || true

$(1)-output:
	mkdir -p ./output
	pulumi stack output --json > $$($(1)_OUTPUT_FILE)
endef

pulumi_OUTPUT_FILE = ./output/wireguard-ec2.json
wireguard_OUTPUT_FILE = ./output/wireguard-$(WIREGUARD_STACK_NAME).json
minipc_OUTPUT_FILE = ./output/minipc.json
```

> ⚠️ **Gotcha (kimi, deepseek, gemma):** minipc uses `minipc-destroy`, not
> `minipc-clean`. The macro produces `minipc-clean`. Add a backward-compatible
> alias: `minipc-destroy: minipc-clean`.
>
> ⚠️ `clean` (line 63) operates on `${STACK_NAME}` (default
> `wireguard-hetzner`), while `wireguard-clean` operates on
> `$(WIREGUARD_STACK_NAME)` which under `ENV=test` is `wireguard-test-hetzner`.
> Don't collapse these into the same target.

### Strategy 5 — Packer triplet macro

**Before:** ~20 lines (2 sets). **After:** ~14 lines.

```makefile
define PACKER_TEMPLATE
$(1)-packer-init:
	cd $$($(1)_PACKER_DIR) && packer init .

$(1)-packer-validate: $(1)-packer-init
	cd $$($(1)_PACKER_DIR) && packer validate $$($(1)_PACKER_VARS) .

$(1)-packer-build: $(1)-packer-validate
	cd $$($(1)_PACKER_DIR) && packer build $$($(1)_PACKER_VARS) .
	@echo "Build complete. Snapshot ID:"
	@jq -r '.builds[-1].artifact_id' $$($(1)_PACKER_MANIFEST) 2>/dev/null || echo "No manifest found"

$(1)-packer-build-debug: $(1)-packer-init
	cd $$($(1)_PACKER_DIR) && PACKER_LOG=1 packer build -debug $$($(1)_PACKER_VARS) .
endef

cloud_PACKER_DIR = $(PACKER_DIR)
cloud_PACKER_MANIFEST = $(PACKER_MANIFEST)
cloud_PACKER_VARS = $(PACKER_VARS_FLAG)
minipc_packer_PACKER_DIR = $(MINIPC_PACKER_DIR)
minipc_packer_PACKER_MANIFEST = $(MINIPC_PACKER_MANIFEST)
minipc_packer_PACKER_VARS =

# Keep backward-compatible bare names:
packer-init := cloud-packer-init
packer-validate := cloud-packer-validate
packer-build := cloud-packer-build
packer-build-debug := cloud-packer-build-debug
```

> ⚠️ **Gotcha (deepseek, kimi, gemma):** `minipc-packer-build-debug` depends
> on `minipc-packer-validate`, but the main `packer-build-debug` depends on
> `packer-init`. The macro above uses `$(1)-packer-init` for both — align them
> deliberately or parametrize the dependency.

### Strategy 6 — `*-full-deploy` / `*-recreate-prebaked` macro

**Before:** ~5 lines. **After:** ~6 lines (neutral, but eliminates repetition).

```makefile
define FULL_DEPLOY_TEMPLATE
$(1)-full-deploy: $(2) $(1)-deploy-prebaked
	@echo "Full $(1) deployment complete!"

$(1)-recreate-prebaked: clean $(2) $(1)-deploy-prebaked
endef

$(eval $(call FULL_DEPLOY_TEMPLATE,temporal,packer-build))
$(eval $(call FULL_DEPLOY_TEMPLATE,wireguard,packer-build))
$(eval $(call FULL_DEPLOY_TEMPLATE,minipc,minipc-packer-build))
```

### Strategy 7 — Merge the two ENV domain `ifeq` blocks

**Before:** 22 lines across lines 120–130 and 212–220, with `WIREGUARD_DOMAIN`
set twice. **After:** 11 lines.

```makefile
ifeq ($(ENV),test)
  TEMPORAL_DOMAIN    ?= temporal-test.dunebot.io
  DUNEBOT_DOMAIN     ?= githubapp-test.dunebot.io
  FRANKY_DOMAIN      ?= franky-test.dunebot.io
  WIREGUARD_DOMAIN   ?= wg-test.fr123k.uk
  WIREGUARD_STACK_NAME ?= wireguard-test-hetzner
  PRIVATE_KEY_FILE   = ./keys/wireguard-test
else
  TEMPORAL_DOMAIN    ?= temporal.dunebot.io
  DUNEBOT_DOMAIN     ?= githubapp.dunebot.io
  FRANKY_DOMAIN      ?= franky.dunebot.io
  WIREGUARD_DOMAIN   ?= wg.fr123k.uk
  WIREGUARD_STACK_NAME ?= wireguard-hetzner
  PRIVATE_KEY_FILE   = ./keys/wireguard
endif
```

> ⚠️ Watch `=` vs `?=` semantics for `PRIVATE_KEY_FILE`.

### Strategy 8 — Unify `sync-versions`

**Before:** 4 lines. **After:** 2 lines.

```makefile
SYNC_SERVICE ?= temporal
sync-versions:
	bash packer/hetzner/$(SYNC_SERVICE)/scripts/sync-versions.sh
sync-versions-wireguard: ; $(MAKE) sync-versions SYNC_SERVICE=wireguard
```

### Strategy 9 — Unify `*-keys`

**Before:** 4 lines. **After:** 2 lines + pattern rule.

```makefile
${TMP_FOLDER}/%_client_publickey: ${TMP_FOLDER}/%_client_privatekey
	wg pubkey < $< > $@

${TMP_FOLDER}/%_client_privatekey: prepare
	wg genkey | tee $@ | wg pubkey > $(@:_privatekey=_publickey)

wireguard-client-keys: ${TMP_FOLDER}/client_privatekey
minipc-keys:           ${TMP_FOLDER}/minipc_client_privatekey
```

### Strategy 10 — Add a comprehensive `.PHONY` block

Currently only `.PHONY: build` (line 1). Every model flagged this as a latent
bug. Add all targets; with `$(eval)`-generated targets you can compute the
list:

```makefile
.PHONY: build go-init init verify verify-linux create preview clean recreate \
        deploy local shell browse output prepare wireguard-client-keys \
        wireguard-public-key validate-wireguard validate-jenkins \
        packer-init packer-validate packer-build packer-build-debug \
        packer-cleanup packer-list cert-generate-wildcard cert-check-expiry \
        $(foreach s,$(SERVICES),$(foreach t,init create preview clean output set-snapshot deploy-prebaked deploy-base full-deploy recreate-prebaked,$s-$t))
```

### Strategy 11 — Drop dead / shadowed code

- Line 10: hard-coded public IP (dead).
- Line 11: `WIREGUARD_SERVER_PUBLIC_KEY` never defined (empty when used at line 94).
- Lines 33, 176, 194, 224, 266, 314: commented-out `pulumi login --local` /
  `# pulumi destroy` — collect or remove.
- Line 71: `local: local-cleanup deploy` — `local-cleanup` is **never defined**;
  this target is broken. Fix or remove (needs owner sign-off).
- Line 9: `WIREGUARD_SERVER_IP=$(shell pulumi stack output publicIp)` runs at
  parse time (kimi flagged). Use `=` (deferred) or move into the target.

---

## 4. Estimated Savings

| Strategy | Lines before | Lines after | Saved | Risk |
|----------|--------------|-------------|-------|------|
| 1. Service table | — | ~12 | +12 (infra) | Low |
| 2. Snapshot/deploy-prebaked/deploy-base macro | ~32 | ~12 | **~20** | Medium |
| 3. `*-init` macro | ~30 | ~15 | **~15** | Low–Med (plugin install) |
| 4. Lifecycle macro | ~24 | ~18 | **~6** | Low (naming alias) |
| 5. Packer macro | ~20 | ~14 | **~6** | Medium (debug dep) |
| 6. Full-deploy macro | ~5 | ~6 | ~0 | Low |
| 7. Merge ENV blocks | 22 | 11 | **~11** | Low (`=` vs `?=`) |
| 8. Unify sync-versions | 4 | 2 | **~2** | Low |
| 9. Unify keys | 4 | 3 | **~1** | Low |
| 10. `.PHONY` block | 1 | 5 | −4 | None (fix) |
| 11. Dead code | ~6 | 0 | **~6** | Low (local needs owner) |
| **Total** | **~382** | **~180–210** | **~170–200** | |

**Net reduction: ~45–52%** (target ~180–210 lines), matching the kimi/gemma
estimates. deepseek's conservative estimate (~335) reflects an incremental
approach that keeps more targets explicit; the aggressive macro approach
reaches the lower number.

---

## 5. Correctness Gotchas (Consensus Across Models)

1. **`pulumi-init` installs plugins** (lines 29–31) but `wireguard-init` and
   `minipc-init` do not. Folding into one template adds those steps
   everywhere — idempotent but a behavior change. Parametrize or keep
   `pulumi-init` separate. *(deepseek, gemma, kimi)*

2. **`--yes` inconsistency on `pulumi refresh`**: `wireguard-deploy-prebaked`
   (line 267) uses `pulumi refresh --yes`; temporal/franky/minipc use plain
   `pulumi refresh`. Pick one. *(gemma)*

3. **Dependency order in `*-deploy-prebaked`**:
   - temporal/franky/minipc: `set-snapshot` before `init`
   - wireguard: `init` before `set-snapshot`
   If `init` creates the stack and `set-snapshot` configures it, set-snapshot
   first may fail on a fresh stack. The wireguard order is arguably correct.
   The macro should use a `$(svc)_DEPLOY_DEPS` knob or default to init-first.
   *(deepseek)*

4. **`minipc` uses `minipc_image_id`, others use `<svc>_snapshot_id`**. Don't
   normalize — it would break `cmd/minipc/local/minipc.go`. Keep config key as
   a per-service variable. *(all models)*

5. **`minipc-packer-build` has `2>/dev/null || echo "No manifest found"`** but
   `packer-build` does not. Preserve via conditional or absorb the fallback
   into the shared template (recommend the latter). *(deepseek, kimi)*

6. **`minipc-packer-build-debug` depends on `validate`**; main
   `packer-build-debug` depends on `init`. Align deliberately. *(deepseek, kimi)*

7. **`minipc-destroy` vs `minipc-clean` naming**: template generates
   `minipc-clean`, breaking backward compat. Add an alias. *(kimi, deepseek)*

8. **`clean` (line 63) uses `${STACK_NAME}`** (default `wireguard-hetzner`),
   while `wireguard-clean` uses `$(WIREGUARD_STACK_NAME)` (which under
   `ENV=test` is `wireguard-test-hetzner`). Don't collapse these into one
   target. *(gemma)*

9. **`minipc-verify` (line 337) and `minipc-shell` (line 340)** reference
   `MINIPC_SSH_KEY_FILE` and `MINIPC_SSH_USER` which are never defined (only
   `PRIVATE_KEY_FILE` and `SSH_USER` exist). Pre-existing bugs — flag, don't
   introduce. *(deepseek, gemma)*

10. **`local` (line 71) depends on undefined `local-cleanup`**. Broken target.
    Owner must decide fix vs remove. *(gemma)*

11. **`WIREGUARD_SERVER_IP` is `:=`-evaluated at parse time** (line 9) via
    `$(shell ...)`, which fails if no stack is selected. Use deferred `=` or
    move into the target. *(kimi)*

12. **`.PHONY` under-declared** — many phony targets (`output`, `build`,
    `shell`) are not marked, so a stray file named `output` would silently
    no-op the target. *(all models)*

13. **ENV `ifeq` blocks must remain outside macros** — they define the
    variables the macros consume. Don't fold them. *(minimax)*

14. **`$$$$` escaping in `$(eval)` recipes** — the snapshot macro needs
    `$$$$SNAPSHOT_ID` (four dollars) to produce `$$SNAPSHOT_ID` in the shell
    recipe. Test with `make -n`. *(kimi)*

---

## 6. Recommended Implementation Order

Consensus across models (kimi's phased plan, endorsed by deepseek/gemma):

### Phase 1 — Low-risk cleanup (~30 lines, no behavior change)
1. **Strategy 11** — drop dead code (line 71 needs owner sign-off).
2. **Strategy 10** — add full `.PHONY` block (safety net before mass edits).
3. **Strategy 7** — merge the two ENV domain `ifeq` blocks.
4. **Strategy 8** — unify `sync-versions`.

### Phase 2 — Mechanical macros (~50 lines, low risk)
5. **Strategy 1** — add the per-service variable table.
6. **Strategy 5** — Packer triplet macro (verify with `make -n`).
7. **Strategy 2** — snapshot/deploy-prebaked/deploy-base macro (watch gotchas
   #2, #3, #5).

### Phase 3 — Lifecycle macros (~60 lines, medium risk)
8. **Strategy 3** — `*-init` macro (watch gotcha #1).
9. **Strategy 4** — create/preview/clean/output macro (add `minipc-destroy`
   alias — gotcha #7).
10. **Strategy 6** — full-deploy/recreate macro.
11. **Strategy 9** — unify `*-keys`.

### Validation gate
After each step, run `make -n <target>` and diff against the pre-refactor
output to catch silent regressions (gemma's recommendation). Leave the file
runnable at every commit.

---

## 7. Final Recommended Structure (Sketch)

```text
[Header: exports, vars, PACKER_VARS_FLAG]                       ~25 lines
[Merged ENV domain ifeq block]                                  ~11 lines
[Service table: 4 services × ~6 lines]                          ~25 lines
[Pulumi login / stack-init helper]                              ~5  lines
[SNAPSHOT_RULES macro + eval]                                   ~15 lines
[INIT_TEMPLATE macro + eval]                                    ~18 lines
[LIFECYCLE_TEMPLATE macro + eval]                               ~20 lines
[PACKER_TEMPLATE macro + eval]                                  ~16 lines
[FULL_DEPLOY_TEMPLATE macro + eval]                             ~8  lines
[One-offs: go-init, build, verify, shell, browse, prepare,       ~25 lines
           keys, cert, sync-versions, minipc-verify/shell]
[.PHONY block]                                                  ~6  lines
─────────────────────────────────────────────────────────────────
Total                                                           ~174 lines
```

Adding a fifth service (e.g. `grafana`) becomes one block in the service
table + one `$(eval)` line, instead of another 15-line copy-paste block.

---

## Appendix — Per-Model Outputs

Full unedited analyses from each model are in:
- `analysis/deepseek-v4-flash.md`
- `analysis/kimi-k2.6.md`
- `analysis/gemma4-31b.md`
- `analysis/minimax-m3.md`

### Where the models agreed
- All four identified the same 13–15 duplicated target families.
- All four recommended `define` + `$(eval)` macros as the primary mechanism.
- All four flagged `.PHONY` under-declaration, `minipc_image_id` vs
  `<svc>_snapshot_id`, and the ENV `ifeq` block duplication.
- All four gave a final estimate in the 45–58% reduction range.

### Where they differed
- **Aggressiveness:** kimi (~160 lines) and gemma (~200) push hard macros;
  deepseek (~335) is conservative and incremental; minimax (~210) uses generic
  helper targets + variable overrides rather than full `$(eval)`.
- **Mechanism:** minimax prefers `CURRENT_STACK`-override helper targets
  (`cleanup-stack`) over per-service `$(eval)` generation — simpler but less
  DRY. kimi/gemma/deepseek prefer generated targets.
- **`local` target:** only gemma explicitly flagged `local-cleanup` as
  undefined; kimi noted commented-out code generally.
- **`WIREGUARD_SERVER_IP` parse-time eval:** only kimi flagged it as a latent
  bug.