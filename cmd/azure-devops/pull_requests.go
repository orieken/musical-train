package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/git"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// PullRequest represents a pull request in Azure DevOps
type PullRequest struct {
	Repository   string    `json:"repository"`
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Creator      string    `json:"creator"`
	Created      time.Time `json:"created"`
	Status       string    `json:"status"`
	TargetBranch string    `json:"targetBranch"`
}

// listOpenPullRequests lists all open pull requests for all repositories in the organization
func listOpenPullRequests(cmd *cobra.Command, args []string) {
	logger.Info("Listing open pull requests")

	// Check if JSON output is requested
	jsonOutput, err := cmd.Flags().GetBool("json")
	if err != nil {
		handleError("Failed to get json flag", err)
		return
	}

	// Get the pull requests
	pullRequests, err := getAllOpenPullRequests()
	if err != nil {
		handleError("Failed to get open pull requests", err)
		return
	}

	// Print the pull requests
	if jsonOutput {
		printPullRequestsAsJSON(pullRequests)
	} else {
		printPullRequestsAsText(pullRequests)
	}

	logger.Info("Pull requests listed successfully")
}

// getAllOpenPullRequests gets all open pull requests for all repositories in the organization
func getAllOpenPullRequests() ([]PullRequest, error) {
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

	// Get all projects
	projects, err := getProjects(connection)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get projects")
	}

	// Get all repositories and pull requests
	var allPullRequests []PullRequest
	for _, project := range projects {
		// Get all repositories for the project
		repositories, err := getRepositories(connection, *project.Name)
		if err != nil {
			logger.Warn("Failed to get repositories for project", "project", *project.Name, "error", err)
			continue
		}

		// Get all pull requests for each repository
		for _, repo := range repositories {
			pullRequests, err := getPullRequests(connection, *project.Name, *repo.Name)
			if err != nil {
				logger.Warn("Failed to get pull requests for repository", "repository", *repo.Name, "error", err)
				continue
			}

			allPullRequests = append(allPullRequests, pullRequests...)
		}
	}

	return allPullRequests, nil
}

// getProjects gets all projects in the organization
func getProjects(connection *azuredevops.Connection) ([]core.TeamProjectReference, error) {
	// Create a client for the Core API
	client, err := core.NewClient(context.Background(), connection)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Core client")
	}

	// Get all projects
	projects, err := client.GetProjects(context.Background(), core.GetProjectsArgs{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get projects")
	}

	return projects.Value, nil
}

// getRepositories gets all repositories for a project
func getRepositories(connection *azuredevops.Connection, projectName string) ([]git.GitRepository, error) {
	// Create a client for the Git API
	client, err := git.NewClient(context.Background(), connection)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Git client")
	}

	// Get all repositories for the project
	repositories, err := client.GetRepositories(context.Background(), git.GetRepositoriesArgs{
		Project: &projectName,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get repositories")
	}

	return *repositories, nil
}

// getPullRequests gets all open pull requests for a repository
func getPullRequests(connection *azuredevops.Connection, projectName, repositoryName string) ([]PullRequest, error) {
	// Create a client for the Git API
	client, err := git.NewClient(context.Background(), connection)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Git client")
	}

	// Set up the pull request search criteria
	status := git.PullRequestStatusValues.Active
	searchCriteria := git.GitPullRequestSearchCriteria{
		Status: &status,
	}

	// Get all pull requests for the repository
	pullRequests, err := client.GetPullRequests(context.Background(), git.GetPullRequestsArgs{
		Project:      &projectName,
		RepositoryId: &repositoryName,
		SearchCriteria: &searchCriteria,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get pull requests")
	}

	// Convert the pull requests to our model
	var result []PullRequest
	for _, pr := range *pullRequests {
		result = append(result, PullRequest{
			Repository:   repositoryName,
			ID:           *pr.PullRequestId,
			Title:        *pr.Title,
			Creator:      *pr.CreatedBy.DisplayName,
			Created:      pr.CreationDate.Time,
			Status:       string(*pr.Status),
			TargetBranch: *pr.TargetRefName,
		})
	}

	return result, nil
}

// printPullRequestsAsText prints pull requests in a human-readable format
func printPullRequestsAsText(pullRequests []PullRequest) {
	if len(pullRequests) == 0 {
		fmt.Println("No open pull requests found.")
		return
	}

	fmt.Printf("Found %d open pull requests:\n\n", len(pullRequests))

	for _, pr := range pullRequests {
		fmt.Printf("Repository: %s\n", pr.Repository)
		fmt.Printf("ID: %d\n", pr.ID)
		fmt.Printf("Title: %s\n", pr.Title)
		fmt.Printf("Creator: %s\n", pr.Creator)
		fmt.Printf("Created: %s\n", pr.Created.Format(time.RFC3339))
		fmt.Printf("Status: %s\n", pr.Status)
		fmt.Printf("Target Branch: %s\n", pr.TargetBranch)
		fmt.Println()
	}
}

// printPullRequestsAsJSON prints pull requests in JSON format
func printPullRequestsAsJSON(pullRequests []PullRequest) {
	// Marshal the pull requests to JSON with indentation
	jsonData, err := json.MarshalIndent(pullRequests, "", "  ")
	if err != nil {
		logger.Error("Failed to marshal pull requests to JSON", "error", err)
		fmt.Println("Error: Failed to marshal pull requests to JSON:", err)
		return
	}

	// Print the JSON
	fmt.Println(string(jsonData))
}
