package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/microsoft/azure-devops-go-api/azuredevops/webapi"
	"github.com/microsoft/azure-devops-go-api/azuredevops/workitemtracking"
)

func TestReadWorkItemsFromFile(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "workitem-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test work item fields
	workItemFields := []WorkItemField{
		{
			Op:    "add",
			Path:  "/fields/System.Title",
			Value: "Test Title",
		},
		{
			Op:    "add",
			Path:  "/fields/System.WorkItemType",
			Value: "Task",
		},
	}

	// Marshal the work item fields to JSON
	jsonData, err := json.MarshalIndent(workItemFields, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal work item fields: %v", err)
	}

	// Create a temporary file path
	tempFile := filepath.Join(tempDir, "test-workitems.json")

	// Write the JSON to a file
	if err := os.WriteFile(tempFile, jsonData, 0644); err != nil {
		t.Fatalf("Failed to write work item fields to file: %v", err)
	}

	// Read the work item fields from the file
	readFields, err := readWorkItemsFromFile(tempFile)
	if err != nil {
		t.Fatalf("readWorkItemsFromFile() returned an error: %v", err)
	}

	// Check that the fields were read correctly
	if len(readFields) != len(workItemFields) {
		t.Errorf("readWorkItemsFromFile() read %d fields, expected %d", len(readFields), len(workItemFields))
	}

	// Check that the fields match
	for i, expectedField := range workItemFields {
		if i >= len(readFields) {
			t.Errorf("readWorkItemsFromFile() read fewer fields than expected")
			break
		}

		actualField := readFields[i]
		if actualField.Op != expectedField.Op {
			t.Errorf("readFields[%d].Op = %s, expected %s", i, actualField.Op, expectedField.Op)
		}
		if actualField.Path != expectedField.Path {
			t.Errorf("readFields[%d].Path = %s, expected %s", i, actualField.Path, expectedField.Path)
		}
		if actualField.Value != expectedField.Value {
			t.Errorf("readFields[%d].Value = %s, expected %s", i, actualField.Value, expectedField.Value)
		}
	}
}

func TestGroupFieldsByWorkItemType(t *testing.T) {
	// Create test work item fields
	workItemFields := []WorkItemField{
		{
			Op:    "add",
			Path:  "/fields/System.Title",
			Value: "Test Title",
		},
		{
			Op:    "add",
			Path:  "/fields/System.WorkItemType",
			Value: "Task",
		},
		{
			Op:    "add",
			Path:  "/fields/System.Description",
			Value: "Test Description",
		},
	}

	// Group the fields by work item type
	groupedFields := groupFieldsByWorkItemType(workItemFields)

	// Check that the fields were grouped correctly
	if len(groupedFields) != 1 {
		t.Errorf("groupFieldsByWorkItemType() returned %d groups, expected 1", len(groupedFields))
	}

	// Check that the group has the expected work item type
	if _, ok := groupedFields["Task"]; !ok {
		t.Errorf("groupFieldsByWorkItemType() did not create a group for 'Task'")
	}

	// Check that the group has all the fields
	taskFields := groupedFields["Task"]
	if len(taskFields) != len(workItemFields) {
		t.Errorf("groupFieldsByWorkItemType() grouped %d fields for 'Task', expected %d", len(taskFields), len(workItemFields))
	}
}

func TestGroupFieldsByWorkItemTypeWithoutType(t *testing.T) {
	// Create test work item fields without a work item type
	workItemFields := []WorkItemField{
		{
			Op:    "add",
			Path:  "/fields/System.Title",
			Value: "Test Title",
		},
		{
			Op:    "add",
			Path:  "/fields/System.Description",
			Value: "Test Description",
		},
	}

	// Group the fields by work item type
	groupedFields := groupFieldsByWorkItemType(workItemFields)

	// Check that the fields were grouped correctly
	if len(groupedFields) != 1 {
		t.Errorf("groupFieldsByWorkItemType() returned %d groups, expected 1", len(groupedFields))
	}

	// Check that the group has the default work item type
	if _, ok := groupedFields["Task"]; !ok {
		t.Errorf("groupFieldsByWorkItemType() did not create a group for the default 'Task'")
	}

	// Check that the group has all the fields
	taskFields := groupedFields["Task"]
	if len(taskFields) != len(workItemFields) {
		t.Errorf("groupFieldsByWorkItemType() grouped %d fields for 'Task', expected %d", len(taskFields), len(workItemFields))
	}
}

func TestConvertFieldsToPatches(t *testing.T) {
	// Create test work item fields
	workItemFields := []WorkItemField{
		{
			Op:    "add",
			Path:  "/fields/System.Title",
			Value: "Test Title",
		},
		{
			Op:    "add",
			Path:  "/fields/System.WorkItemType",
			Value: "Task",
		},
	}

	// Convert the fields to patches
	patches := convertFieldsToPatches(workItemFields)

	// Check that the patches were created correctly
	if len(patches) != len(workItemFields) {
		t.Errorf("convertFieldsToPatches() created %d patches, expected %d", len(patches), len(workItemFields))
	}

	// Check that the patches have the expected values
	for i, patch := range patches {
		expectedField := workItemFields[i]

		// Check the operation
		expectedOp := webapi.OperationValues.Add
		if *patch.Op != expectedOp {
			t.Errorf("patches[%d].Op = %s, expected %s", i, *patch.Op, expectedOp)
		}

		// Check the path
		if *patch.Path != expectedField.Path {
			t.Errorf("patches[%d].Path = %s, expected %s", i, *patch.Path, expectedField.Path)
		}

		// Check the value
		if patch.Value != expectedField.Value {
			t.Errorf("patches[%d].Value = %v, expected %s", i, patch.Value, expectedField.Value)
		}
	}
}

func TestGetJSONFilePath(t *testing.T) {
	// Create a mock flag provider
	provider := &MockFlagProvider{
		flags: map[string]string{
			"json": "test-file.json",
		},
	}

	// Get the JSON file path
	path, err := getJSONFilePath(provider)
	if err != nil {
		t.Fatalf("getJSONFilePath() returned an error: %v", err)
	}

	// Check that the path is correct
	expectedPath := "test-file.json"
	if path != expectedPath {
		t.Errorf("getJSONFilePath() returned %s, expected %s", path, expectedPath)
	}
}

func TestGetWorkItemTitle(t *testing.T) {
	// Create a test work item with a title
	title := "Test Title"
	fields := make(map[string]interface{})
	fields["System.Title"] = title
	workItem := workitemtracking.WorkItem{
		Fields: &fields,
	}

	// Get the title
	actualTitle := getWorkItemTitle(workItem)

	// Check that the title is correct
	if actualTitle != title {
		t.Errorf("getWorkItemTitle() returned %s, expected %s", actualTitle, title)
	}
}

func TestGetWorkItemTitleWithNilFields(t *testing.T) {
	// Create a test work item with nil fields
	workItem := workitemtracking.WorkItem{
		Fields: nil,
	}

	// Get the title
	actualTitle := getWorkItemTitle(workItem)

	// Check that the title is the default
	expectedTitle := "Unknown"
	if actualTitle != expectedTitle {
		t.Errorf("getWorkItemTitle() returned %s, expected %s", actualTitle, expectedTitle)
	}
}

func TestGetWorkItemTitleWithoutTitle(t *testing.T) {
	// Create a test work item without a title
	fields := make(map[string]interface{})
	workItem := workitemtracking.WorkItem{
		Fields: &fields,
	}

	// Get the title
	actualTitle := getWorkItemTitle(workItem)

	// Check that the title is the default
	expectedTitle := "Unknown"
	if actualTitle != expectedTitle {
		t.Errorf("getWorkItemTitle() returned %s, expected %s", actualTitle, expectedTitle)
	}
}

// MockFlagProvider is a mock implementation of the FlagProvider interface for testing
type MockFlagProvider struct {
	flags map[string]string
}

// GetStringFlag returns a string flag value
func (p *MockFlagProvider) GetStringFlag(name string) (string, error) {
	if value, ok := p.flags[name]; ok {
		return value, nil
	}
	return "", nil
}
