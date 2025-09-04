#!/bin/bash

# DevEnv - NVM Installation Script

set -e

echo "Installing NVM..."


# Download and install NVM
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash

# Source nvm to make it available in current session
export NVM_DIR="$HOME/.nvm"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
[ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion"


# Verify installation
if command -v nvm &> /dev/null && command -v node &> /dev/null; then
    echo "NVM and Node.js installed successfully"
    nvm --version
else
    echo "NVM installation failed"
    exit 1
fi

echo "NVM installation completed"