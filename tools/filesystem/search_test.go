package filesystem

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestSearch(t *testing.T) {
	// Setup: Create a temporary directory with test files
	dir := t.TempDir()
	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}
	defer os.Chdir(oldwd)

	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Failed to change working directory: %v", err)
	}

	// Create test files
	testFiles := []struct {
		name    string
		content string
	}{
		{"file1.txt", "This is line 1.\nThis is line 2.\nThis is line 3.\n"},
		{"file2.txt", "Another line 1.\nAnother line 2.\n"},
		{".hidden.txt", "This file should be ignored.\n"},
	}

	for _, tf := range testFiles {
		if err := os.WriteFile(tf.name, []byte(tf.content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", tf.name, err)
		}
	}

	// Test cases
	tests := []struct {
		name     string
		query    string
		expected map[string][]Match
	}{
		{
			name:  "Search for 'line 2'",
			query: "line 2",
			expected: map[string][]Match{
				filepath.Join(dir, "file1.txt"): {
					{LineNumber: 2, Line: "This is line 2."},
				},
				filepath.Join(dir, "file2.txt"): {
					{LineNumber: 2, Line: "Another line 2."},
				},
			},
		},
		{
			name:     "Search for non-existent pattern",
			query:    "non-existent-pattern",
			expected: map[string][]Match{},
		}, {
			name:     "Search for content in the hidden file",
			query:    "ignored",
			expected: map[string][]Match{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			args := SearchArgs{Query: tc.query}
			result, err := Search(context.Background(), args)
			if err != nil {
				t.Errorf("Search failed: %v", err)
			}

			// Check if the number of files matches
			if len(result.Matches) != len(tc.expected) {
				t.Errorf("Expected %d files, got %d", len(tc.expected), len(result.Matches))
			}

			// Check each file and its matches
			for file, matches := range tc.expected {
				if resultMatches, ok := result.Matches[file]; !ok {
					t.Errorf("Expected file %s not found in result", file)
				} else {
					if len(resultMatches) != len(matches) {
						t.Errorf("Expected %d matches for file %s, got %d", len(matches), file, len(resultMatches))
					} else {
						for i, expectedMatch := range matches {
							if i >= len(resultMatches) {
								t.Errorf("Missing match in file %s", file)
								continue
							}
							resultMatch := resultMatches[i]
							if resultMatch.LineNumber != expectedMatch.LineNumber || resultMatch.Line != expectedMatch.Line {
								t.Errorf("Match mismatch in file %s: expected %+v, got %+v", file, expectedMatch, resultMatch)
							}
						}
					}
				}
			}
		})
	}
}
