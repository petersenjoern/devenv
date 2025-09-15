#!/bin/bash

# DevEnv - Visual Studio Code Installation Script
# Installs VS Code with extensions

set -e

echo "Installing Visual Studio Code..."

# Add Microsoft GPG key and repository
cd /tmp
wget -qO- https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor > packages.microsoft.gpg
sudo install -D -o root -g root -m 644 packages.microsoft.gpg /etc/apt/keyrings/packages.microsoft.gpg
echo "deb [arch=amd64,arm64,armhf signed-by=/etc/apt/keyrings/packages.microsoft.gpg] https://packages.microsoft.com/repos/code stable main" | sudo tee /etc/apt/sources.list.d/vscode.list > /dev/null
rm -f packages.microsoft.gpg
cd -

# Update and install
sudo apt update
sudo apt install -y code

# Create config directory
mkdir -p ~/.config/Code/User

# Verify installation
if command -v code &> /dev/null; then
    echo "VS Code installed successfully"
    code --version | head -1
else
    echo "VS Code installation failed"
    exit 1
fi

echo "VS Code installation complete"