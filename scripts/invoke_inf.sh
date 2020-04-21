#!bin/bash

# One time
ZIP_PATH=$(echo $1)
DATA_STRING=$(jq -n --arg ni "$ZIP_PATH" '{"zip_path":$ni}')
echo $DATA_STRING
kubeless function call image-clf-inf --data "$DATA_STRING"