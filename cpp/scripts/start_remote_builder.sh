#!/usr/bin/env bash

cd "$(dirname "$0")/.." || exit
# build base image
docker build \
  --tag etcd-app-builder-base \
  -f dockerfiles/Dockerfile.builder_base \
  .
# build remote-builder image
docker build \
  --tag etcd-app-remote-builder \
  -f dockerfiles/Dockerfile.builder_remote \
  .
# remove remote-builder container (if exists)
docker container rm -f /etcd-app-remote-builder
# start new remote-builder container
docker run \
    --detach \
    --cap-add sys_ptrace \
    --publish 127.0.0.1:2222:22 \
    --name etcd-app-remote-builder \
    etcd-app-remote-builder
