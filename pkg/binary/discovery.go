package binary

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// BinaryPrefix defines the valid prefixes for master-mold subcommand binaries
type BinaryPrefix string

const (
	// MMPrefix is the short prefix for master-mold binaries
	MMPrefix BinaryPrefix = "mm-"
	// MasterMoldPrefix is the long prefix for master-mold binaries
	MasterMoldPrefix BinaryPrefix = "master-mold-"
)

// ValidPrefixes returns all valid binary prefixes
func ValidPrefixes() []BinaryPrefix {
	return []BinaryPrefix{MMPrefix, MasterMoldPrefix}
}

// HasValidPrefix checks if a filename has a valid master-mold binary prefix
func HasValidPrefix(filename string) bool {
	for _, prefix := range ValidPrefixes() {
		if strings.HasPrefix(filename, string(prefix)) {
			return true
		}
	}
	return false
}

// IsExecutable checks if a file is executable
func IsExecutable(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.Mode()&0111 != 0
}

// FindInDirectory finds all master-mold binaries in a specific directory
func FindInDirectory(dir string) ([]string, error) {
	var binaries []string

	// Skip if the directory doesn't exist
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return binaries, nil
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read directory: %s", dir)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if HasValidPrefix(name) {
			fullPath := filepath.Join(dir, name)
			if IsExecutable(fullPath) {
				binaries = append(binaries, fullPath)
			}
		}
	}

	return binaries, nil
}

// FindInPath finds all master-mold binaries in the system PATH
func FindInPath() ([]string, error) {
	var allBinaries []string

	pathEnv := os.Getenv("PATH")
	paths := strings.Split(pathEnv, string(os.PathListSeparator))

	for _, path := range paths {
		binaries, err := FindInDirectory(path)
		if err != nil {
			// Skip directories we can't read
			continue
		}
		allBinaries = append(allBinaries, binaries...)
	}

	return allBinaries, nil
}

// FindAll finds all master-mold binaries in both the specified directory and PATH
func FindAll(baseDir string) ([]string, error) {
	// Create the base directory if it doesn't exist
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		if err := os.MkdirAll(baseDir, 0755); err != nil {
			return nil, errors.Wrap(err, "failed to create base directory")
		}
	}

	// Find binaries in the base directory
	baseDirBinaries, err := FindInDirectory(baseDir)
	if err != nil {
		return nil, err
	}

	// Find binaries in PATH
	pathBinaries, err := FindInPath()
	if err != nil {
		return nil, err
	}

	// Combine the results
	return append(baseDirBinaries, pathBinaries...), nil
}

// ExtractCommandName extracts the command name from a binary path
func ExtractCommandName(binaryPath string) string {
	commandName := filepath.Base(binaryPath)
	
	for _, prefix := range ValidPrefixes() {
		if strings.HasPrefix(commandName, string(prefix)) {
			return strings.TrimPrefix(commandName, string(prefix))
		}
	}
	
	return commandName
}