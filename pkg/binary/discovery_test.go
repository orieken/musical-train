package binary

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHasValidPrefix(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     bool
	}{
		{
			name:     "mm prefix",
			filename: "mm-test",
			want:     true,
		},
		{
			name:     "master-mold prefix",
			filename: "master-mold-test",
			want:     true,
		},
		{
			name:     "no valid prefix",
			filename: "test",
			want:     false,
		},
		{
			name:     "empty string",
			filename: "",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasValidPrefix(tt.filename); got != tt.want {
				t.Errorf("HasValidPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsExecutable(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create an executable file
	execFile := filepath.Join(tempDir, "exec-file")
	if err := os.WriteFile(execFile, []byte("test"), 0755); err != nil {
		t.Fatalf("Failed to create executable file: %v", err)
	}

	// Create a non-executable file
	nonExecFile := filepath.Join(tempDir, "non-exec-file")
	if err := os.WriteFile(nonExecFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create non-executable file: %v", err)
	}

	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "executable file",
			path: execFile,
			want: true,
		},
		{
			name: "non-executable file",
			path: nonExecFile,
			want: false,
		},
		{
			name: "non-existent file",
			path: filepath.Join(tempDir, "non-existent"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsExecutable(tt.path); got != tt.want {
				t.Errorf("IsExecutable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractCommandName(t *testing.T) {
	tests := []struct {
		name       string
		binaryPath string
		want       string
	}{
		{
			name:       "mm prefix",
			binaryPath: "/usr/bin/mm-test",
			want:       "test",
		},
		{
			name:       "master-mold prefix",
			binaryPath: "/usr/bin/master-mold-test",
			want:       "test",
		},
		{
			name:       "no valid prefix",
			binaryPath: "/usr/bin/test",
			want:       "test",
		},
		{
			name:       "with directory",
			binaryPath: "/path/to/mm-test",
			want:       "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractCommandName(tt.binaryPath); got != tt.want {
				t.Errorf("ExtractCommandName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindInDirectory(t *testing.T) {
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

	// Test FindInDirectory
	binaries, err := FindInDirectory(tempDir)
	if err != nil {
		t.Fatalf("FindInDirectory() error = %v", err)
	}

	// We should find 2 executable binaries with valid prefixes
	if len(binaries) != 2 {
		t.Errorf("FindInDirectory() found %d binaries, want 2", len(binaries))
	}

	// Check that the binaries have the expected names
	foundTest1 := false
	foundTest3 := false
	for _, binary := range binaries {
		name := filepath.Base(binary)
		if name == "mm-test1" {
			foundTest1 = true
		} else if name == "master-mold-test3" {
			foundTest3 = true
		}
	}

	if !foundTest1 {
		t.Errorf("FindInDirectory() did not find mm-test1")
	}
	if !foundTest3 {
		t.Errorf("FindInDirectory() did not find master-mold-test3")
	}
}