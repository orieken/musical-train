package main

import (
	"fmt"
	"os"
	"path/filepath"

	"log/slog"

	"github.com/oscarrieken/master-mold/pkg/binary"
	"github.com/oscarrieken/master-mold/pkg/display"
)

// initLogger initializes the logger
func initLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

// findBinaries finds all master-mold binaries
func findBinaries(logger *slog.Logger) ([]string, error) {
	// Get the home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	// Add the master-mold directory to the search paths
	masterMoldDir := filepath.Join(homeDir, ".master-mold")

	// Find all binaries
	return binary.FindAll(masterMoldDir)
}

// checkIfRunningAsSubcommand checks if we're running as a subcommand of master-mold
func checkIfRunningAsSubcommand(logger *slog.Logger) {
	isSubcommand, err := binary.IsRunningAsSubcommand()
	if err != nil {
		logger.Warn("Failed to determine if running as subcommand", "error", err)
		return
	}

	if isSubcommand {
		fmt.Println("\nRunning as a subcommand of master-mold")
	} else {
		fmt.Println("\nRunning as a standalone command")
	}
}

func main() {
	// Initialize the logger
	logger := initLogger()
	logger.Info("Running mm-list-binaries subcommand")

	// Find all binaries
	binaries, err := findBinaries(logger)
	if err != nil {
		logger.Error("Failed to find binaries", "error", err)
		os.Exit(1)
	}

	// Display the binaries
	display.PrintBinaryPaths(binaries)

	// Check if we're running as a subcommand
	checkIfRunningAsSubcommand(logger)

	logger.Info("mm-list-binaries completed successfully")
}
