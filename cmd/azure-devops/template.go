package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// WorkItemField represents a field in an Azure DevOps work item
type WorkItemField struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"`
}

// DefaultTemplateFileName is the default name for the template file
const DefaultTemplateFileName = "work-item-template.json"

// generateWorkItemTemplate generates a template JSON file for creating work items
func generateWorkItemTemplate(cmd *cobra.Command, args []string) {
	generateWorkItemTemplateWithExit(cmd, args, defaultExitFunc)
}

// generateWorkItemTemplateWithExit generates a template JSON file for creating work items with a custom exit function
func generateWorkItemTemplateWithExit(cmd *cobra.Command, args []string, exit exitFunc) {
	logger.Info("Generating work item template")

	// Generate the template file
	err := generateTemplateFile(DefaultTemplateFileName)
	if err != nil {
		handleTemplateErrorWithExit("Failed to generate template", err, exit)
		return
	}

	// Print success message
	printTemplateSuccessMessage(DefaultTemplateFileName)

	logger.Info("Template generation completed successfully")
}

// generateTemplateFile generates a template file with the given filename
func generateTemplateFile(filename string) error {
	// Create the template
	template, err := createWorkItemTemplate()
	if err != nil {
		return errors.Wrap(err, "failed to create work item template")
	}

	// Write the template to a file
	if err := writeTemplateToFile(template, filename); err != nil {
		return errors.Wrap(err, "failed to write template to file")
	}

	return nil
}

// exitFunc is a function that exits the program
type exitFunc func(int)

// defaultExitFunc is the default exit function
var defaultExitFunc exitFunc = os.Exit

// handleTemplateError handles errors during template generation
func handleTemplateError(message string, err error) {
	handleTemplateErrorWithExit(message, err, defaultExitFunc)
}

// handleTemplateErrorWithExit handles errors during template generation with a custom exit function
func handleTemplateErrorWithExit(message string, err error, exit exitFunc) {
	logger.Error(message, "error", err)
	fmt.Printf("Error: %s: %v\n", message, err)
	exit(1)
}

// printTemplateSuccessMessage prints a success message after template generation
func printTemplateSuccessMessage(filename string) {
	absPath, err := filepath.Abs(filename)
	if err != nil {
		absPath = filename
	}
	fmt.Printf("Template file created: %s\n", absPath)
	fmt.Println("Edit this file and use it with the 'create' command to create work items.")
}

// createWorkItemTemplate creates a template for an Azure DevOps work item
func createWorkItemTemplate() ([]WorkItemField, error) {
	// Create a template with common fields
	template := []WorkItemField{
		{
			Op:    "add",
			Path:  "/fields/System.Title",
			Value: "Example Title: Update this value",
		},
		{
			Op:    "add",
			Path:  "/fields/System.WorkItemType",
			Value: "Task | Bug | User Story | Feature",
		},
		{
			Op:    "add",
			Path:  "/fields/System.Description",
			Value: "Example Description: Provide a detailed description here.",
		},
		{
			Op:    "add",
			Path:  "/fields/System.AreaPath",
			Value: "YourProject\\YourArea",
		},
		{
			Op:    "add",
			Path:  "/fields/System.IterationPath",
			Value: "YourProject\\Iteration 1",
		},
	}

	return template, nil
}

// writeTemplateToFile writes the template to a file
func writeTemplateToFile(template []WorkItemField, filename string) error {
	// Marshal the template to JSON with indentation
	jsonData, err := json.MarshalIndent(template, "", "  ")
	if err != nil {
		return errors.Wrap(err, "failed to marshal template to JSON")
	}

	// Write the JSON to a file
	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
		return errors.Wrap(err, "failed to write template to file")
	}

	return nil
}
