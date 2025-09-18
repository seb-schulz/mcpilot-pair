package filesystem

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
)

// ReadFileArgs sind die Argumente für das read_file-Tool.
type ReadFileArgs struct {
	Path string `json:"path"`
}

// ReadFileResult ist das Ergebnis des read_file-Tools.
type ReadFileResult struct {
	Content string `json:"content"`
}

// ReadFile liest den Inhalt einer Datei.
func ReadFile(ctx context.Context, args ReadFileArgs) (ReadFileResult, error) {
	content, err := ioutil.ReadFile(args.Path)
	if err != nil {
		return ReadFileResult{}, err
	}
	return ReadFileResult{Content: string(content)}, nil
}

// WriteFileArgs sind die Argumente für das write_file-Tool.
type WriteFileArgs struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

// WriteFileResult ist das Ergebnis des write_file-Tools.
type WriteFileResult struct {
	Success bool `json:"success"`
}

// WriteFile schreibt Inhalt in eine Datei.
func WriteFile(ctx context.Context, args WriteFileArgs) (WriteFileResult, error) {
	if err := ioutil.WriteFile(args.Path, []byte(args.Content), 0644); err != nil {
		return WriteFileResult{}, err
	}
	return WriteFileResult{Success: true}, nil
}

// ListFilesArgs sind die Argumente für das list_files-Tool.
type ListFilesArgs struct {
	Path      string `json:"path"`
	Recursive bool   `json:"recursive,omitempty"`
}

// ListFilesResult ist das Ergebnis des list_files-Tools.
type ListFilesResult struct {
	Files []string `json:"files"`
}

// ListFiles listet Dateien und Verzeichnisse in einem Pfad auf.
func ListFiles(ctx context.Context, args ListFilesArgs) (ListFilesResult, error) {
	var files []string
	if args.Recursive {
		err := filepath.Walk(args.Path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			files = append(files, path)
			return nil
		})
		if err != nil {
			return ListFilesResult{}, err
		}
	} else {
		fileInfos, err := ioutil.ReadDir(args.Path)
		if err != nil {
			return ListFilesResult{}, err
		}
		for _, info := range fileInfos {
			files = append(files, info.Name())
		}
	}
	return ListFilesResult{Files: files}, nil
}

// FileExistsArgs sind die Argumente für das file_exists-Tool.
type FileExistsArgs struct {
	Path string `json:"path"`
}

// FileExistsResult ist das Ergebnis des file_exists-Tools.
type FileExistsResult struct {
	Exists bool `json:"exists"`
}

// FileExists prüft, ob eine Datei oder ein Verzeichnis existiert.
func FileExists(ctx context.Context, args FileExistsArgs) (FileExistsResult, error) {
	_, err := os.Stat(args.Path)
	if os.IsNotExist(err) {
		return FileExistsResult{Exists: false}, nil
	}
	if err != nil {
		return FileExistsResult{}, err
	}
	return FileExistsResult{Exists: true}, nil
}
