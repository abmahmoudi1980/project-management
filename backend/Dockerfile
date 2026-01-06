# Multi-stage build for Go backend
FROM golang:1.21-alpine AS builder

# Install dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server .

# Final stage
FROM alpine:latest

# Install runtime deps (ca-certs for HTTPS, wget for HEALTHCHECK)
RUN apk --no-cache add ca-certificates wget

# Create non-root user
RUN addgroup -S app && adduser -S -G app app

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server ./server

# Create .env file location (will be mounted or provided at runtime)
# Copy migration files if needed for initialization
COPY migration/ ./migration/

# Upload storage (can be mounted as a named volume)
RUN mkdir -p /app/uploads && chown -R app:app /app

USER app

# Expose port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:3000/api/health || exit 1

# Run the application
CMD ["./server"]
