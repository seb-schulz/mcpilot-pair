package make

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

var allowedTargets = map[string]bool{
	"all":   true,
	"build": true,
	"test":  true,
	"clean": true,
}

// RunMake executes `make -C <directory> <target>` in a bash environment.
func RunMake(ctx context.Context, args RunMakeArgs) (RunMakeResult, error) {
	// Validate target
	if !allowedTargets[args.Target] {
		return RunMakeResult{}, fmt.Errorf("target '%s' is not allowed", args.Target)
	}

	// Build command
	cmd := exec.Command("bash", "-c", "make -C "+args.Directory+" "+args.Target)

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	// Run command
	err := cmd.Run()

	// Capture exit code
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			return RunMakeResult{}, fmt.Errorf("command failed: %v", err)
		}
	}

	return RunMakeResult{
		Stdout:   stdoutBuf.String(),
		Stderr:   stderrBuf.String(),
		Success:  err == nil,
		ExitCode: exitCode,
	}, nil
}
