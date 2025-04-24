package test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestListBinariesCommand tests the list-binaries command
func TestListBinariesCommand(t *testing.T) {
	// Skip if the master-mold binary doesn't exist
	if _, err := os.Stat("../master-mold"); os.IsNotExist(err) {
		t.Skip("master-mold binary not found, skipping integration test")
	}

	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "master-mold-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a config directory in the temp directory
	configDir := filepath.Join(tempDir, "config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}

	// Create a config file in the config directory
	configFile := filepath.Join(configDir, "config.toml")
	configContent := `
base_dir = "${HOME}/.master-mold"
timeout = 10
`
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Create a .master-mold directory in the temp directory
	masterMoldDir := filepath.Join(tempDir, ".master-mold")
	if err := os.MkdirAll(masterMoldDir, 0755); err != nil {
		t.Fatalf("Failed to create .master-mold dir: %v", err)
	}

	// Create a test binary in the temp directory
	testBinary := filepath.Join(tempDir, "mm-test-binary")
	if err := os.WriteFile(testBinary, []byte("#!/bin/sh\necho 'Test binary executed'"), 0755); err != nil {
		t.Fatalf("Failed to create test binary: %v", err)
	}

	// Get the current PATH
	oldPath := os.Getenv("PATH")
	newPath := tempDir + string(os.PathListSeparator) + oldPath

	// Run the list-binaries command using go run
	cmd := exec.Command("go", "run", "../cmd/master-mold/main.go", "list-binaries")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Set environment variables for the command
	cmd.Env = append(os.Environ(), 
		fmt.Sprintf("HOME=%s", tempDir),
		fmt.Sprintf("PATH=%s", newPath))

	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to run list-binaries command: %v\nStderr: %s", err, stderr.String())
	}

	// Check that the output contains the test binary
	output := stdout.String()
	if !strings.Contains(output, "test-binary") {
		t.Errorf("list-binaries output does not contain test-binary: %s", output)
	}
}

// TestSubcommandExecution tests the execution of a subcommand
func TestSubcommandExecution(t *testing.T) {
	// Skip if the master-mold binary doesn't exist
	if _, err := os.Stat("../master-mold"); os.IsNotExist(err) {
		t.Skip("master-mold binary not found, skipping integration test")
	}

	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "master-mold-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a config directory in the temp directory
	configDir := filepath.Join(tempDir, "config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}

	// Create a config file in the config directory
	configFile := filepath.Join(configDir, "config.toml")
	configContent := `
base_dir = "${HOME}/.master-mold"
timeout = 10
`
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Create a .master-mold directory in the temp directory
	masterMoldDir := filepath.Join(tempDir, ".master-mold")
	if err := os.MkdirAll(masterMoldDir, 0755); err != nil {
		t.Fatalf("Failed to create .master-mold dir: %v", err)
	}

	// Create a test binary in the temp directory
	testBinary := filepath.Join(tempDir, "mm-test-command")
	if err := os.WriteFile(testBinary, []byte("#!/bin/sh\necho 'Test command executed with args: $@'"), 0755); err != nil {
		t.Fatalf("Failed to create test binary: %v", err)
	}

	// Get the current PATH
	oldPath := os.Getenv("PATH")
	newPath := tempDir + string(os.PathListSeparator) + oldPath

	// Run the test command using go run
	cmd := exec.Command("go", "run", "../cmd/master-mold/main.go", "test-command", "arg1", "arg2")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Set environment variables for the command
	cmd.Env = append(os.Environ(), 
		fmt.Sprintf("HOME=%s", tempDir),
		fmt.Sprintf("PATH=%s", newPath))

	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to run test command: %v\nStderr: %s", err, stderr.String())
	}

	// Check that the output contains the expected message
	output := stdout.String()
	if !strings.Contains(output, "Test command executed with args:") {
		t.Errorf("test command output does not contain expected message: %s", output)
	}
	if !strings.Contains(output, "arg1") || !strings.Contains(output, "arg2") {
		t.Errorf("test command output does not contain expected arguments: %s", output)
	}
}

// TestMain is the entry point for the integration tests
func TestMain(m *testing.M) {
	// Build the master-mold binary if it doesn't exist
	if _, err := os.Stat("../master-mold"); os.IsNotExist(err) {
		cmd := exec.Command("go", "build", "-o", "../master-mold", "../cmd/master-mold")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Failed to build master-mold binary: %v\n", err)
			os.Exit(1)
		}
	}

	// Run the tests
	os.Exit(m.Run())
}
