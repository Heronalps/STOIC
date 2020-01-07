#!/bin/bash
kubeless function delete $1
sleep 3
kubeless function deploy $1 --runtime gpupython$2 \
                                         --from-file ./apps/$1.py \
                                         --handler $1.handler \
                                         --timeout 10800
sleep 3
kubectl patch deployment $1 --patch "$(cat ./scripts/patch-$1.yaml)"
                                                                            