# Define variables
KO_DOCKER_REPO := ghcr.io/carverauto/eventrunner-api
VERSION := v0.0.03

# Default target
all: build

# Build the binary locally
build:
	go build -o eventrunner-api .

# Build and push the container image using ko
ko-build:
	@export KO_DOCKER_REPO=$(KO_DOCKER_REPO) && \
	export VERSION=$(VERSION) && \
	ko build --bare --tags $(VERSION) ./

# Clean up
clean:
	rm -f eventrunner-api

.PHONY: all build ko-build clean