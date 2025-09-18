package filesystem

import (
	"testing"
)

// TestGetSafePath tests the getSafePath function with various scenarios.
func TestGetSafePath(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		expectError bool
	}{
		{
			name:        "Valid path",
			path:        "testfile.txt",
			expectError: false,
		},
		{
			name:        "Path traversal",
			path:        "../testfile.txt",
			expectError: true,
		},
		{
			name:        "Dotfile",
			path:        ".hiddenfile",
			expectError: true,
		},
		{
			name:        "Nested dotfile",
			path:        "dir/.hiddenfile",
			expectError: true,
		},
		{
			name:        "Current directory",
			path:        "./testfile.txt",
			expectError: false,
		},
		{
			name:        "Subdirectory",
			path:        "subdir/testfile.txt",
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := getSafePath(tc.path)
			if tc.expectError && err == nil {
				t.Errorf("Expected error for %s, got nil", tc.path)
			}
			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error for %s: %v", tc.path, err)
			}
		})
	}
}
