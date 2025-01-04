#!/bin/bash
set -e

# Environment setup for preview
export ENVIRONMENT=preview
export GIN_MODE=debug
export RATE_LIMIT_DISABLED=true
export PORT=8090

# Ensure proper permissions
chmod +x backend/main
mkdir -p backend/data
chmod -R 777 backend/data

# Build and deploy preview environment
docker compose -f docker-compose.preview.yml down --remove-orphans
docker compose -f docker-compose.preview.yml build --no-cache
docker compose -f docker-compose.preview.yml up -d

# Wait for services to be healthy
echo "Waiting for services to be healthy..."
sleep 10

# Check service health
docker compose -f docker-compose.preview.yml ps
