#!/bin/bash

# DevEnv - Neovim Installation Script
# Installs Neovim editor with LazyVim

set -e

echo "Installing Neovim..."

# Install Neovim from official releases
cd /tmp
wget -O nvim.tar.gz "https://github.com/neovim/neovim/releases/download/stable/nvim-linux-x86_64.tar.gz"
tar -xf nvim.tar.gz
sudo install nvim-linux-x86_64/bin/nvim /usr/local/bin/nvim
sudo cp -R nvim-linux-x86_64/lib /usr/local/
sudo cp -R nvim-linux-x86_64/share /usr/local/
rm -rf nvim-linux-x86_64 nvim.tar.gz
cd -

# Install supporting tools
sudo apt install -y luarocks

# Create nvim config directory if it doesn't exist
mkdir -p ~/.config/nvim

# Only set up LazyVim if config doesn't exist
if [ ! -f ~/.config/nvim/init.lua ]; then
    echo "Setting up LazyVim configuration..."
    # Clone LazyVim starter
    git clone https://github.com/LazyVim/starter ~/.config/nvim
    # Remove .git to allow user to add to their own repo
    rm -rf ~/.config/nvim/.git
    echo "LazyVim configuration installed"
fi

# Verify installation
if command -v nvim &> /dev/null; then
    echo "Neovim installed successfully"
    nvim --version | head -1
else
    echo "Neovim installation failed"
    exit 1
fi

echo "Neovim installation complete"