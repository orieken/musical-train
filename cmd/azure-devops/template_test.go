package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestCreateWorkItemTemplate(t *testing.T) {
	// Call the function
	template, err := createWorkItemTemplate()

	// Check for errors
	if err != nil {
		t.Fatalf("createWorkItemTemplate() returned an error: %v", err)
	}

	// Check that the template is not nil
	if template == nil {
		t.Fatal("createWorkItemTemplate() returned nil template")
	}

	// Check that the template has the expected number of fields
	expectedFieldCount := 5
	if len(template) != expectedFieldCount {
		t.Errorf("createWorkItemTemplate() returned %d fields, expected %d", len(template), expectedFieldCount)
	}

	// Check that all fields have the expected operation type
	for i, field := range template {
		if field.Op != "add" {
			t.Errorf("template[%d].Op = %s, expected 'add'", i, field.Op)
		}
	}

	// Check that the template contains the required fields
	requiredPaths := []string{
		"/fields/System.Title",
		"/fields/System.WorkItemType",
		"/fields/System.Description",
		"/fields/System.AreaPath",
		"/fields/System.IterationPath",
	}

	for _, requiredPath := range requiredPaths {
		found := false
		for _, field := range template {
			if field.Path == requiredPath {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("template does not contain required field: %s", requiredPath)
		}
	}
}

func TestWriteTemplateToFile(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "template-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test template
	template := []WorkItemField{
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

	// Create a temporary file path
	tempFile := filepath.Join(tempDir, "test-template.json")

	// Write the template to the file
	err = writeTemplateToFile(template, tempFile)
	if err != nil {
		t.Fatalf("writeTemplateToFile() returned an error: %v", err)
	}

	// Check that the file exists
	if _, err := os.Stat(tempFile); os.IsNotExist(err) {
		t.Fatalf("writeTemplateToFile() did not create the file: %s", tempFile)
	}

	// Read the file content
	data, err := os.ReadFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to read the template file: %v", err)
	}

	// Unmarshal the JSON
	var readTemplate []WorkItemField
	if err := json.Unmarshal(data, &readTemplate); err != nil {
		t.Fatalf("Failed to unmarshal the template JSON: %v", err)
	}

	// Check that the template was written correctly
	if len(readTemplate) != len(template) {
		t.Errorf("writeTemplateToFile() wrote %d fields, expected %d", len(readTemplate), len(template))
	}

	// Check that the fields match
	for i, expectedField := range template {
		if i >= len(readTemplate) {
			t.Errorf("writeTemplateToFile() wrote fewer fields than expected")
			break
		}

		actualField := readTemplate[i]
		if actualField.Op != expectedField.Op {
			t.Errorf("readTemplate[%d].Op = %s, expected %s", i, actualField.Op, expectedField.Op)
		}
		if actualField.Path != expectedField.Path {
			t.Errorf("readTemplate[%d].Path = %s, expected %s", i, actualField.Path, expectedField.Path)
		}
		if actualField.Value != expectedField.Value {
			t.Errorf("readTemplate[%d].Value = %s, expected %s", i, actualField.Value, expectedField.Value)
		}
	}
}

func TestGenerateTemplateFile(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "template-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a temporary file path
	tempFile := filepath.Join(tempDir, "test-template.json")

	// Generate the template file
	err = generateTemplateFile(tempFile)
	if err != nil {
		t.Fatalf("generateTemplateFile() returned an error: %v", err)
	}

	// Check that the file exists
	if _, err := os.Stat(tempFile); os.IsNotExist(err) {
		t.Fatalf("generateTemplateFile() did not create the file: %s", tempFile)
	}

	// Read the file content
	data, err := os.ReadFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to read the template file: %v", err)
	}

	// Unmarshal the JSON
	var template []WorkItemField
	if err := json.Unmarshal(data, &template); err != nil {
		t.Fatalf("Failed to unmarshal the template JSON: %v", err)
	}

	// Check that the template has the expected number of fields
	expectedFieldCount := 5
	if len(template) != expectedFieldCount {
		t.Errorf("generateTemplateFile() created a template with %d fields, expected %d", len(template), expectedFieldCount)
	}
}

func TestHandleTemplateErrorWithExit(t *testing.T) {
	// Create a test error
	testErr := fmt.Errorf("test error")
	testMessage := "test message"

	// Create a custom exit function that records the exit code
	var exitCode int
	testExitFunc := func(code int) {
		exitCode = code
	}

	// Redirect stdout to capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Create a logger for the test
	oldLogger := logger
	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Call the function
	handleTemplateErrorWithExit(testMessage, testErr, testExitFunc)

	// Restore stdout and logger
	w.Close()
	os.Stdout = oldStdout
	logger = oldLogger

	// Read the captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Check that the exit code is 1
	if exitCode != 1 {
		t.Errorf("handleTemplateErrorWithExit() exit code = %d, want 1", exitCode)
	}

	// Check that the output contains the expected message
	expectedOutput := fmt.Sprintf("Error: %s: %v", testMessage, testErr)
	if !strings.Contains(output, expectedOutput) {
		t.Errorf("handleTemplateErrorWithExit() output = %q, want to contain %q", output, expectedOutput)
	}
}

func TestPrintTemplateSuccessMessage(t *testing.T) {
	// Create a temporary file path
	filename := "test-template.json"

	// Redirect stdout to capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the function
	printTemplateSuccessMessage(filename)

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read the captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Check that the output contains the expected messages
	if !strings.Contains(output, "Template file created:") {
		t.Errorf("printTemplateSuccessMessage() output does not contain 'Template file created:': %s", output)
	}
	if !strings.Contains(output, filename) {
		t.Errorf("printTemplateSuccessMessage() output does not contain the filename: %s", output)
	}
	if !strings.Contains(output, "Edit this file and use it with the 'create' command") {
		t.Errorf("printTemplateSuccessMessage() output does not contain the expected instructions: %s", output)
	}
}
