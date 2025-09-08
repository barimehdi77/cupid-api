# Build stage
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o ./bin/api ./cmd/api/
RUN go build -o ./bin/fetch ./cmd/fetch/

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates postgresql-client

# Create non-root user
RUN adduser -D -s /bin/sh appuser

# Set working directory
WORKDIR /app

# Copy binaries from builder stage
COPY --from=builder /app/bin/api ./api
COPY --from=builder /app/bin/fetch ./fetch

# Copy migration files
COPY --from=builder /app/cmd/migrate/migrations ./migrations

# Copy environment file template
COPY --from=builder /app/integration.env.example ./.env

# Change ownership to appuser
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/health || exit 1

# Default command (can be overridden)
CMD ["./api"]
