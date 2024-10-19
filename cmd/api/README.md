# eventrunner-api

## Description

## Build

```shell
docker buildx create --use\n
docker buildx build \\n  --platform linux/amd64,linux/arm64 \\n  -t ghcr.io/carverauto/eventrunner-base:v2 \\n  --push .
docker buildx imagetools inspect ghcr.io/carverauto/eventrunner-base:v2\n
make ko-build VERSION=v0.0.03
```