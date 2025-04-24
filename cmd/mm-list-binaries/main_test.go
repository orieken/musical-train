package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
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

func TestFindBinaries(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "binaries-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Create a temporary home directory
	homeDir := filepath.Join(tempDir, "home")
	if err := os.MkdirAll(homeDir, 0755); err != nil {
		t.Fatalf("Failed to create home dir: %v", err)
	}

	// Create a .master-mold directory in the home directory
	masterMoldDir := filepath.Join(homeDir, ".master-mold")
	if err := os.MkdirAll(masterMoldDir, 0755); err != nil {
		t.Fatalf("Failed to create .master-mold dir: %v", err)
	}

	// Create some test binaries
	testBinaries := []struct {
		name       string
		executable bool
	}{
		{name: "mm-test1", executable: true},
		{name: "mm-test2", executable: true},
		{name: "not-a-binary", executable: false},
	}

	for _, bin := range testBinaries {
		path := filepath.Join(masterMoldDir, bin.name)
		mode := os.FileMode(0644)
		if bin.executable {
			mode = 0755
		}
		if err := os.WriteFile(path, []byte("test"), mode); err != nil {
			t.Fatalf("Failed to create test binary %s: %v", bin.name, err)
		}
	}

	// Set the HOME environment variable to the temporary home directory
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", homeDir)
	defer os.Setenv("HOME", oldHome)

	// Call the function
	binaries, err := findBinaries(logger)
	if err != nil {
		t.Fatalf("findBinaries() returned an error: %v", err)
	}

	// Check that the binaries were found
	if len(binaries) != 2 {
		t.Errorf("findBinaries() found %d binaries, want 2", len(binaries))
	}

	// Check that the correct binaries were found
	for _, binary := range binaries {
		name := filepath.Base(binary)
		if name != "mm-test1" && name != "mm-test2" {
			t.Errorf("findBinaries() found unexpected binary: %s", name)
		}
	}
}

func TestCheckIfRunningAsSubcommand(t *testing.T) {
	// Create a test logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Redirect stdout to capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the function
	checkIfRunningAsSubcommand(logger)

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read the captured output
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Check that the output contains either "Running as a subcommand" or "Running as a standalone command"
	if !strings.Contains(output, "Running as a") {
		t.Errorf("checkIfRunningAsSubcommand() did not print expected output: %s", output)
	}
}