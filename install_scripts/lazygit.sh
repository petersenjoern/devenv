#!/bin/bash

# DevEnv - Lazygit Installation Script
# Installs lazygit terminal UI for git

set -e

echo "Installing lazygit..."

cd /tmp
LAZYGIT_VERSION=$(curl -s "https://api.github.com/repos/jesseduffield/lazygit/releases/latest" | grep -Po '"tag_name": "v\K[^"]*')
curl -sLo lazygit.tar.gz "https://github.com/jesseduffield/lazygit/releases/latest/download/lazygit_${LAZYGIT_VERSION}_Linux_x86_64.tar.gz"
tar -xf lazygit.tar.gz lazygit
sudo install lazygit /usr/local/bin
rm lazygit.tar.gz lazygit
cd -

# Create config directory
mkdir -p ~/.config/lazygit/
touch ~/.config/lazygit/config.yml

# Verify installation
if command -v lazygit &> /dev/null; then
    echo "lazygit installed successfully"
    lazygit --version
else
    echo "lazygit installation failed"
    exit 1
fi

echo "lazygit installation complete"