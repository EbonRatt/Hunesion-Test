#!/bin/bash
# Run Agent Script (Linux/macOS)
# This script runs the ansible-control agent

echo "Starting Ansible Control Agent..."

# Navigate to project root
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
GO_DIR="$PROJECT_ROOT/ansible-control/go"

# Check for executable in bin directory first (recommended), then in go directory
BIN_DIR="$PROJECT_ROOT/ansible-control/bin"
EXE_PATH="$BIN_DIR/ansible-control"
if [ ! -f "$EXE_PATH" ]; then
    # Fallback to go directory
    EXE_PATH="$GO_DIR/ansible-control"
    if [ ! -f "$EXE_PATH" ]; then
        echo "Error: ansible-control executable not found!"
        echo "Please build the application first:"
        echo "  cd ansible-control/go"
        echo "  mkdir ../bin"
        echo "  go build -o ../bin/ansible-control main.go"
        exit 1
    fi
fi

# Check if Docker services are running
echo "Checking Docker services..."
if ! docker ps --filter "name=ansible" --format "{{.Names}}" | grep -q ansible; then
    echo "Warning: Ansible container is not running."
    echo "Starting services..."
    "$SCRIPT_DIR/start-services.sh"
    if [ $? -ne 0 ]; then
        echo "Error: Failed to start services."
        exit 1
    fi
fi

# Change to executable directory
EXE_DIR="$(dirname "$EXE_PATH")"
cd "$EXE_DIR"

echo ""
echo "========================================"
echo "  Ansible Control Agent"
echo "  Main File: main.go"
echo "  Executable: ansible-control"
echo "========================================"
echo ""

# Run the agent
./ansible-control
