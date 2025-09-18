package filesystem

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// getSafePath ensures the path is within the working directory (including symlinks)
// and does not contain dotfiles/dotdirs. It is platform-independent and works on both Unix and Windows.
func getSafePath(p string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("could not get working directory: %v", err)
	}

	p = filepath.Clean(p)
	abs, err := filepath.Abs(p)
	if err != nil {
		return "", fmt.Errorf("could not resolve absolute path: %v", err)
	}

	absEval, err := filepath.EvalSymlinks(abs)
	if err != nil {
		log.Printf("Warning: could not resolve symlinks for %s: %v", abs, err)
		absEval = abs
	}

	wdEval, err := filepath.EvalSymlinks(wd)
	if err != nil {
		log.Printf("Warning: could not resolve working directory symlinks: %v", err)
		wdEval = wd
	}

	rel, err := filepath.Rel(wdEval, absEval)
	if err != nil {
		return "", fmt.Errorf("could not compute relative path: %v", err)
	}

	// No upward traversal allowed
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path traversal attempt detected: %s", rel)
	}

	// Prohibit dotfiles/dotdirs (any segment starting with '.')
	if containsDotSegment(rel) {
		return "", fmt.Errorf("access to dotfiles/dotdirs not allowed: %s", rel)
	}

	return absEval, nil
}

// containsDotSegment checks if the relative path contains any segment starting with '.'.
func containsDotSegment(rel string) bool {
	for _, seg := range strings.Split(rel, string(filepath.Separator)) {
		if seg == "" || seg == "." {
			continue
		}
		if strings.HasPrefix(seg, ".") {
			return true
		}
	}
	return false
}

// ReadFileArgs are the arguments for the read_file tool.
type ReadFileArgs struct {
	Path string `json:"path"`
}

// ReadFileResult is the result of the read_file tool.
type ReadFileResult struct {
	Content string `json:"content"`
}

// ReadFile reads the content of a file within the working directory.
func ReadFile(ctx context.Context, args ReadFileArgs) (ReadFileResult, error) {
	safePath, err := getSafePath(args.Path)
	if err != nil {
		log.Printf("Invalid path: %v", err)
		return ReadFileResult{}, fmt.Errorf("invalid path: %v", err)
	}

	content, err := os.ReadFile(safePath)
	if err != nil {
		return ReadFileResult{}, fmt.Errorf("failed to read file: %v", err)
	}
	return ReadFileResult{Content: string(content)}, nil
}

// WriteFileArgs are the arguments for the write_file tool.
type WriteFileArgs struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

// WriteFileResult is the result of the write_file tool.
type WriteFileResult struct {
	Success bool `json:"success"`
}

// WriteFile writes content to a file within the working directory.
func WriteFile(ctx context.Context, args WriteFileArgs) (WriteFileResult, error) {
	safePath, err := getSafePath(args.Path)
	if err != nil {
		log.Printf("Invalid path: %v", err)
		return WriteFileResult{}, fmt.Errorf("invalid path: %v", err)
	}

	if err := os.WriteFile(safePath, []byte(args.Content), 0644); err != nil {
		return WriteFileResult{}, fmt.Errorf("failed to write file %s: %v", safePath, err)
	}
	return WriteFileResult{Success: true}, nil
}

// ListFilesArgs are the arguments for the list_files tool.
type ListFilesArgs struct {
	Path      string `json:"path"`
	Recursive bool   `json:"recursive,omitempty"`
}

// ListFilesResult is the result of the list_files tool.
type ListFilesResult struct {
	Files []string `json:"files"`
}

// ListFiles lists files and directories within the working directory.
func ListFiles(ctx context.Context, args ListFilesArgs) (ListFilesResult, error) {
	safePath, err := getSafePath(args.Path)
	if err != nil {
		log.Printf("Invalid path: %v", err)
		return ListFilesResult{}, fmt.Errorf("invalid path: %v", err)
	}

	var files []string
	if args.Recursive {
		err := filepath.Walk(safePath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Printf("Error walking path %s: %v", path, err)
				return err
			}
			if !strings.HasPrefix(info.Name(), ".") {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return ListFilesResult{}, fmt.Errorf("failed to list files: %v", err)
		}
	} else {
		fileInfos, err := os.ReadDir(safePath)
		if err != nil {
			return ListFilesResult{}, fmt.Errorf("failed to read directory: %v", err)
		}
		for _, info := range fileInfos {
			if !strings.HasPrefix(info.Name(), ".") {
				files = append(files, info.Name())
			}
		}
	}
	return ListFilesResult{Files: files}, nil
}

// FileExistsArgs are the arguments for the file_exists tool.
type FileExistsArgs struct {
	Path string `json:"path"`
}

// FileExistsResult is the result of the file_exists tool.
type FileExistsResult struct {
	Exists bool `json:"exists"`
}

// FileExists checks if a file or directory exists within the working directory.
func FileExists(ctx context.Context, args FileExistsArgs) (FileExistsResult, error) {
	safePath, err := getSafePath(args.Path)
	if err != nil {
		log.Printf("Invalid path: %v", err)
		return FileExistsResult{}, fmt.Errorf("invalid path: %v", err)
	}

	_, err = os.Stat(safePath)
	if os.IsNotExist(err) {
		return FileExistsResult{Exists: false}, nil
	}
	if err != nil {
		return FileExistsResult{}, fmt.Errorf("failed to check file existence: %v", err)
	}
	return FileExistsResult{Exists: true}, nil
}
