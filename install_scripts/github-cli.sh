#!/bin/bash

# DevEnv - GitHub CLI Installation Script
# Installs GitHub CLI (gh)

set -e

echo "Installing GitHub CLI..."

# Add GitHub CLI repository
curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | sudo dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg
sudo chmod go+r /usr/share/keyrings/githubcli-archive-keyring.gpg

echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | sudo tee /etc/apt/sources.list.d/github-cli.list

# Update package list and install
sudo apt update
sudo apt install -y gh

# Verify installation
if command -v gh &> /dev/null; then
    echo "GitHub CLI installed successfully"
    gh --version
else
    echo "GitHub CLI installation failed"
    exit 1
fi

echo "GitHub CLI installation complete"
echo "Use 'gh auth login' to authenticate"