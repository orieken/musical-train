package command

import (
	"testing"

	"log/slog"
	"os"

	"github.com/oscarrieken/master-mold/pkg/config"
)

// MockHandler is a mock implementation of the Handler interface for testing
type MockHandler struct {
	ExecuteCalled bool
	Args          []string
	ReturnError   error
}

func (m *MockHandler) Execute(args []string) error {
	m.ExecuteCalled = true
	m.Args = args
	return m.ReturnError
}

func TestRegistry_Register(t *testing.T) {
	// Create a registry
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cfg := &config.Config{}
	registry := NewRegistry(cfg, logger)

	// Create a mock handler
	handler := &MockHandler{}

	// Register the handler
	registry.Register("test", handler)

	// Check that the handler was registered
	if _, ok := registry.handlers["test"]; !ok {
		t.Errorf("Register() did not register the handler")
	}
}

func TestRegistry_RegisterFunc(t *testing.T) {
	// Create a registry
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cfg := &config.Config{}
	registry := NewRegistry(cfg, logger)

	// Create a handler function
	var executeCalled bool
	var executeArgs []string
	handlerFunc := func(args []string) error {
		executeCalled = true
		executeArgs = args
		return nil
	}

	// Register the handler function
	registry.RegisterFunc("test", handlerFunc)

	// Check that the handler was registered
	if _, ok := registry.handlers["test"]; !ok {
		t.Errorf("RegisterFunc() did not register the handler")
	}

	// Execute the handler
	args := []string{"arg1", "arg2"}
	registry.handlers["test"].Execute(args)

	// Check that the handler function was called with the correct arguments
	if !executeCalled {
		t.Errorf("RegisterFunc() handler was not called")
	}
	if len(executeArgs) != len(args) {
		t.Errorf("RegisterFunc() handler was called with %d args, want %d", len(executeArgs), len(args))
	}
	for i, arg := range args {
		if executeArgs[i] != arg {
			t.Errorf("RegisterFunc() handler arg %d = %v, want %v", i, executeArgs[i], arg)
		}
	}
}

func TestRegistry_Get(t *testing.T) {
	// Create a registry
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cfg := &config.Config{}
	registry := NewRegistry(cfg, logger)

	// Create a mock handler
	handler := &MockHandler{}

	// Register the handler
	registry.Register("test", handler)

	// Get the handler
	got, ok := registry.Get("test")
	if !ok {
		t.Errorf("Get() returned ok = false, want true")
	}
	if got != handler {
		t.Errorf("Get() returned handler = %v, want %v", got, handler)
	}

	// Get a non-existent handler
	got, ok = registry.Get("non-existent")
	if ok {
		t.Errorf("Get() returned ok = true for non-existent handler, want false")
	}
	if got != nil {
		t.Errorf("Get() returned handler = %v for non-existent handler, want nil", got)
	}
}

func TestRegistry_Execute(t *testing.T) {
	// Create a registry
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cfg := &config.Config{}
	registry := NewRegistry(cfg, logger)

	// Create a mock handler
	handler := &MockHandler{}

	// Register the handler
	registry.Register("test", handler)

	// Execute the handler
	args := []string{"arg1", "arg2"}
	err := registry.Execute("test", args)
	if err != nil {
		t.Errorf("Execute() returned error = %v, want nil", err)
	}

	// Check that the handler was called with the correct arguments
	if !handler.ExecuteCalled {
		t.Errorf("Execute() did not call the handler")
	}
	if len(handler.Args) != len(args) {
		t.Errorf("Execute() called handler with %d args, want %d", len(handler.Args), len(args))
	}
	for i, arg := range args {
		if handler.Args[i] != arg {
			t.Errorf("Execute() handler arg %d = %v, want %v", i, handler.Args[i], arg)
		}
	}
}