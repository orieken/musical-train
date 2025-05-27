package main

import (
	"os"

	"github.com/spf13/cobra"
	"log/slog"
)

var logger *slog.Logger

func main() {
	// Initialize the logger
	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	logger.Info("Starting Azure DevOps subcommand")

	// Create the root command
	var rootCmd = &cobra.Command{
		Use:     "azure-devops",
		Short:   "Manage Azure DevOps work items",
		Long:    "Provides commands to create and manage work items in Azure DevOps.",
		Aliases: []string{"ado"},
	}

	// Create the work-items subcommand
	var workItemsCmd = &cobra.Command{
		Use:   "work-items",
		Short: "Manage work items in Azure DevOps",
		Long:  "Create and manage work items in Azure DevOps.",
	}

	// Create the create subcommand
	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create work items from a JSON file",
		Long:  "Creates work items in Azure DevOps based on data provided in a JSON file.",
		Run:   createWorkItems,
	}

	// Create the template subcommand
	var templateCmd = &cobra.Command{
		Use:   "template",
		Short: "Generate a template JSON file for creating work items",
		Long:  "Generates a template JSON file that can be used as a starting point for creating work items.",
		Run:   generateWorkItemTemplate,
	}

	// Create the assigned subcommand
	var assignedCmd = &cobra.Command{
		Use:   "assigned",
		Short: "List work items assigned to a user",
		Long:  "Lists all work items assigned to a user and displays the work item and time logged.",
		Run:   listAssignedWorkItems,
	}

	// Create the pull-requests subcommand
	var prCmd = &cobra.Command{
		Use:   "pull-requests",
		Short: "Manage pull requests",
		Long:  "Provides commands to manage pull requests in Azure DevOps.",
	}

	// Create the list-open subcommand
	var listOpenCmd = &cobra.Command{
		Use:   "list-open",
		Short: "List all open pull requests",
		Long:  "Lists all open pull requests for all repositories in the organization.",
		Run:   listOpenPullRequests,
	}

	// Add flags to the commands
	createCmd.Flags().String("json", "", "Path to the JSON file containing work item definitions")
	createCmd.MarkFlagRequired("json")

	assignedCmd.Flags().String("user", "", "Username to filter work items by")
	assignedCmd.MarkFlagRequired("user")
	assignedCmd.Flags().Bool("json", false, "Output the results in JSON format")

	listOpenCmd.Flags().Bool("json", false, "Output the results in JSON format")

	// Add subcommands to their parent commands
	workItemsCmd.AddCommand(createCmd)
	workItemsCmd.AddCommand(templateCmd)
	workItemsCmd.AddCommand(assignedCmd)
	prCmd.AddCommand(listOpenCmd)

	rootCmd.AddCommand(workItemsCmd)
	rootCmd.AddCommand(prCmd)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		logger.Error("Error executing command", "error", err)
		os.Exit(1)
	}

	logger.Info("Azure DevOps subcommand completed successfully")
}

// Note: The implementations for the command handlers (createWorkItems and generateWorkItemTemplate)
// are defined in separate files (create.go and template.go)
