#!/bin/bash

set -e

echo "🚀 Installing DevOps Memory CLI..."

# Detect OS
OS="$(uname)"

if [[ "$OS" == "Darwin" ]]; then
    URL="https://github.com/Hardik19102003/devops-memory-assistant/releases/latest/download/devops-memory"
elif [[ "$OS" == "Linux" ]]; then
    URL="https://github.com/Hardik19102003/devops-memory-assistant/releases/latest/download/devops-memory-linux"
else
    echo "❌ Unsupported OS"
    exit 1
fi

# Download binary
curl -L $URL -o devops-memory

# Make executable
chmod +x devops-memory

# Move to PATH
sudo mv devops-memory /usr/local/bin/

echo "✅ Installed successfully!"
echo "👉 Run: devops-memory"
