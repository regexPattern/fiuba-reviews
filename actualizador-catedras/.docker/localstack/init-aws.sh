#!/usr/bin/env bash

set -e

cd `mktemp -d`

curl -L -O https://github.com/user-attachments/files/19150131/planes-parseados.tar.gz
tar -xvf planes-parseados.tar.gz
rm planes-parseados.tar.gz

export AWS_ACCESS_KEY_ID=000000000000
export AWS_SECRET_ACCESS_KEY=000000000000

awslocal s3 mb s3://${AWS_S3_BUCKET}

objs=(
  "ingenieria-civil.json;'ingenieria civil';1;2025"
  "ingenieria-electronica.json;'ingenieria electronica';1;2025"
  "ingenieria-en-energia-electrica.json;'ingenieria en energia electrica';1;2025"
  "ingenieria-en-informatica.json;'ingenieria en informatica';1;2025"
  "ingenieria-en-petroleo.json;'ingenieria en petroleo';1;2025"
  "ingenieria-industrial.json;'ingenieria industrial';1;2025"
  "ingenieria-mecanica.json;'ingenieria mecanica';1;2025"
  "ingenieria-quimica.json;'ingenieria quimica';2;2024"
)

for obj in "${objs[@]}"; do
  archivo=$(echo "$obj" | awk -F';' '{print $1}')
  carrera=$(echo "$obj" | awk -F';' '{print $2}')
  carrera=$(echo "$carrera" | sed "s/'//g")
  cuatri=$(echo "$obj" | awk -F';' '{print $3}')
  anio=$(echo "$obj" | awk -F';' '{print $4}')

  awslocal s3 cp $archivo s3://${AWS_S3_BUCKET}/$archivo \
    --metadata "carrera=$carrera,cuatri-numero=$cuatri,cuatri-anio=$anio"
done
