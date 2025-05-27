# Master-Mold List Binaries

A Go CLI application that lists all available Master-Mold binaries. This command helps users discover what subcommands are available in the Master-Mold CLI.

## Features

- Discovers all Master-Mold binaries in the system
- Displays the paths of the discovered binaries
- Indicates whether it's running as a subcommand or standalone command
- Structured logging for debugging and information

## Installation

```bash
# Build the CLI
go build -o mm-list-binaries ./cmd/mm-list-binaries
```

## Usage

### As a Master-Mold Subcommand

```bash
./master-mold list-binaries
```

### As a Standalone Command

```bash
./mm-list-binaries
```

## How It Works

The command performs the following steps:

1. Initializes a logger for structured logging
2. Searches for Master-Mold binaries in:
   - The system's PATH
   - The `~/.master-mold` directory
3. Displays the paths of all discovered binaries
4. Indicates whether it's running as a subcommand of Master-Mold or as a standalone command

## Example Output

```
Available Master-Mold binaries:
/usr/local/bin/k8s-pods
/home/user/.master-mold/bin/azure-devops
/home/user/.master-mold/bin/activemq-cli

Running as a subcommand of master-mold
```

## Integration with Master-Mold

This command is designed to work both as a standalone command and as a subcommand of the Master-Mold CLI. When run as a subcommand, it helps users discover what other subcommands are available.

The Master-Mold CLI suggests running this command when no command is specified:

```
Usage: master-mold <command> [options]
Run 'master-mold list-binaries' to see available commands
```

## Error Handling

The command includes error handling for various scenarios:

- Failure to get the home directory
- Failure to find binaries
- Failure to determine if running as a subcommand