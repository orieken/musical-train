package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/workitemtracking"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// AssignedWorkItem represents a work item assigned to a user
type AssignedWorkItem struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Type        string    `json:"type"`
	State       string    `json:"state"`
	AssignedTo  string    `json:"assignedTo"`
	TimeLogged  float64   `json:"timeLogged"`
	CreatedDate time.Time `json:"createdDate"`
}

// listAssignedWorkItems lists all work items assigned to a user
func listAssignedWorkItems(cmd *cobra.Command, args []string) {
	logger.Info("Listing work items assigned to user")

	// Get the username from the flag
	username, err := cmd.Flags().GetString("user")
	if err != nil {
		handleError("Failed to get user flag", err)
		return
	}

	// Check if JSON output is requested
	jsonOutput, err := cmd.Flags().GetBool("json")
	if err != nil {
		handleError("Failed to get json flag", err)
		return
	}

	// Get the work items
	workItems, err := getAssignedWorkItems(username)
	if err != nil {
		handleError("Failed to get assigned work items", err)
		return
	}

	// Print the work items
	if jsonOutput {
		printWorkItemsAsJSON(workItems)
	} else {
		printWorkItemsAsText(workItems)
	}

	logger.Info("Work items listed successfully")
}

// getAssignedWorkItems gets all work items assigned to a user
func getAssignedWorkItems(username string) ([]AssignedWorkItem, error) {
	// Get the Azure DevOps connection details from environment variables
	connectionDetails, err := getAzureDevOpsConnectionDetails()
	if err != nil {
		return nil, err
	}

	// Create a connection to Azure DevOps
	connection := azuredevops.NewPatConnection(
		fmt.Sprintf("https://dev.azure.com/%s", connectionDetails.Organization),
		connectionDetails.Token,
	)

	// Create a client for the Work Item Tracking API
	client, err := workitemtracking.NewClient(context.Background(), connection)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Work Item Tracking client")
	}

	// Build the WIQL query to find work items assigned to the user
	wiql := fmt.Sprintf("SELECT [System.Id], [System.Title], [System.WorkItemType], [System.State], [System.AssignedTo], [Microsoft.VSTS.Scheduling.CompletedWork] FROM WorkItems WHERE [System.AssignedTo] = '%s' ORDER BY [System.ChangedDate] DESC", username)

	// Execute the WIQL query
	wiqlArgs := workitemtracking.QueryByWiqlArgs{
		Wiql: &workitemtracking.Wiql{
			Query: &wiql,
		},
		Project: &connectionDetails.Project,
	}

	queryResult, err := client.QueryByWiql(context.Background(), wiqlArgs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute WIQL query")
	}

	if queryResult.WorkItems == nil || len(*queryResult.WorkItems) == 0 {
		return []AssignedWorkItem{}, nil
	}

	// Get the work item IDs
	var workItemIDs []int
	for _, workItemRef := range *queryResult.WorkItems {
		workItemIDs = append(workItemIDs, *workItemRef.Id)
	}


	// Get each work item individually
	var result []AssignedWorkItem
	for _, workItemID := range workItemIDs {
		// Get the work item
		workItem, err := client.GetWorkItem(
			context.Background(),
			workitemtracking.GetWorkItemArgs{
				Id:      &workItemID,
				Project: &connectionDetails.Project,
			},
		)
		if err != nil {
			logger.Warn("Failed to get work item", "id", workItemID, "error", err)
			continue
		}

		if workItem.Fields == nil {
			continue
		}

		fields := *workItem.Fields
		assignedTo := getFieldValue(fields, "System.AssignedTo.displayName", "Unknown")
		timeLogged := getFieldValueFloat(fields, "Microsoft.VSTS.Scheduling.CompletedWork", 0.0)
		createdDate := getFieldValueTime(fields, "System.CreatedDate", time.Time{})

		result = append(result, AssignedWorkItem{
			ID:          *workItem.Id,
			Title:       getFieldValue(fields, "System.Title", "Unknown"),
			Type:        getFieldValue(fields, "System.WorkItemType", "Unknown"),
			State:       getFieldValue(fields, "System.State", "Unknown"),
			AssignedTo:  assignedTo,
			TimeLogged:  timeLogged,
			CreatedDate: createdDate,
		})
	}

	return result, nil
}

// getFieldValue gets a string field value from work item fields
func getFieldValue(fields map[string]interface{}, fieldName string, defaultValue string) string {
	if value, ok := fields[fieldName]; ok {
		if strValue, ok := value.(string); ok {
			return strValue
		}
	}
	return defaultValue
}

// getFieldValueFloat gets a float field value from work item fields
func getFieldValueFloat(fields map[string]interface{}, fieldName string, defaultValue float64) float64 {
	if value, ok := fields[fieldName]; ok {
		switch v := value.(type) {
		case float64:
			return v
		case float32:
			return float64(v)
		case int:
			return float64(v)
		}
	}
	return defaultValue
}

// getFieldValueTime gets a time field value from work item fields
func getFieldValueTime(fields map[string]interface{}, fieldName string, defaultValue time.Time) time.Time {
	if value, ok := fields[fieldName]; ok {
		if strValue, ok := value.(string); ok {
			if t, err := time.Parse(time.RFC3339, strValue); err == nil {
				return t
			}
		}
	}
	return defaultValue
}

// printWorkItemsAsText prints work items in a human-readable format
func printWorkItemsAsText(workItems []AssignedWorkItem) {
	if len(workItems) == 0 {
		fmt.Println("No work items found.")
		return
	}

	fmt.Printf("Found %d work items:\n\n", len(workItems))

	for _, item := range workItems {
		fmt.Printf("ID: %d\n", item.ID)
		fmt.Printf("Title: %s\n", item.Title)
		fmt.Printf("Type: %s\n", item.Type)
		fmt.Printf("State: %s\n", item.State)
		fmt.Printf("Assigned To: %s\n", item.AssignedTo)
		fmt.Printf("Time Logged: %.2f hours\n", item.TimeLogged)
		fmt.Printf("Created Date: %s\n", item.CreatedDate.Format(time.RFC3339))
		fmt.Println()
	}
}

// printWorkItemsAsJSON prints work items in JSON format
func printWorkItemsAsJSON(workItems []AssignedWorkItem) {
	// Marshal the work items to JSON with indentation
	jsonData, err := json.MarshalIndent(workItems, "", "  ")
	if err != nil {
		logger.Error("Failed to marshal work items to JSON", "error", err)
		fmt.Println("Error: Failed to marshal work items to JSON:", err)
		return
	}

	// Print the JSON
	fmt.Println(string(jsonData))
}
