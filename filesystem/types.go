package filesystem

// FileInfo contains metadata about a file.
type FileInfo struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	IsDir   bool   `json:"is_dir"`
	ModTime int64  `json:"mod_time"`
}

// GetFileInfoArgs are the arguments for the get_file_info tool.
type GetFileInfoArgs struct {
	Path string `json:"path"`
}

// GetFileInfoResult is the result of the get_file_info tool.
type GetFileInfoResult struct {
	Info FileInfo `json:"info"`
}