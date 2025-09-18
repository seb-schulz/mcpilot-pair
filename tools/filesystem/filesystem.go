package filesystem

import (
	"context"
	"os"
	"path/filepath"
	"strings"
)

// getSafePath normalizes the path and ensures it is within the working directory.
func getSafePath(path string) (string, error) {
	// Get the absolute working directory
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Normalize the input path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	// Ensure the path is within the working directory
	if !strings.HasPrefix(absPath, wd) {
		return "", os.ErrPermission
	}

	return absPath, nil
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
		return ReadFileResult{}, os.ErrPermission
	}

	content, err := os.ReadFile(safePath)
	if err != nil {
		return ReadFileResult{}, err
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
		return WriteFileResult{}, os.ErrPermission
	}

	if err := os.WriteFile(safePath, []byte(args.Content), 0644); err != nil {
		return WriteFileResult{}, err
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
		return ListFilesResult{}, os.ErrPermission
	}

	var files []string
	if args.Recursive {
		err := filepath.Walk(safePath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// Ignore hidden files/directories (starting with .)
			if strings.HasPrefix(info.Name(), ".") {
				if path == safePath {
					return filepath.SkipDir
				}
				return nil
			}
			files = append(files, path)
			return nil
		})
		if err != nil {
			return ListFilesResult{}, err
		}
	} else {
		fileInfos, err := os.ReadDir(safePath)
		if err != nil {
			return ListFilesResult{}, err
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
		return FileExistsResult{}, os.ErrPermission
	}

	_, err = os.Stat(safePath)
	if os.IsNotExist(err) {
		return FileExistsResult{Exists: false}, nil
	}
	if err != nil {
		return FileExistsResult{}, err
	}
	return FileExistsResult{Exists: true}, nil
}