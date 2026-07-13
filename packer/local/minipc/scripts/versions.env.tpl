#!/bin/bash
set -euo pipefail

# versions.env.tpl - Template for version pins
# This file is rendered by Packer with the actual version values.

FRANKY_VERSION=${FRANKY_VERSION}