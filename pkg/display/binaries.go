package display

import (
	"fmt"

	"github.com/oscarrieken/master-mold/pkg/binary"
)

// BinaryInfo represents information about a binary
type BinaryInfo struct {
	Name     string
	FullPath string
}

// FormatBinaryInfo formats binary information for display
func FormatBinaryInfo(info BinaryInfo) string {
	return fmt.Sprintf("  - %s (%s)", info.Name, info.FullPath)
}

// ProcessBinaries processes a list of binary paths and returns unique binary information
func ProcessBinaries(binaryPaths []string) []BinaryInfo {
	var result []BinaryInfo
	seenCommands := make(map[string]bool)

	for _, binaryPath := range binaryPaths {
		commandName := binary.ExtractCommandName(binaryPath)

		// Skip if we've already seen this command
		if seenCommands[commandName] {
			continue
		}

		seenCommands[commandName] = true
		result = append(result, BinaryInfo{
			Name:     commandName,
			FullPath: binaryPath,
		})
	}

	return result
}

// PrintBinaries prints a list of binaries to stdout
func PrintBinaries(binaries []BinaryInfo) {
	if len(binaries) == 0 {
		fmt.Println("No subcommand binaries found.")
		return
	}

	fmt.Println("Available subcommands:")
	for _, info := range binaries {
		fmt.Println(FormatBinaryInfo(info))
	}
}

// PrintBinaryPaths prints a list of binary paths to stdout
func PrintBinaryPaths(binaryPaths []string) {
	binaries := ProcessBinaries(binaryPaths)
	PrintBinaries(binaries)
}
