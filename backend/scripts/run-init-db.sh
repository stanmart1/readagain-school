#!/bin/bash

# Build and run the production database initialization script

echo "Building database initialization script..."
cd "$(dirname "$0")"

# Build the Go binary
go build -o init-production-db init-production-db.go

if [ $? -eq 0 ]; then
    echo "✓ Build successful"
    echo "Running database initialization..."
    ./init-production-db
else
    echo "✗ Build failed"
    exit 1
fi
