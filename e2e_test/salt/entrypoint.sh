#!/bin/sh
set -e

# SALT_MASTER_CONFIG/SALT_MINION_CONFIG hold a JSON object (valid YAML) with
# config overrides, mirroring the env vars used by the old saltstack/salt image.
if [ -n "$SALT_MASTER_CONFIG" ]; then
  printf '%s\n' "$SALT_MASTER_CONFIG" > /etc/salt/master.d/docker.conf
fi

if [ -n "$SALT_MINION_CONFIG" ]; then
  printf '%s\n' "$SALT_MINION_CONFIG" > /etc/salt/minion.d/docker.conf
fi

exec "$@"
