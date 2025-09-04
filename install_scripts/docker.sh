#!/bin/bash

# DevEnv - Docker Installation Script
# Installs Docker Engine

set -e

echo "Installing Docker..."

# Add the official Docker repo
sudo install -m 0755 -d /etc/apt/keyrings
sudo wget -qO /etc/apt/keyrings/docker.asc https://download.docker.com/linux/ubuntu/gpg
sudo chmod a+r /etc/apt/keyrings/docker.asc

echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

sudo apt update

# Install Docker engine and standard plugins
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin docker-ce-rootless-extras

# Give current user privileged Docker access
sudo usermod -aG docker ${USER}

# Restart Docker service
sudo systemctl restart docker
sudo systemctl enable docker

# Verify installation
if command -v docker &> /dev/null; then
    echo "Docker installed successfully"
    docker --version
    docker compose version
else
    echo "Docker installation failed"
    exit 1
fi

echo "Docker installation complete"
echo "Note: Log out and back in for group permissions to take effect"