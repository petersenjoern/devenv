# DevEnv Application Specification

## 1. Overview

DevEnv is a Go-based command-line application that automates developer
environment setup for personal use. It supports WSL (Ubuntu) and native Linux
(Ubuntu) environments, providing an interactive TUI for selecting and
installing development tools with minimal default configurations.

### 1.1 Core Objectives

- Automate installation of personal developer tools (zsh, fzf, vscode, vim, neovim, etc.)
- Provide minimal/default configurations that come with tool installation
- Support both WSL and native Ubuntu environments
- Maintain consistent configuration path interfaces
- Operate statelessly without persistent installation tracking

## 2. System Requirements

### 2.1 Supported Environments

- Ubuntu (native Linux)
- Windows Subsystem for Linux (WSL) with Ubuntu

### 2.2 Prerequisites

- Go runtime environment for execution
- Internet connectivity for package downloads
- Sudo privileges for system package installations

## 3. Architecture

### 3.1 Technology Stack

- **Language**: Go
- **TUI Library**: 'huh' for interactive selection interfaces, Bubble Tea to extend functionality for TUI if 'huh' isnt enough
- **Configuration Format**: Embedded YAML with consistent schema
- **Package Management**: System package managers (apt), custom scripts, manual instructions

### 3.2 Project Structure

```
devenv/
├── cmd/
│   ├── install.go       # Interactive installation command
│   └── status.go        # Status reporting command
├── internal/
│   ├── config/          # Configuration parsing and management
│   ├── installer/       # Installation logic for different methods
│   ├── detector/        # Binary and config detection
│   └── tui/             # Terminal UI components
├── templates/           # Minimal configuration templates for the tools and app to install
│   ├── zsh.conf
│   ├── vscode.json
│   ├── mise.toml
│   └── ...
├── install_scripts/     # Complex installation scripts for the tools and app to install
│   ├── mise.sh
│   └── fzf.sh
├── config.yaml          # Embedded tool definitions
└── main.go
```

## Development Commands

### Setup

- `mise install-deps` - Install all dependencies
- `mise trust` - Trust the mise configuration file (required first time)

### Go Development

- `go mod tidy` - Clean up dependencies and download what's needed
- `go mod verify` - verify dependencies
- `go test` - Run all tests in current package
- `go test -v` - Run tests with verbose output
- `go test ./...` - Run all tests in current directory and subdirectories
- `go test ./path/to/package` - Run tests in specific package
- `go test -run TestAdd` - Run specific test function
- `go test -run "Test.\*Add"` - Run tests matching a pattern
