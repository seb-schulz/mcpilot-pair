package make

// RunMakeArgs are the arguments for the run_make tool.
type RunMakeArgs struct {
	Target    string `json:"target"`    // e.g., "all", "build", "test"
	Directory string `json:"directory"` // optional, relative path; defaults to root
}

// RunMakeResult is the result of the run_make tool.
type RunMakeResult struct {
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	Success  bool   `json:"success"`
	ExitCode int    `json:"exit_code"` // Exit code of the command
}
