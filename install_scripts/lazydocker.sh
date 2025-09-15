#!/bin/bash

# DevEnv - Lazydocker Installation Script
# Installs lazydocker terminal UI for docker

set -e

echo "Installing lazydocker..."

curl https://raw.githubusercontent.com/jesseduffield/lazydocker/master/scripts/install_update_linux.sh | bash

# Verify installation
if command -v lazydocker &> /dev/null; then
    echo "lazydocker installed successfully"
    lazydocker --version
else
    echo "lazydocker installation failed"
    exit 1
fi

echo "lazydocker installation complete"