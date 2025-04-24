package display

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/oscarrieken/master-mold/pkg/binary"
)

func TestFormatBinaryInfo(t *testing.T) {
	tests := []struct {
		name string
		info BinaryInfo
		want string
	}{
		{
			name: "basic info",
			info: BinaryInfo{
				Name:     "test",
				FullPath: "/usr/bin/mm-test",
			},
			want: "  - test (/usr/bin/mm-test)",
		},
		{
			name: "empty name",
			info: BinaryInfo{
				Name:     "",
				FullPath: "/usr/bin/mm-test",
			},
			want: "  -  (/usr/bin/mm-test)",
		},
		{
			name: "empty path",
			info: BinaryInfo{
				Name:     "test",
				FullPath: "",
			},
			want: "  - test ()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatBinaryInfo(tt.info); got != tt.want {
				t.Errorf("FormatBinaryInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcessBinaries(t *testing.T) {
	tests := []struct {
		name        string
		binaryPaths []string
		wantCount   int
		wantNames   []string
	}{
		{
			name: "unique binaries",
			binaryPaths: []string{
				"/usr/bin/mm-test1",
				"/usr/bin/mm-test2",
				"/usr/bin/master-mold-test3",
			},
			wantCount: 3,
			wantNames: []string{"test1", "test2", "test3"},
		},
		{
			name: "duplicate binaries",
			binaryPaths: []string{
				"/usr/bin/mm-test1",
				"/usr/local/bin/mm-test1",
				"/usr/bin/mm-test2",
			},
			wantCount: 2,
			wantNames: []string{"test1", "test2"},
		},
		{
			name:        "empty list",
			binaryPaths: []string{},
			wantCount:   0,
			wantNames:   []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ProcessBinaries(tt.binaryPaths)
			if len(got) != tt.wantCount {
				t.Errorf("ProcessBinaries() returned %d binaries, want %d", len(got), tt.wantCount)
			}

  	// Check that the result contains all the expected names
			for i, wantName := range tt.wantNames {
				if i >= len(got) {
					t.Errorf("ProcessBinaries() missing expected name at index %d: %s", i, wantName)
					continue
				}
				if got[i].Name != wantName {
					t.Errorf("ProcessBinaries()[%d].Name = %v, want %v", i, got[i].Name, wantName)
				}

				// Verify that the binary path is correctly processed
				expectedPath := ""
				for _, path := range tt.binaryPaths {
					if binary.ExtractCommandName(path) == wantName {
						expectedPath = path
						break
					}
				}
				if got[i].FullPath != expectedPath {
					t.Errorf("ProcessBinaries()[%d].FullPath = %v, want %v", i, got[i].FullPath, expectedPath)
				}
			}
		})
	}
}

func TestPrintBinaries(t *testing.T) {
	tests := []struct {
		name     string
		binaries []BinaryInfo
		want     string
	}{
		{
			name: "multiple binaries",
			binaries: []BinaryInfo{
				{Name: "test1", FullPath: "/usr/bin/mm-test1"},
				{Name: "test2", FullPath: "/usr/bin/mm-test2"},
			},
			want: "Available subcommands:\n  - test1 (/usr/bin/mm-test1)\n  - test2 (/usr/bin/mm-test2)\n",
		},
		{
			name:     "no binaries",
			binaries: []BinaryInfo{},
			want:     "No subcommand binaries found.\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Redirect stdout to capture output
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			PrintBinaries(tt.binaries)

			// Restore stdout
			w.Close()
			os.Stdout = oldStdout

			// Read the captured output
			var buf bytes.Buffer
			io.Copy(&buf, r)
			got := buf.String()

			if got != tt.want {
				t.Errorf("PrintBinaries() output = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrintBinaryPaths(t *testing.T) {
	tests := []struct {
		name        string
		binaryPaths []string
		wantContains string
	}{
		{
			name: "multiple binaries",
			binaryPaths: []string{
				"/usr/bin/mm-test1",
				"/usr/bin/mm-test2",
			},
			wantContains: "Available subcommands:",
		},
		{
			name:        "no binaries",
			binaryPaths: []string{},
			wantContains: "No subcommand binaries found.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Redirect stdout to capture output
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			PrintBinaryPaths(tt.binaryPaths)

			// Restore stdout
			w.Close()
			os.Stdout = oldStdout

			// Read the captured output
			var buf bytes.Buffer
			io.Copy(&buf, r)
			got := buf.String()

			if !strings.Contains(got, tt.wantContains) {
				t.Errorf("PrintBinaryPaths() output = %v, does not contain %v", got, tt.wantContains)
			}
		})
	}
}
