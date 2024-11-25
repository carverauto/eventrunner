# Define variables
KO_DOCKER_REPO := ghcr.io/carverauto/api-admin
VERSION := v0.0.01

# Default target
all: build

# Build the binary locally
build:
	go build -o api-admin .

# Build and push the container image using ko
ko-build:
	@export KO_DOCKER_REPO=$(KO_DOCKER_REPO) && \
	export VERSION=$(VERSION) && \
	ko build --bare --tags $(VERSION) ./

# Clean up
clean:
	rm -f api-admin

.PHONY: all build ko-build clean