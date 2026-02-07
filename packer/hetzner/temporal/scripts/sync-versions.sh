#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
GOMOD="${SCRIPT_DIR}/../go.mod"
VARIABLES="${SCRIPT_DIR}/../variables.pkr.hcl"

if [[ ! -f "$GOMOD" ]]; then
  echo "ERROR: go.mod not found at $GOMOD" >&2
  exit 1
fi

if [[ ! -f "$VARIABLES" ]]; then
  echo "ERROR: variables.pkr.hcl not found at $VARIABLES" >&2
  exit 1
fi

# Extract version from go.mod for a given module path, stripping the v prefix.
get_version() {
  local module="$1"
  local version
  version=$(grep -E "^\s+${module} v" "$GOMOD" | awk '{print $2}' | sed 's/^v//')
  if [[ -z "$version" ]]; then
    echo "WARNING: module $module not found in go.mod" >&2
    return 1
  fi
  echo "$version"
}

# Update a variable's default value in variables.pkr.hcl.
update_variable() {
  local var_name="$1"
  local new_version="$2"

  if ! grep -q "variable \"${var_name}\"" "$VARIABLES"; then
    echo "WARNING: variable ${var_name} not found in variables.pkr.hcl" >&2
    return 1
  fi

  sed -i.bak -E "/variable \"${var_name}\"/,/^\}/ s/(default[[:space:]]*=[[:space:]]*\").*(\")/\1${new_version}\2/" "$VARIABLES"
  echo "  ${var_name} = ${new_version}"
}

# Module-to-variable mapping using parallel arrays (bash 3 compatible)
MODULES=(
  "github.com/temporalio/cli"
  "go.temporal.io/server"
  "github.com/temporalio/ui-server/v2"
  "github.com/containifyci/secret-operator"
  "github.com/containifyci/temporal-worker"
  "github.com/containifyci/oauth2-storage"
  "github.com/containifyci/dunebot"
)

VARIABLES_NAMES=(
  "temporal_cli_version"
  "temporal_server_version"
  "temporal_ui_version"
  "secret_operator_version"
  "temporal_worker_version"
  "oauth2_storage_version"
  "dunebot_version"
)

echo "Syncing versions from go.mod to variables.pkr.hcl..."

errors=0
for i in "${!MODULES[@]}"; do
  module="${MODULES[$i]}"
  var_name="${VARIABLES_NAMES[$i]}"
  if version=$(get_version "$module"); then
    update_variable "$var_name" "$version" || ((errors++))
  else
    ((errors++))
  fi
done

# Clean up sed backup file
rm -f "${VARIABLES}.bak"

if [[ $errors -gt 0 ]]; then
  echo "Completed with $errors warning(s)."
  exit 1
fi

echo "All versions synced successfully."
