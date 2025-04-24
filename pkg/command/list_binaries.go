package command

import (
	"github.com/pkg/errors"
	"github.com/oscarrieken/master-mold/pkg/binary"
	"github.com/oscarrieken/master-mold/pkg/config"
	"github.com/oscarrieken/master-mold/pkg/display"
)

// ListBinariesHandler handles the list-binaries command
type ListBinariesHandler struct {
	config *config.Config
}

// NewListBinariesHandler creates a new list-binaries command handler
func NewListBinariesHandler(config *config.Config) *ListBinariesHandler {
	return &ListBinariesHandler{
		config: config,
	}
}

// Execute executes the list-binaries command
func (h *ListBinariesHandler) Execute(args []string) error {
	// Ensure the base directory exists
	if err := config.EnsureBaseDirExists(h.config); err != nil {
		return errors.Wrap(err, "failed to ensure base directory exists")
	}

	// Find all binaries
	baseDir := config.GetExpandedBaseDir(h.config)
	binaryPaths, err := binary.FindAll(baseDir)
	if err != nil {
		return errors.Wrap(err, "failed to find binaries")
	}

	// Display the binaries
	display.PrintBinaryPaths(binaryPaths)
	
	return nil
}

// RegisterListBinariesCommand registers the list-binaries command
func RegisterListBinariesCommand(registry *Registry) {
	registry.Register("list-binaries", NewListBinariesHandler(registry.Config()))
}