#!/bin/bash

# DevEnv - Docker Installation Script
# Installs Docker Engine

set -e

echo "Installing Docker..."

# Add the official Docker repo
# Add Docker's official GPG key:
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc

# Add the repository to Apt sources:
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update

# Install Docker engine and standard plugins
sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

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