#!/bin/bash

# DevEnv - Btop Installation Script
# Installs btop system monitor

set -e

echo "Installing btop..."

# Install btop
sudo apt update
sudo apt install -y btop

# Create config directory
mkdir -p ~/.config/btop/themes

# Verify installation
if command -v btop &> /dev/null; then
    echo "btop installed successfully"
    btop --version
else
    echo "btop installation failed"
    exit 1
fi

echo "btop installation complete"