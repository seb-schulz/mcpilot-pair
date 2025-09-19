package make

// RunMakeArgs are the arguments for the run_make tool.
type RunMakeArgs struct {
	Target    string `json:"target" jsonschema:"the make target to execute (e.g., 'all', 'build', 'test')"`
	Directory string `json:"directory,omitempty" jsonschema:"The optional relative path to execute the make command. If omitted, the root directory is used. Only specify if you explicitly want to run make in a subdirectory"`
}

// RunMakeResult is the result of the run_make tool.
type RunMakeResult struct {
	Stdout   string `json:"stdout" jsonschema:"the standard output of the make command"`
	Stderr   string `json:"stderr" jsonschema:"the standard error output of the make command"`
	Success  bool   `json:"success" jsonschema:"indicates whether the make command executed successfully"`
	ExitCode int    `json:"exit_code" jsonschema:"the exit code of the make command"`
}
