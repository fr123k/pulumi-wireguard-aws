#!/bin/bash
set -euxo pipefail

# render.sh
# Renders {{ VAR }} placeholders in a cloud-init template with real values.
# Every {{ VAR }} in the template is replaced with the value of the matching
# environment variable VAR. If VAR is unset or empty, the placeholder is
# replaced with an empty string so the cloud-init script's own defaults
# (e.g. `[ -z "$VAR" ] && VAR=...`) take effect.
#
# This is generic: any placeholder of the form {{ NAME }} is substituted from
# the environment, so new templates (wireguard, temporal, minipc, ...) need no
# changes here. Placeholders whose variable is not set in the environment are
# left as empty strings, not left as literal "{{ NAME }}".
#
# Usage:
#   DOMAIN=wg.test.local MINIPC_USER=frank ./render.sh \
#     /tmp/cloud-init-template.txt /tmp/cloud-init-rendered.sh

TEMPLATE="$1"
OUTPUT="$2"

if [ ! -f "$TEMPLATE" ]; then
  echo "ERROR: template not found: $TEMPLATE" >&2
  exit 1
fi

cp "$TEMPLATE" "$OUTPUT"

# Collect every distinct {{ NAME }} placeholder from the template.
# Handles optional inner whitespace: {{ NAME }}, {{NAME}}.
PLACEHOLDERS=$(grep -oE '\{\{[[:space:]]*[A-Z0-9_]+[[:space:]]*\}\}' "$OUTPUT" \
  | sed -E 's/\{\{[[:space:]]*([A-Z0-9_]+)[[:space:]]*\}\}/\1/' \
  | sort -u)

for VAR in $PLACEHOLDERS; do
  # Read the env var value (empty if unset). envsubst-style substitution via sed.
  VALUE="${!VAR:-}"
  # Escape sed replacement metacharacters (&, /, \n) in the value.
  ESC_VALUE=$(printf '%s' "$VALUE" | sed -e 's/[\\&/]/\\&/g')
  sed -i -E "s/\{\{[[:space:]]*${VAR}[[:space:]]*\}\}/${ESC_VALUE}/g" "$OUTPUT"
done

chmod +x "$OUTPUT"
echo "=== Rendered cloud-init script to ${OUTPUT} ==="

# Show what was rendered (first 5 lines)
head -5 "$OUTPUT"