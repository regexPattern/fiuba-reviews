#!/usr/bin/env bash

set -e

curl -L \
  -o /docker-entrypoint-initdb.d/public-schema.sql.gz \
  https://github.com/user-attachments/files/19075454/public-schema.sql.gz

exec /usr/local/bin/docker-entrypoint.sh "$@"
