#!/bin/bash

# DevEnv - FZF Installation Script
# Installs fzf fuzzy finder

set -e

echo "Installing fzf..."

# Clone fzf repository
git clone --depth 1 https://github.com/junegunn/fzf.git ~/.fzf

# Install fzf
~/.fzf/install --all

# Verify installation
if command -v fzf &> /dev/null; then
    echo "fzf installed successfully"
    fzf --version
else
    echo "fzf installation failed"
    exit 1
fi

echo "fzf installation complete"