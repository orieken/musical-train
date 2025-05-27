# Azure DevOps CLI

A Go CLI application that interacts with Azure DevOps. The CLI allows managing work items and pull requests in Azure DevOps.

## Features

- Create work items from a JSON file
- Generate a template JSON file for creating work items
- List work items assigned to a user and display time logged
- List open pull requests across all repositories

## Installation

```bash
# Build the CLI
go build -o azure-devops ./cmd/azure-devops
```

## Usage

### Work Items

#### Create Work Items

Create work items in Azure DevOps based on data provided in a JSON file:

```bash
./azure-devops work-items create --json path/to/workitems.json
```

Options:
- `--json`: Path to the JSON file containing work item definitions (required)

#### Generate Template

Generate a template JSON file that can be used as a starting point for creating work items:

```bash
./azure-devops work-items template
```

This will output a template JSON structure that you can save to a file and modify for your needs.

#### List Assigned Work Items

List all work items assigned to a user and display the work item and time logged:

```bash
./azure-devops work-items assigned --user "John Doe"
```

Options:
- `--user`: Username to filter work items by (required)
- `--json`: Output the results in JSON format

Example output:
```
Found 2 work items:

ID: 123
Title: Fix login bug
Type: Bug
State: Active
Assigned To: John Doe
Time Logged: 4.50 hours
Created Date: 2023-05-15T10:30:00Z

ID: 456
Title: Implement new feature
Type: User Story
State: Active
Assigned To: John Doe
Time Logged: 8.25 hours
Created Date: 2023-05-10T09:15:00Z
```

### Pull Requests

#### List Open Pull Requests

List all open pull requests for all repositories in the organization:

```bash
./azure-devops pull-requests list-open
```

Options:
- `--json`: Output the results in JSON format

## JSON Format for Work Items

The JSON file for creating work items should follow this structure:

```json
[
  {
    "title": "Example Work Item",
    "description": "This is an example work item",
    "type": "Task",
    "priority": 2,
    "tags": ["example", "documentation"]
  }
]
```

## Error Handling

The CLI includes comprehensive error handling for various scenarios:

- Invalid command-line arguments
- Connection failures to Azure DevOps
- Invalid JSON format
- Authentication issues

## Aliases

The CLI supports the following aliases:

- `ado` as an alias for `azure-devops`
