#!/usr/bin/env bash

cd /docker-entrypoint-initdb.d/

curl -L -O https://github.com/user-attachments/files/19073885/public-schema.sql.gz
tar -xvf public-schema.sql.gz
