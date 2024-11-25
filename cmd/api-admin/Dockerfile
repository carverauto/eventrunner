# Start from a minimal Alpine image to set up our directory structure
FROM --platform=linux/amd64 alpine:3.14 AS builder

# Create necessary directories
WORKDIR /app
RUN mkdir -p /app/configs /app/migrations

# Copy migrations and configs
COPY migrations/ /app/migrations/
COPY configs/.env /app/configs/.env
COPY configs/.staging.env /app/configs/.staging.env

# Now, create our final image based on distroless
FROM --platform=linux/amd64 cgr.dev/chainguard/go:latest
# debug using a container with a shell
#FROM --platform=linux/amd64 alpine:3.14

# Copy the directory structure from the builder stage
COPY --from=builder /app /app

# Set working directory
WORKDIR /app

# Set environment variables
ENV APP_ENV=staging
ENV GOFR_MIGRATIONS_DIR=/app/migrations

# The ko-built binary will be added as "app" in the container
# We don't set an ENTRYPOINT here as ko will do that for us
ENTRYPOINT ["sh"]
