#!/bin/bash

set -e

echo "Starting E2E tests..."

# Create results directory
mkdir -p results

# Start services and run tests
docker-compose -f docker-compose.e2e.yml up --build --abort-on-container-exit

# Cleanup
echo "Cleaning up..."
docker-compose -f docker-compose.e2e.yml down --volumes --remove-orphans

echo "E2E tests completed!"