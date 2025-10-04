#!/usr/bin/env bash
set -euo pipefail

JAEGER_HOST=${JAEGER_HOST:-localhost:16686}

echo "Visit http://$JAEGER_HOST to explore ArcBridge traces."
echo "Example: curl 'http://$JAEGER_HOST/api/traces?service=controlplane'"
