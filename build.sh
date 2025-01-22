#!/usr/bin/env bash

REGISTRY='registry.cn-beijing.aliyuncs.com'
NAMESPACE='licsber'
NAME='myfakessh'
# 2023-01-22-1415-CST
TAG=`date +%F-%H%M-CST`

IMAGE="$REGISTRY/$NAMESPACE/$NAME"
echo "$IMAGE:$TAG"
PLATFORM='linux/amd64,linux/arm64'
sudo docker buildx build \
  --platform $PLATFORM \
  -t "$IMAGE:$TAG" \
  -t "$IMAGE:latest" \
  --pull --push .
