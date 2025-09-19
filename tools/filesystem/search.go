package filesystem

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// SearchArgs are the arguments for the search tool.

type SearchArgs struct {
	Query string `json:"query" jsonschema:"the regex pattern to search for in files"`
}

// SearchResult is the result of the search tool.
type SearchResult struct {
	Matches map[string][]Match `json:"matches" jsonschema:"a map of file paths to a list of matches, each containing the line number and the matching line"`
}

// Match represents a single match in a file, including the line number and the matching line.
type Match struct {
	LineNumber int    `json:"line_number"`
	Line       string `json:"line"`
}

// Search searches for a regex pattern in files within the working directory.
func Search(ctx context.Context, args SearchArgs) (SearchResult, error) {
	re, err := regexp.Compile(args.Query)
	if err != nil {
		return SearchResult{}, fmt.Errorf("invalid regex pattern: %v", err)
	}
	wd, err := os.Getwd()
	if err != nil {
		return SearchResult{}, fmt.Errorf("could not get working directory: %v", err)
	}

	result := SearchResult{
		Matches: make(map[string][]Match),
	}
	err = filepath.Walk(wd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || strings.HasPrefix(info.Name(), ".") {
			return nil
		}
		safePath, err := getSafePath(path)
		if err != nil {
			return nil // Skip files outside the working directory
		}

		file, err := os.Open(safePath)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		lineNumber := 1
		for scanner.Scan() {
			line := scanner.Text()
			if re.MatchString(line) {
				result.Matches[safePath] = append(result.Matches[safePath], Match{
					LineNumber: lineNumber,
					Line:       line,
				})
			}
			lineNumber++
		}

		if err := scanner.Err(); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return SearchResult{}, fmt.Errorf("failed to search files: %v", err)
	}

	return result, nil
}
