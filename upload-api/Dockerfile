# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o upload-api ./cmd/api

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates

# Copy binary from builder
COPY --from=builder /app/upload-api .

# Create storage directory
RUN mkdir -p /app/storage/covers /app/storage/books /app/storage/profiles /app/storage/thumbnails

# Expose port
EXPOSE 8001

# Run the application
CMD ["./upload-api"]
