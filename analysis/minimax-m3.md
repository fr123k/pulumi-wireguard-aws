This Makefile is a prime candidate for refactoring. It currently uses a "copy-paste-modify" pattern to implement nearly identical lifecycles for different services (`temporal`, `franky`, `wireguard`, and `minipc`).

### 1. Identified Redundancy Patterns

The Makefile repeats the same logical workflow for every service. The duplication falls into these specific buckets:

*   **The Snapshot Lifecycle:** `temporal-set-snapshot` (L168), `franky-set-snapshot` (L186), `wireguard-set-snapshot` (L258), and `minipc-set-snapshot` (L369) are identical except for the config key and the manifest file path.
*   **Prebaked Deployment:** `*-deploy-prebaked` (L175, L193, L265, L376) always follows the pattern: *Set Snapshot $\rightarrow$ Init $\rightarrow$ Refresh $\rightarrow$ Up*.
*   **Base Deployment:** `*-deploy-base` (L180, L198, L270, L380) always performs: *Remove Snapshot ID $\rightarrow$ Up*.
*   **Stack Initialization:** `pulumi-init` (L28), `wireguard-init` (L222), and `minipc-init` (L312) share 90% of their logic (login, stack init, select, and setting common configs like `ssh_key_file`).
*   **Infrastructure Cleanup:** `clean` (L63), `wireguard-clean` (L242), and `minipc-destroy` (L328) all perform the same `pulumi destroy` and `pulumi stack rm` sequence.
*   **Packer Boilerplate:** The `packer-init/validate/build` triplet (L143-150) is repeated exactly for the `minipc` targets (L350-357).
*   **Output Export:** `output` (L81), `wireguard-output` (L246), and `minipc-output` (L332) are identical except for the filename.

### 2. Categorization of Duplication

| Group | Template Logic | Variance |
| :--- | :--- | :--- |
| **Cloud Services** | `temporal`, `franky`, `wireguard` | Config key (e.g., `temporal_snapshot_id` vs `wireguard_snapshot_id`) |
| **Specialized Hardware** | `minipc` | Different Packer directory, different config keys (`minipc_image_id`) |
| **Pulumi Lifecycle** | `init` $\rightarrow$ `preview` $\rightarrow$ `create` $\rightarrow$ `output` $\rightarrow$ `clean` | Stack name and a few specific `config set` calls |
| **Packer Lifecycle** | `init` $\rightarrow$ `validate` $\rightarrow$ `build` $\rightarrow$ `debug` | Working directory (`PACKER_DIR`) |

### 3. Proposed Reduction Strategies

#### A. Use a `define` function for repetitive Shell logic
Instead of writing the `if [ -z "$(SNAPSHOT_ID)" ]` block four times, create a macro.

#### B. Parameterized Pattern Rules or Target Templates
Since the services follow the same pattern, we can use a variable-based approach. Instead of `temporal-deploy-prebaked` and `franky-deploy-prebaked`, we use a single `deploy-prebaked` target that relies on a `SERVICE` variable.

#### C. Centralized Pulumi Helper Targets
Create a generic `pulumi-manage` target that takes a stack name and action, reducing the `clean` and `output` targets to single-line calls.

---

### 4. Before/After Sketches

#### Example 1: Snapshot Logic
**Before:** (L168, L186, L258, L369)
```makefile
wireguard-set-snapshot:
    @if [ -z "$(SNAPSHOT_ID)" ]; then \
        SNAPSHOT_ID=$$(jq -r '.builds[-1].artifact_id' $(PACKER_MANIFEST)); \
    fi; \
    pulumi config set wireguard_snapshot_id $$SNAPSHOT_ID
```

**After:**
```makefile
define set_snapshot
	@if [ -z "$(SNAPSHOT_ID)" ]; then \
		SNAPSHOT_ID=$$(jq -r '.builds[-1].artifact_id' $(PACKER_MANIFEST)); \
	fi; \
	pulumi config set $(1)_snapshot_id $$SNAPSHOT_ID
endef

# Usage in target:
$(SERVICE)-set-snapshot:
	$(call set_snapshot,$(SERVICE))
```

#### Example 2: Packer Pipeline
**Before:** (L143-150 and L350-357)
Two sets of 4 targets each.

**After:**
```makefile
# Use a variable for the directory
PACKER_DIR_SELECTED ?= $(PACKER_DIR)

packer-init:
	cd $(PACKER_DIR_SELECTED) && packer init .

# To run for minipc: make packer-init PACKER_DIR_SELECTED=$(MINIPC_PACKER_DIR)
```

#### Example 3: Stack Management
**Before:** `wireguard-clean` (L242) and `minipc-destroy` (L328).

**After:**
```makefile
# Generic cleanup target
cleanup-stack:
	pulumi destroy --yes -s $(CURRENT_STACK)
	pulumi stack rm -f --yes $(CURRENT_STACK) || true

wireguard-clean: CURRENT_STACK=$(WIREGUARD_STACK_NAME) ; $(MAKE) cleanup-stack
minipc-destroy: CURRENT_STACK=$(MINIPC_STACK_NAME) ; $(MAKE) cleanup-stack
```

### 5. Impact Estimate

| Change | Est. Lines Removed | Risk Level | Notes |
| :--- | :--- | :--- | :--- |
| **Snapshot Macro** | ~40 lines | Low | Direct replacement of identical logic. |
| **Packer Consolidation** | ~20 lines | Low | Uses variable overrides. |
| **Pulumi Lifecycle Folding** | ~60 lines | Medium | Requires careful handling of `CURRENT_STACK`. |
| **Init/Deploy Templates** | ~50 lines | Medium | Some services have unique config keys. |
| **Total** | **~170+ lines** | **Low/Medium** | Reduces file size by ~45%. |

### 6. Correctness Gotchas

1.  **Variable Shadowing:** Ensure `$(PROJECT)` and `$(VM)` are not overwritten accidentally when switching between `temporal` and `minipc`.
2.  **`.PHONY`:** When moving to generic targets (e.g., `cleanup-stack`), ensure these helpers are marked `.PHONY` to avoid conflicts with files named `cleanup-stack`.
3.  **Conditional Blocks:** The `ifeq ($(ENV),test)` blocks (L120, L212) must remain outside the generic macros as they define the variables the macros use.
4.  **Slight Deviations:** `minipc-set-snapshot` uses `minipc_image_id` instead of `_snapshot_id`. The macro must allow the config key to be passed as an argument rather than assuming a suffix.