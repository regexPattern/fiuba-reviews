#!/usr/bin/env bash

set -e

cd `mktemp -d`

curl -L -O https://github.com/user-attachments/files/19076297/planes-parseados.tar.gz
tar -xvf planes-parseados.tar.gz
rm planes-parseados.tar.gz

export AWS_ACCESS_KEY_ID=000000000000
export AWS_SECRET_ACCESS_KEY=000000000000

awslocal s3 mb s3://${AWS_S3_BUCKET}

objs=(
  ingenieria-civil-1C-2025.json
  ingenieria-electronica-1C-2025.json
  ingenieria-en-informatica-1C-2025.json
  ingenieria-en-informatica-2C-2024.json
  ingenieria-en-petroleo-1C-2025.json
  ingenieria-industrial-1C-2025.json
  ingenieria-mecanica-1C-2025.json
  ingenieria-quimica-2C-2024.json
)

for o in ${objs[@]}; do
  if [[ $o =~ ^(.*)-([0-9])C-([0-9]{4})\.json$ ]]; then
    carrera=${BASH_REMATCH[1]}
    cuatri=${BASH_REMATCH[2]}
    anio=${BASH_REMATCH[3]}

    awslocal s3 cp $o s3://${AWS_S3_BUCKET}/$o \
      --metadata "carrera=$carrera,cuatri-numero=$cuatri,cuatri-anio=$anio"
  fi
done
