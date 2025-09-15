#!/bin/bash

# DevEnv - Eza Installation Script
# Installs Eza

set -e

echo "Installing Eza..."


sudo mkdir -p /etc/apt/keyrings
wget -qO- https://raw.githubusercontent.com/eza-community/eza/main/deb.asc | sudo gpg --dearmor -o /etc/apt/keyrings/gierens.gpg
echo "deb [signed-by=/etc/apt/keyrings/gierens.gpg] http://deb.gierens.de stable main" | sudo tee /etc/apt/sources.list.d/gierens.list
sudo chmod 644 /etc/apt/keyrings/gierens.gpg /etc/apt/sources.list.d/gierens.list
sudo apt update
sudo apt install -y eza

# Verify installation
if command -v eza &> /dev/null; then
    echo "Eza installed successfully"
    eza --version
else
    echo "Eza installation failed"
    exit 1
fi

echo "Eza installation complete"