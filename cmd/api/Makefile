# Define variables
KO_DOCKER_REPO := ghcr.io/carverauto/eventrunner
VERSION := v0.0.1

# Default target
all: build

# Build target
build:
        @export KO_DOCKER_REPO=$(KO_DOCKER_REPO) && export VERSION=$(VERSION) && ko build --bare --tags $(VERSION) ./

.PHONY: all build
