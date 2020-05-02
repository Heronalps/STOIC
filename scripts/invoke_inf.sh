#!bin/bash

# One time
S3_KEY=$(echo $1)
DATA_STRING=$(jq -n --arg ni "$S3_KEY" '{"key_name":$ni}')
echo $DATA_STRING
kubeless function call image-clf-inf --data "$DATA_STRING"