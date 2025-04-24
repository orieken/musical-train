package command

import (
	"github.com/pkg/errors"
	"github.com/oscarrieken/master-mold/pkg/binary"
	"github.com/oscarrieken/master-mold/pkg/config"
)

// SubcommandExecutor executes subcommands
type SubcommandExecutor struct {
	config *config.Config
	registry *Registry
}

// NewSubcommandExecutor creates a new subcommand executor
func NewSubcommandExecutor(config *config.Config, registry *Registry) *SubcommandExecutor {
	return &SubcommandExecutor{
		config: config,
		registry: registry,
	}
}

// Execute executes a subcommand
func (e *SubcommandExecutor) Execute(name string, args []string) error {
	// Find the executable
	baseDir := config.GetExpandedBaseDir(e.config)
	cmdPath, err := binary.FindExecutable(name, baseDir)
	if err != nil {
		return errors.Wrapf(err, "subcommand '%s' not found", name)
	}

	// Execute the command
	return binary.Execute(cmdPath, args, e.registry.Logger())
}

// RegisterSubcommandExecutor registers the subcommand executor with the registry
func RegisterSubcommandExecutor(registry *Registry) {
	executor := NewSubcommandExecutor(registry.Config(), registry)

	// Set the subcommand executor function
	registry.subcommandExecutor = executor.Execute
}
