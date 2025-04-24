package command

import (
	"log/slog"

	"github.com/oscarrieken/master-mold/pkg/config"
)

// Handler defines the interface for command handlers
type Handler interface {
	// Execute executes the command with the given arguments
	Execute(args []string) error
}

// HandlerFunc is a function type that implements the Handler interface
type HandlerFunc func(args []string) error

// Execute calls the handler function
func (f HandlerFunc) Execute(args []string) error {
	return f(args)
}

// Registry is a registry of command handlers
type Registry struct {
	handlers          map[string]Handler
	config            *config.Config
	logger            *slog.Logger
	subcommandExecutor func(name string, args []string) error
}

// NewRegistry creates a new command registry
func NewRegistry(config *config.Config, logger *slog.Logger) *Registry {
	return &Registry{
		handlers: make(map[string]Handler),
		config:   config,
		logger:   logger,
	}
}

// Register registers a command handler
func (r *Registry) Register(name string, handler Handler) {
	r.handlers[name] = handler
}

// RegisterFunc registers a function as a command handler
func (r *Registry) RegisterFunc(name string, fn func(args []string) error) {
	r.Register(name, HandlerFunc(fn))
}

// Get returns the handler for the given command
func (r *Registry) Get(name string) (Handler, bool) {
	handler, ok := r.handlers[name]
	return handler, ok
}

// Execute executes the given command with the given arguments
func (r *Registry) Execute(name string, args []string) error {
	handler, ok := r.Get(name)
	if !ok {
		// If the command is not found in the registry, try to execute it as a subcommand
		if r.subcommandExecutor != nil {
			return r.subcommandExecutor(name, args)
		}
		return r.ExecuteSubcommand(name, args)
	}

	r.logger.Info("Executing command", "command", name)
	return handler.Execute(args)
}

// ExecuteSubcommand executes a subcommand
func (r *Registry) ExecuteSubcommand(name string, args []string) error {
	// This will be implemented in a separate file
	return nil
}

// Config returns the configuration
func (r *Registry) Config() *config.Config {
	return r.config
}

// Logger returns the logger
func (r *Registry) Logger() *slog.Logger {
	return r.logger
}
