package config

import (
	"os"
	"path/filepath"
	"testing"

	"log/slog"
)

func TestDefaultConfig(t *testing.T) {
	// Get the default configuration
	config := DefaultConfig()

	// Check that the default values are set correctly
	if config.BaseDir != "${HOME}/.master-mold" {
		t.Errorf("DefaultConfig().BaseDir = %s, want ${HOME}/.master-mold", config.BaseDir)
	}
	if config.Timeout != 10 {
		t.Errorf("DefaultConfig().Timeout = %d, want 10", config.Timeout)
	}
}

func TestGetExpandedBaseDir(t *testing.T) {
	// Create a test configuration
	config := &Config{
		BaseDir: "${HOME}/test-dir",
	}

	// Get the expanded base directory
	expandedDir := GetExpandedBaseDir(config)

	// Check that the environment variables are expanded
	homeDir, _ := os.UserHomeDir()
	expectedDir := filepath.Join(homeDir, "test-dir")
	if expandedDir != expectedDir {
		t.Errorf("GetExpandedBaseDir() = %s, want %s", expandedDir, expectedDir)
	}
}

func TestEnsureBaseDirExists(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test configuration with a non-existent directory
	testDir := filepath.Join(tempDir, "test-dir")
	config := &Config{
		BaseDir: testDir,
	}

	// Ensure the base directory exists
	if err := EnsureBaseDirExists(config); err != nil {
		t.Fatalf("EnsureBaseDirExists() returned an error: %v", err)
	}

	// Check that the directory was created
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		t.Errorf("EnsureBaseDirExists() did not create the directory: %s", testDir)
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

	// Test loading configuration when no config file exists
	t.Run("no config file", func(t *testing.T) {
		// Load the configuration
		config, err := LoadConfig([]string{tempDir}, logger)
		if err != nil {
			t.Fatalf("LoadConfig() returned an error: %v", err)
		}

		// Check that the default values are set
		if config.BaseDir != "${HOME}/.master-mold" {
			t.Errorf("LoadConfig().BaseDir = %s, want ${HOME}/.master-mold", config.BaseDir)
		}
		if config.Timeout != 10 {
			t.Errorf("LoadConfig().Timeout = %d, want 10", config.Timeout)
		}

		// Check that the config file was created
		configFile := filepath.Join(tempDir, "config.toml")
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			t.Errorf("LoadConfig() did not create the config file: %s", configFile)
		}
	})

	// Test loading configuration from an existing config file
	t.Run("existing config file", func(t *testing.T) {
		// Create a config file
		configFile := filepath.Join(tempDir, "config.toml")
		configContent := `
base_dir = "/custom/dir"
timeout = 20
`
		if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
			t.Fatalf("Failed to create config file: %v", err)
		}

		// Load the configuration
		config, err := LoadConfig([]string{tempDir}, logger)
		if err != nil {
			t.Fatalf("LoadConfig() returned an error: %v", err)
		}

		// Check that the values from the config file are used
		if config.BaseDir != "/custom/dir" {
			t.Errorf("LoadConfig().BaseDir = %s, want /custom/dir", config.BaseDir)
		}
		if config.Timeout != 20 {
			t.Errorf("LoadConfig().Timeout = %d, want 20", config.Timeout)
		}
	})
}