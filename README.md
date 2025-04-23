# Master-Mold CLI

A Go-based CLI tool that utilizes a binary execution model for its subcommands.

## Overview

Master-Mold is a CLI tool that manages subcommands as separate executables, enhancing modularity and extensibility. The project is structured as a monorepo, with the core CLI and subcommand binaries residing in separate directories.

## Project Structure

```
master-mold/
├── cmd/
│   ├── master-mold/       # Main CLI application
│   │   └── main.go
│   └── mm-list-binaries/  # List binaries subcommand
│       └── main.go
├── pkg/
│   ├── binary/            # Binary discovery and execution
│   │   ├── discovery.go
│   │   └── execution.go
│   ├── command/           # Command handling
│   │   ├── handler.go
│   │   ├── list_binaries.go
│   │   ├── registry.go
│   │   └── subcommand.go
│   ├── config/            # Configuration management
│   │   └── config.go
│   └── display/           # Display utilities
│       └── binaries.go
├── test/                  # Integration tests
│   └── integration_test.go
├── config/                # Configuration files
│   └── config.toml
├── go.mod
└── go.sum
```

## Features

- **Binary Execution Model**: Discovers and executes external binaries for subcommands.
- **Subcommand Discovery**: Searches for binaries in the system's PATH and in a dedicated directory (`~/.master-mold`).
- **Configuration Management**: Uses TOML for application configuration with environment variable substitution.
- **Error Handling**: Implements robust error handling with informative error messages.
- **Logging**: Uses structured logging for significant events, errors, and debugging information.
- **Command Registry**: Provides a flexible way to register and execute commands.
- **Unit Tests**: Comprehensive unit tests for all components.
- **Integration Tests**: End-to-end tests to ensure the binary execution model works correctly.
- **Azure DevOps Integration**: Create work items in Azure DevOps using a JSON payload.

## Clean Code Practices

The codebase has been refactored to follow clean code practices:

- **Single Responsibility Principle**: Each function and class has a single responsibility.
- **Extract Till You Drop**: Complex functions have been broken down into smaller, more focused functions.
- **Low Cyclomatic Complexity**: Functions have been kept simple with a cyclomatic complexity between 5-7.
- **Comprehensive Testing**: Unit tests and integration tests ensure code quality and correctness.
- **Clear Naming**: Functions and variables have clear, descriptive names.
- **Consistent Error Handling**: Errors are wrapped with context for better debugging.
- **Modular Design**: The codebase is organized into logical modules with clear boundaries.

## Usage

### Building the CLI

```bash
go build -o master-mold ./cmd/master-mold
go build -o mm-list-binaries ./cmd/mm-list-binaries
```

### Running the CLI

```bash
# List available subcommands
./master-mold list-binaries

# Run a subcommand
./master-mold <subcommand> [options]
```

### Installing Subcommands

Subcommands can be installed by placing executables with the prefix `mm-` or `master-mold-` in:

- Any directory in the system's PATH
- The `~/.master-mold` directory

## Testing

### Running Unit Tests

```bash
go test ./pkg/...
```

### Running Integration Tests

```bash
go test ./test/...
```

## Configuration

The CLI uses a TOML configuration file located at `config/config.toml` or `~/.master-mold/config.toml`. The configuration supports environment variable substitution.

Example configuration:

```toml
# Base directory for master-mold (supports environment variable substitution)
base_dir = "${HOME}/.master-mold"

# Timeout in seconds for command execution
timeout = 10
```

## Azure DevOps Integration

The Azure DevOps subcommand provides functionality to create work items in Azure DevOps.

### Building the Subcommand

```bash
go build -o mm-azure-devops ./cmd/azure-devops
cp mm-azure-devops ~/.master-mold/
```

### Environment Variables

The following environment variables are required for authentication and configuration:

- `AZURE_DEVOPS_PAT`: Your Azure DevOps Personal Access Token
- `AZURE_DEVOPS_ORG`: Your Azure DevOps Organization name
- `AZURE_DEVOPS_PROJECT`: Your Azure DevOps Project name
- `AZURE_DEVOPS_API_VERSION` (optional): The API version to use (defaults to "7.0")

### Usage

#### Generating a Template

```bash
master-mold azure-devops work-items template
# or using the alias
master-mold ado work-items template
```

This will generate a `work-item-template.json` file in the current directory with the following structure:

```json
[
  {
    "op": "add",
    "path": "/fields/System.Title",
    "value": "Example Title: Update this value"
  },
  {
    "op": "add",
    "path": "/fields/System.WorkItemType",
    "value": "Task | Bug | User Story | Feature"
  },
  {
    "op": "add",
    "path": "/fields/System.Description",
    "value": "Example Description: Provide a detailed description here."
  },
  {
    "op": "add",
    "path": "/fields/System.AreaPath",
    "value": "YourProject\\YourArea"
  },
  {
    "op": "add",
    "path": "/fields/System.IterationPath",
    "value": "YourProject\\Iteration 1"
  }
]
```

Edit this file to customize the work item details.

#### Creating Work Items

```bash
master-mold azure-devops work-items create --json work-items.json
# or using the alias
master-mold ado work-items create --json work-items.json
```

This will create work items in Azure DevOps based on the JSON file provided.
