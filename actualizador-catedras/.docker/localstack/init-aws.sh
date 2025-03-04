#!/usr/bin/env bash

set -e

cd `mktemp -d`

curl -L -O https://github.com/user-attachments/files/19073886/planes-parseados.tar.gz
tar -xvf planes-parseados.tar.gz
rm planes-parseados.tar.gz

export AWS_ACCESS_KEY_ID=000000000000
export AWS_SECRET_ACCESS_KEY=000000000000

awslocal s3 mb s3://${AWS_S3_BUCKET}
awslocal s3 cp --recursive . s3://${AWS_S3_BUCKET}/
