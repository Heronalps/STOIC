#!bin/bash

# One time
NUM_IMAGE=$(echo $1 | bc)
DATA_STRING=$(jq -n --arg ni "$NUM_IMAGE" '{"num_image":$ni}')
echo $DATA_STRING
kubeless function call image-clf-inf --data "$DATA_STRING"