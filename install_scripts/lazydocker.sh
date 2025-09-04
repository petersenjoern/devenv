#!/bin/bash

# DevEnv - Lazydocker Installation Script
# Installs lazydocker terminal UI for docker

set -e

echo "Installing lazydocker..."

cd /tmp
LAZYDOCKER_VERSION=$(curl -s "https://api.github.com/repos/jesseduffield/lazydocker/releases/latest" | grep -Po '"tag_name": "v\K[^"]*')
curl -sLo lazydocker.tar.gz "https://github.com/jesseduffield/lazydocker/releases/latest/download/lazydocker_${LAZYDOCKER_VERSION}_Linux_x86_64.tar.gz"
tar -xf lazydocker.tar.gz lazydocker
sudo install lazydocker /usr/local/bin
rm lazydocker.tar.gz lazydocker
cd -

# Verify installation
if command -v lazydocker &> /dev/null; then
    echo "lazydocker installed successfully"
    lazydocker --version
else
    echo "lazydocker installation failed"
    exit 1
fi

echo "lazydocker installation complete"