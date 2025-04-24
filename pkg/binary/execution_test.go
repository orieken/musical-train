package binary

import (
	"os"
	"path/filepath"
	"testing"

	"log/slog"
)

func TestFindExecutable(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

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

	// Add the temp directory to PATH temporarily
	oldPath := os.Getenv("PATH")
	os.Setenv("PATH", tempDir+string(os.PathListSeparator)+oldPath)
	defer os.Setenv("PATH", oldPath)

	tests := []struct {
		name      string
		command   string
		baseDir   string
		wantError bool
	}{
		{
			name:      "find mm-test1 in PATH",
			command:   "test1",
			baseDir:   "/non-existent-dir",
			wantError: false,
		},
		{
			name:      "find master-mold-test3 in PATH",
			command:   "test3",
			baseDir:   "/non-existent-dir",
			wantError: false,
		},
		{
			name:      "find mm-test1 in baseDir",
			command:   "test1",
			baseDir:   tempDir,
			wantError: false,
		},
		{
			name:      "non-existent command",
			command:   "non-existent",
			baseDir:   tempDir,
			wantError: true,
		},
		{
			name:      "non-executable command",
			command:   "test2",
			baseDir:   tempDir,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmdPath, err := FindExecutable(tt.command, tt.baseDir)
			if (err != nil) != tt.wantError {
				t.Errorf("FindExecutable() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError {
				if cmdPath == "" {
					t.Errorf("FindExecutable() returned empty path")
				}
				// Check that the path exists and is executable
				if !IsExecutable(cmdPath) {
					t.Errorf("FindExecutable() returned non-executable path: %s", cmdPath)
				}
			}
		})
	}
}

func TestExecute(t *testing.T) {
	// This test is more of an integration test and would require a real executable
	// For unit testing, we'll just check that the function doesn't panic
	
	// Create a logger that discards output
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))
	
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// On Unix-like systems, we can use /bin/echo as a test executable
	// On Windows, we would need a different approach
	echoPath := "/bin/echo"
	if _, err := os.Stat(echoPath); os.IsNotExist(err) {
		t.Skip("Skipping test on non-Unix platform")
	}
	
	// Test that Execute doesn't panic
	err = Execute(echoPath, []string{"test"}, logger)
	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}
}