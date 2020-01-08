#!/bin/bash
# $1 - App: image-clf-inf
# $2 - Python version: 3.6 / 3.7
# $3 - Patched GPU number: 0-8

kubeless function delete $1
sleep 3
kubeless function deploy $1 --runtime gpupython$2 \
                            --from-file ./apps/$1.py \
                            --handler $1.handler \
                            --timeout 10800
sleep 3

kubectl patch deployment $1 --patch "$(cat ./scripts/patch.yaml | 
yq w - spec.template.spec.containers[0].resources.requests[nvidia.com/gpu] $3 | 
yq w - spec.template.spec.containers[0].resources.limits[nvidia.com/gpu] $3 | 
yq w - spec.template.spec.containers[0].name $1 )"