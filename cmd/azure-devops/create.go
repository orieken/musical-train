package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/webapi"
	"github.com/microsoft/azure-devops-go-api/azuredevops/workitemtracking"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// Environment variable names
const (
	EnvAzureDevOpsToken       = "AZURE_DEVOPS_PAT"
	EnvAzureDevOpsOrg         = "AZURE_DEVOPS_ORG"
	EnvAzureDevOpsProject     = "AZURE_DEVOPS_PROJECT"
	EnvAzureDevOpsAPIVersion  = "AZURE_DEVOPS_API_VERSION"
	DefaultAzureDevOpsAPIVersion = "7.0"
)

// createWorkItems creates work items in Azure DevOps based on a JSON file
func createWorkItems(cmd *cobra.Command, args []string) {
	logger.Info("Creating work items")

	// Create a flag provider from the command
	provider := &CommandFlagProvider{cmd: cmd}

	// Get the JSON file path from the flag
	jsonFilePath, err := getJSONFilePath(provider)
	if err != nil {
		handleError("Failed to get JSON file path", err)
		return
	}

	// Process the work items
	err = processWorkItems(jsonFilePath)
	if err != nil {
		handleError("Failed to process work items", err)
		return
	}

	logger.Info("Work items created successfully")
}

// FlagProvider is an interface for getting flag values
type FlagProvider interface {
	GetStringFlag(name string) (string, error)
}

// CommandFlagProvider adapts a cobra.Command to the FlagProvider interface
type CommandFlagProvider struct {
	cmd *cobra.Command
}

// GetStringFlag gets a string flag value from the command
func (p *CommandFlagProvider) GetStringFlag(name string) (string, error) {
	return p.cmd.Flags().GetString(name)
}

// getJSONFilePath gets the JSON file path from the flag provider
func getJSONFilePath(provider FlagProvider) (string, error) {
	jsonFilePath, err := provider.GetStringFlag("json")
	if err != nil {
		return "", errors.Wrap(err, "failed to get json flag")
	}
	return jsonFilePath, nil
}

// handleError logs an error and exits the program
func handleError(message string, err error) {
	logger.Error(message, "error", err)
	fmt.Printf("Error: %s: %v\n", message, err)
	os.Exit(1)
}

// processWorkItems reads work items from a file and creates them in Azure DevOps
func processWorkItems(jsonFilePath string) error {
	// Read the JSON file
	workItemFields, err := readWorkItemsFromFile(jsonFilePath)
	if err != nil {
		return errors.Wrap(err, "failed to read work items from file")
	}

	// Get the Azure DevOps connection details from environment variables
	connectionDetails, err := getAzureDevOpsConnectionDetails()
	if err != nil {
		return err
	}

	// Create the work items
	createdWorkItems, err := createAzureDevOpsWorkItems(connectionDetails, workItemFields)
	if err != nil {
		return errors.Wrap(err, "failed to create work items")
	}

	// Print the created work items
	printCreatedWorkItems(createdWorkItems)

	return nil
}

// printCreatedWorkItems prints information about created work items
func printCreatedWorkItems(workItems []workitemtracking.WorkItem) {
	fmt.Println("Successfully created work items:")
	for _, workItem := range workItems {
		fmt.Printf("  - ID: %d, Title: %s\n", *workItem.Id, getWorkItemTitle(workItem))
	}
}

// ConnectionDetails holds the details needed to connect to Azure DevOps
type ConnectionDetails struct {
	Token       string
	Organization string
	Project     string
	APIVersion  string
}

// getAzureDevOpsConnectionDetails gets the Azure DevOps connection details from environment variables
func getAzureDevOpsConnectionDetails() (*ConnectionDetails, error) {
	// Get the token
	token := os.Getenv(EnvAzureDevOpsToken)
	if token == "" {
		return nil, fmt.Errorf("Azure DevOps Personal Access Token not found. Set the %s environment variable", EnvAzureDevOpsToken)
	}

	// Get the organization
	org := os.Getenv(EnvAzureDevOpsOrg)
	if org == "" {
		return nil, fmt.Errorf("Azure DevOps Organization not found. Set the %s environment variable", EnvAzureDevOpsOrg)
	}

	// Get the project
	project := os.Getenv(EnvAzureDevOpsProject)
	if project == "" {
		return nil, fmt.Errorf("Azure DevOps Project not found. Set the %s environment variable", EnvAzureDevOpsProject)
	}

	// Get the API version (optional)
	apiVersion := os.Getenv(EnvAzureDevOpsAPIVersion)
	if apiVersion == "" {
		apiVersion = DefaultAzureDevOpsAPIVersion
	}

	return &ConnectionDetails{
		Token:       token,
		Organization: org,
		Project:     project,
		APIVersion:  apiVersion,
	}, nil
}

// readWorkItemsFromFile reads work item fields from a JSON file
func readWorkItemsFromFile(filePath string) ([]WorkItemField, error) {
	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read JSON file")
	}

	// Unmarshal the JSON
	var workItemFields []WorkItemField
	if err := json.Unmarshal(data, &workItemFields); err != nil {
		return nil, errors.Wrap(err, "failed to parse JSON")
	}

	return workItemFields, nil
}

// createAzureDevOpsWorkItems creates work items in Azure DevOps
func createAzureDevOpsWorkItems(connectionDetails *ConnectionDetails, fields []WorkItemField) ([]workitemtracking.WorkItem, error) {
	// Create a client for the Work Item Tracking API
	client, err := createAzureDevOpsClient(connectionDetails)
	if err != nil {
		return nil, err
	}

	// Group fields by work item type
	workItemsByType := groupFieldsByWorkItemType(fields)

	// Create the work items
	return createWorkItemsByType(client, connectionDetails.Project, workItemsByType)
}

// createAzureDevOpsClient creates a client for the Azure DevOps API
func createAzureDevOpsClient(connectionDetails *ConnectionDetails) (workitemtracking.Client, error) {
	// Create a connection to Azure DevOps
	connection := azuredevops.NewPatConnection(
		fmt.Sprintf("https://dev.azure.com/%s", connectionDetails.Organization),
		connectionDetails.Token,
	)

	// Create a client for the Work Item Tracking API
	client, err := workitemtracking.NewClient(context.Background(), connection)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Azure DevOps client")
	}

	return client, nil
}

// createWorkItemsByType creates work items for each work item type
func createWorkItemsByType(client workitemtracking.Client, project string, workItemsByType map[string][]WorkItemField) ([]workitemtracking.WorkItem, error) {
	var createdWorkItems []workitemtracking.WorkItem

	for workItemType, workItemFields := range workItemsByType {
		// Create the work item
		workItem, err := createWorkItem(client, project, workItemType, workItemFields)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create work item of type %s", workItemType)
		}

		createdWorkItems = append(createdWorkItems, *workItem)
	}

	return createdWorkItems, nil
}

// groupFieldsByWorkItemType groups fields by work item type
func groupFieldsByWorkItemType(fields []WorkItemField) map[string][]WorkItemField {
	result := make(map[string][]WorkItemField)

	// Find the work item type field
	var workItemType string
	for _, field := range fields {
		if field.Path == "/fields/System.WorkItemType" {
			workItemType = field.Value
			break
		}
	}

	// If no work item type is found, use "Task" as default
	if workItemType == "" {
		workItemType = "Task"
	}

	// Group all fields under this work item type
	result[workItemType] = fields

	return result
}

// createWorkItem creates a work item in Azure DevOps
func createWorkItem(client workitemtracking.Client, project string, workItemType string, fields []WorkItemField) (*workitemtracking.WorkItem, error) {
	// Convert fields to JSON patches
	patches := convertFieldsToPatches(fields)

	// Create the work item
	return createWorkItemWithPatches(client, project, workItemType, patches)
}

// convertFieldsToPatches converts work item fields to JSON patch operations
func convertFieldsToPatches(fields []WorkItemField) []webapi.JsonPatchOperation {
	patches := make([]webapi.JsonPatchOperation, len(fields))

	for i, field := range fields {
		// Convert string to Operation type
		op := webapi.OperationValues.Add
		if field.Op != "add" {
			logger.Warn("Unsupported operation type", "op", field.Op)
		}

		patches[i] = webapi.JsonPatchOperation{
			Op:    &op,
			Path:  &field.Path,
			Value: field.Value,
		}
	}

	return patches
}

// createWorkItemWithPatches creates a work item with the given patches
func createWorkItemWithPatches(client workitemtracking.Client, project string, workItemType string, patches []webapi.JsonPatchOperation) (*workitemtracking.WorkItem, error) {
	// Create the work item
	workItem, err := client.CreateWorkItem(
		context.Background(),
		workitemtracking.CreateWorkItemArgs{
			Project:  &project,
			Type:     &workItemType,
			Document: &patches,
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create work item")
	}

	return workItem, nil
}

// getWorkItemTitle gets the title of a work item
func getWorkItemTitle(workItem workitemtracking.WorkItem) string {
	if workItem.Fields == nil {
		return "Unknown"
	}

	title, ok := (*workItem.Fields)["System.Title"]
	if !ok {
		return "Unknown"
	}

	return fmt.Sprintf("%v", title)
}
