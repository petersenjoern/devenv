#!/bin/bash

# DevEnv - Zsh Installation Script
# Installs Zsh shell with Oh My Zsh

set -e

echo "Installing Zsh..."

# Install Zsh
sudo apt update
sudo apt install -y zsh

# Install Oh My Zsh (non-interactive)
RUNZSH=no CHSH=no sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"

# Install plugins
ZSH_CUSTOM=${ZSH_CUSTOM:-~/.oh-my-zsh/custom}

# zsh-autosuggestions
if [ ! -d "$ZSH_CUSTOM/plugins/zsh-autosuggestions" ]; then
    git clone https://github.com/zsh-users/zsh-autosuggestions "$ZSH_CUSTOM/plugins/zsh-autosuggestions"
fi

# zsh-syntax-highlighting
if [ ! -d "$ZSH_CUSTOM/plugins/zsh-syntax-highlighting" ]; then
    git clone https://github.com/zsh-users/zsh-syntax-highlighting.git "$ZSH_CUSTOM/plugins/zsh-syntax-highlighting"
fi

# zsh-z (z command)
if [ ! -d "$ZSH_CUSTOM/plugins/zsh-z" ]; then
    git clone https://github.com/agkozak/zsh-z "$ZSH_CUSTOM/plugins/zsh-z"
fi

# powerlevel10k theme
if [ ! -d "$ZSH_CUSTOM/themes/powerlevel10k" ]; then
    git clone --depth=1 https://github.com/romkatv/powerlevel10k.git "$ZSH_CUSTOM/themes/powerlevel10k"
fi

# Verify installation
if command -v zsh &> /dev/null; then
    echo "Zsh installed successfully"
    zsh --version
else
    echo "Zsh installation failed"
    exit 1
fi

echo "Zsh installation complete with Oh My Zsh and plugins"
echo "Note: Use 'chsh -s $(which zsh)' to set zsh as default shell"