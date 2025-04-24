package command

import (
	"os"
	"path/filepath"
	"testing"

	"log/slog"

	"github.com/oscarrieken/master-mold/pkg/config"
)

func TestListBinariesHandler_Execute(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test config
	cfg := &config.Config{
		BaseDir: tempDir,
	}

	// Create some test files
	files := []struct {
		name       string
		executable bool
	}{
		{name: "mm-test1", executable: true},
		{name: "mm-test2", executable: false},
		{name: "master-mold-test3", executable: true},
		{name: "test4", executable: true},
	}

	for _, f := range files {
		path := filepath.Join(tempDir, f.name)
		mode := os.FileMode(0644)
		if f.executable {
			mode = 0755
		}
		if err := os.WriteFile(path, []byte("test"), mode); err != nil {
			t.Fatalf("Failed to create file %s: %v", f.name, err)
		}
	}

	// Create the handler
	handler := NewListBinariesHandler(cfg)

	// Redirect stdout to discard output
	oldStdout := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w

	// Execute the handler
	err = handler.Execute(nil)

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Check for errors
	if err != nil {
		t.Errorf("Execute() returned error = %v", err)
	}

	// We don't need to check the output content here, as we've already tested
	// the display functions separately. We just need to make sure the function
	// completes without errors.
}

func TestRegisterListBinariesCommand(t *testing.T) {
	// Create a registry
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cfg := &config.Config{}
	registry := NewRegistry(cfg, logger)

	// Register the list-binaries command
	RegisterListBinariesCommand(registry)

	// Check that the command was registered
	handler, ok := registry.Get("list-binaries")
	if !ok {
		t.Errorf("RegisterListBinariesCommand() did not register the command")
	}
	if handler == nil {
		t.Errorf("RegisterListBinariesCommand() registered a nil handler")
	}

	// Check that the handler is of the correct type
	_, ok = handler.(*ListBinariesHandler)
	if !ok {
		t.Errorf("RegisterListBinariesCommand() registered handler of type %T, want *ListBinariesHandler", handler)
	}
}
