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

# Install goose for database migrations
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates postgresql-client make

# Copy goose binary from builder stage
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Create non-root user
RUN adduser -D -s /bin/sh appuser

# Set working directory
WORKDIR /app

# Copy binaries from builder stage
COPY --from=builder /app/bin/api ./api
COPY --from=builder /app/bin/fetch ./fetch

# Copy migration files maintaining directory structure
COPY --from=builder /app/cmd/migrate ./cmd/migrate

# Copy makefile
COPY --from=builder /app/makefile ./makefile

# Copy startup script and make it executable
COPY --from=builder /app/startup.sh ./startup.sh
RUN chmod +x ./startup.sh

# Change ownership to appuser
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/health || exit 1

# Run the startup script
CMD ["./startup.sh"]
