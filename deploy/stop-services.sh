#!/bin/bash
# Stop Docker Services Script (Linux/macOS)

echo "Stopping Docker services..."

# Navigate to project root
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
cd "$PROJECT_ROOT"

# Stop services
docker compose down

if [ $? -eq 0 ]; then
    echo "Services stopped successfully!"
else
    echo "Error: Failed to stop services."
    exit 1
fi
