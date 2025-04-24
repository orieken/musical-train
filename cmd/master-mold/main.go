package main

import (
	"fmt"
	"os"

	"log/slog"

	"github.com/oscarrieken/master-mold/pkg/command"
	"github.com/oscarrieken/master-mold/pkg/config"
)

// initLogger initializes the logger
func initLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

// loadConfig loads the configuration
func loadConfig(logger *slog.Logger) (*config.Config, error) {
	return loadConfigWithPaths(logger, []string{
		"./config",
		"$HOME/.master-mold",
	})
}

// loadConfigWithPaths loads the configuration from the specified paths
func loadConfigWithPaths(logger *slog.Logger, configPaths []string) (*config.Config, error) {
	return config.LoadConfig(configPaths, logger)
}

// CommandExecutor is an interface for executing commands
type CommandExecutor interface {
	Execute(commandName string, args []string) error
}

// handleCommands handles command execution
func handleCommands(registry CommandExecutor) error {
	if len(os.Args) < 2 {
		fmt.Println("Usage: master-mold <command> [options]")
		fmt.Println("Run 'master-mold list-binaries' to see available commands")
		return fmt.Errorf("no command specified")
	}

	commandName := os.Args[1]
	args := os.Args[2:]

	return registry.Execute(commandName, args)
}

func main() {
	// Initialize the logger
	logger := initLogger()
	logger.Info("Starting master-mold CLI")

	// Load the configuration
	cfg, err := loadConfig(logger)
	if err != nil {
		logger.Error("Error loading configuration", "error", err)
		os.Exit(1)
	}

	// Create the command registry
	registry := command.NewRegistry(cfg, logger)
	command.RegisterCommands(registry)

	// Handle commands
	if err := handleCommands(registry); err != nil {
		logger.Error("Error executing command", "error", err)
		os.Exit(1)
	}

	logger.Info("Command completed successfully")
}
