#!/bin/bash
# Start Docker Services Script (Linux/macOS)
# This script starts Kafka and Ansible containers

echo "Starting Docker services..."

# Navigate to project root
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
cd "$PROJECT_ROOT"

# Check if Docker is running
if ! docker ps > /dev/null 2>&1; then
    echo "Error: Docker is not running. Please start Docker."
    exit 1
fi

# Start services
echo "Starting Ansible container..."
docker compose up -d

if [ $? -eq 0 ]; then
    echo "Services started successfully!"
    echo ""
    echo "Running containers:"
    docker compose ps
else
    echo "Error: Failed to start services."
    exit 1
fi
