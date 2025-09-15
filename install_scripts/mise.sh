#!/bin/bash

# DevEnv - Mise Installation Script
# Installs mise (formerly rtx) runtime manager

set -e

echo "Installing mise..."

# Download and install mise
curl https://mise.run | sh

# Add to PATH in current session
export PATH="$HOME/.local/bin:$PATH"

# Verify installation
if command -v mise &> /dev/null; then
    echo "mise installed successfully"
    mise --version
else
    echo "mise installation failed"
    exit 1
fi

echo "mise installation complete"