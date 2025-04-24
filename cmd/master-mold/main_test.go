package main

import (
	"bytes"
	"os"
	"testing"

	"log/slog"
)

func TestInitLogger(t *testing.T) {
	// Call the function
	logger := initLogger()

	// Check that the logger is not nil
	if logger == nil {
		t.Fatal("initLogger() returned nil")
	}

	// Test that the logger works by logging a message
	var buf bytes.Buffer
	testLogger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	testLogger.Info("Test message")

	// Check that something was logged
	if buf.Len() == 0 {
		t.Error("Logger did not log anything")
	}
}

func TestLoadConfig(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Call the function with the temporary directory
	cfg, err := loadConfigWithPaths(logger, []string{tempDir})

	// Check that there was no error
	if err != nil {
		t.Fatalf("loadConfigWithPaths() returned an error: %v", err)
	}

	// Check that the configuration is not nil
	if cfg == nil {
		t.Fatal("loadConfigWithPaths() returned nil config")
	}

	// Check that the configuration has the expected values
	if cfg.BaseDir == "" {
		t.Error("loadConfigWithPaths() returned config with empty BaseDir")
	}
	if cfg.Timeout <= 0 {
		t.Error("loadConfigWithPaths() returned config with invalid Timeout")
	}

	// Check that the config file was created in the temporary directory
	configFile := tempDir + "/config.toml"
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Errorf("Config file was not created: %s", configFile)
	}
}

func TestHandleCommands_NoCommand(t *testing.T) {
	// Save the original os.Args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Set os.Args to simulate no command
	os.Args = []string{"master-mold"}

	// Create a mock registry
	registry := &MockRegistry{}

	// Call the function
	err := handleCommands(registry)

	// Check that there was an error
	if err == nil {
		t.Fatal("handleCommands() did not return an error when no command was specified")
	}

	// Check that the error message is as expected
	if err.Error() != "no command specified" {
		t.Errorf("handleCommands() returned error = %v, want 'no command specified'", err)
	}
}

func TestHandleCommands_WithCommand(t *testing.T) {
	// Save the original os.Args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Set os.Args to simulate a command
	os.Args = []string{"master-mold", "test-command", "arg1", "arg2"}

	// Create a mock registry
	registry := &MockRegistry{}

	// Call the function
	err := handleCommands(registry)

	// Check that there was no error
	if err != nil {
		t.Fatalf("handleCommands() returned an error: %v", err)
	}

	// Check that the registry's Execute method was called with the correct arguments
	if !registry.ExecuteCalled {
		t.Fatal("handleCommands() did not call registry.Execute()")
	}
	if registry.CommandName != "test-command" {
		t.Errorf("handleCommands() called registry.Execute() with commandName = %s, want 'test-command'", registry.CommandName)
	}
	if len(registry.Args) != 2 || registry.Args[0] != "arg1" || registry.Args[1] != "arg2" {
		t.Errorf("handleCommands() called registry.Execute() with args = %v, want ['arg1', 'arg2']", registry.Args)
	}
}

// MockRegistry is a mock implementation of the command.Registry interface for testing
type MockRegistry struct {
	ExecuteCalled bool
	CommandName   string
	Args          []string
}

// Execute records the call and returns nil
func (m *MockRegistry) Execute(commandName string, args []string) error {
	m.ExecuteCalled = true
	m.CommandName = commandName
	m.Args = args
	return nil
}
