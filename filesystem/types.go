package filesystem

// FileInfo enthält Metadaten zu einer Datei.
type FileInfo struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	IsDir   bool   `json:"is_dir"`
	ModTime int64  `json:"mod_time"`
}

// GetFileInfoArgs sind die Argumente für das get_file_info-Tool.
type GetFileInfoArgs struct {
	Path string `json:"path"`
}

// GetFileInfoResult ist das Ergebnis des get_file_info-Tools.
type GetFileInfoResult struct {
	Info FileInfo `json:"info"`
}
