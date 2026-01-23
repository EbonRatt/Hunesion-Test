#!/bin/bash
# Build Script for Linux/macOS
# Builds the executable to ../bin directory to keep source clean

echo "Building Ansible Control Agent..."

# Create bin directory if it doesn't exist
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
BIN_DIR="$(dirname "$SCRIPT_DIR")/bin"
mkdir -p "$BIN_DIR"

# Build
EXE_PATH="$BIN_DIR/ansible-control"
echo "Building to: $EXE_PATH"
go build -o "$EXE_PATH" main.go

if [ $? -eq 0 ]; then
    echo "Build successful!"
    echo "Executable: $EXE_PATH"
else
    echo "Build failed!"
    exit 1
fi
