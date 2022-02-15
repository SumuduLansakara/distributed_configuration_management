# How to build

## Using remote-builder container

1. Build base builder image.

```shell
docker build \
  --tag etcd-app-builder-base \
  -f dockerfiles/Dockerfile.builder_base \
  .
```

2. Build remote-builder image.

```shell
docker build \
  --tag etcd-app-remote-builder \
  -f dockerfiles/Dockerfile.builder_remote \
  .
```

3. Run remote-builder container

```shell
docker run \
    --detach \
    --cap-add sys_ptrace \
    --publish 127.0.0.1:2222:22 \
    --name etcd-app-remote-builder \
    etcd-app-remote-builder
```

4. Stop and remove remote-builder container

```shell
docker container rm -f /etcd-app-remote-builder
```
