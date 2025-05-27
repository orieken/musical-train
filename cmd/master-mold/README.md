# Master-Mold CLI

The main entry point for the Master-Mold CLI tool. This CLI serves as a command dispatcher that discovers and executes subcommands as separate binaries.

## Features

- Binary execution model for subcommands
- Configuration management with TOML
- Structured logging
- Command registry for registering and executing commands
- Error handling with informative error messages

## Installation

```bash
# Build the CLI
go build -o master-mold ./cmd/master-mold
```

## Configuration

The CLI looks for configuration files in the following locations:
- `./config`
- `$HOME/.master-mold`

### Configuration File Example

```toml
# Master-Mold CLI Configuration

# Base directory for master-mold (supports environment variable substitution)
base_dir = "${HOME}/.master-mold"

# Timeout in seconds for command execution
timeout = 10

# Binary discovery paths
[binary]
paths = ["${HOME}/.master-mold/bin", "/usr/local/bin"]
```

## Usage

### Basic Usage

```bash
./master-mold <command> [options]
```

If no command is specified, the CLI will suggest running `master-mold list-binaries` to see available commands.

### List Available Commands

To see all available commands:

```bash
./master-mold list-binaries
```

### Execute a Subcommand

To execute a subcommand:

```bash
./master-mold <subcommand> [options]
```

For example:

```bash
./master-mold k8s-pods --namespace=kube-system
```

## Architecture

The Master-Mold CLI follows a binary execution model:

1. The main CLI (`master-mold`) serves as a command dispatcher
2. Subcommands are implemented as separate binaries
3. The CLI discovers these binaries in the system's PATH and in dedicated directories
4. When a command is executed, the CLI finds the corresponding binary and executes it

This architecture provides several benefits:
- Modularity: Each subcommand is a separate binary
- Extensibility: New commands can be added without modifying the main CLI
- Independence: Subcommands can be developed and deployed independently

## Error Handling

The CLI includes comprehensive error handling for various scenarios:

- Missing commands
- Configuration errors
- Command execution failures