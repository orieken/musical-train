package binary

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
	"log/slog"
)

// FindExecutable finds the executable for a given command name
func FindExecutable(command string, baseDir string) (string, error) {
	// Try both naming conventions
	binNames := []string{
		string(MMPrefix) + command,
		string(MasterMoldPrefix) + command,
	}

	// Expand environment variables in the base directory
	expandedBaseDir := os.ExpandEnv(baseDir)

	// First, look for the binary in PATH
	for _, binName := range binNames {
		cmdPath, err := exec.LookPath(binName)
		if err == nil {
			// Found the binary in PATH
			return cmdPath, nil
		}
	}

	// If not found in PATH, check the base directory
	for _, binName := range binNames {
		fullPath := filepath.Join(expandedBaseDir, binName)
		if IsExecutable(fullPath) {
			return fullPath, nil
		}
	}

	return "", errors.Errorf("subcommand '%s' not found", command)
}

// Execute executes a subcommand binary
func Execute(cmdPath string, args []string, logger *slog.Logger) error {
	logger.Info("Executing binary", "path", cmdPath, "args", args)

	// Create the command
	cmd := exec.Command(cmdPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Execute the command
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "failed to execute binary '%s'", cmdPath)
	}

	return nil
}

// ExecuteSubcommand finds and executes a subcommand
func ExecuteSubcommand(command string, args []string, baseDir string, logger *slog.Logger) error {
	// Find the executable
	cmdPath, err := FindExecutable(command, baseDir)
	if err != nil {
		return err
	}

	logger.Info("Executing subcommand", "command", command, "binary", cmdPath)

	// Execute the command
	return Execute(cmdPath, args, logger)
}

// IsRunningAsSubcommand checks if the current process is running as a subcommand of master-mold
func IsRunningAsSubcommand() (bool, error) {
	parentPID := os.Getppid()
	parentProcess, err := exec.Command("ps", "-o", "comm=", "-p", fmt.Sprintf("%d", parentPID)).Output()
	if err != nil {
		return false, errors.Wrap(err, "failed to get parent process information")
	}

	parentName := string(parentProcess)
	return filepath.Base(parentName) == "master-mold", nil
}
