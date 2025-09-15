FROM ubuntu:22.04

# Prevent interactive prompts during package installation
ENV DEBIAN_FRONTEND=noninteractive
ENV TZ=UTC

# Install basic dependencies for testing
COPY certificates/zscaler.crt /usr/local/share/ca-certificates/zscaler.crt
RUN apt-get update && apt-get install -y \
    sudo \
    ca-certificates \
    build-essential \
    gnupg \
    lsb-release 
RUN update-ca-certificates

# Create a test user with sudo privileges
RUN useradd -m -s /bin/bash testuser && \
    echo "testuser ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers

# Set up working directory
WORKDIR /home/testuser/devenv
USER testuser


# Copy the built binary and test files
COPY --chown=testuser:testuser bin/devenv /usr/local/bin/devenv
COPY --chown=testuser:testuser install_scripts/ ./install_scripts/
RUN find ./install_scripts -type f -name "*.sh" -exec chmod +x {} +
# COPY --chown=testuser:testuser test/integration/testdata/ ./testdata/

# Install Go for testing
RUN curl -fsSL https://go.dev/dl/go1.21.5.linux-amd64.tar.gz | sudo tar -C /usr/local -xzf -
ENV PATH="/usr/local/go/bin:${PATH}"

CMD ["/bin/bash"]