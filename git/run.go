package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// RunError is returned when a git command exits with a non-zero status.
type RunError struct {
	Args   []string
	Dir    string
	Stdout string
	Stderr string
	Code   int
}

func (e *RunError) Error() string {
	return fmt.Sprintf("git %s failed (exit %d):\n%s\n%s",
		strings.Join(e.Args, " "), e.Code, e.Stdout, e.Stderr)
}

// run executes a git command in the given directory and returns trimmed stdout.
// Returns a *RunError if the command exits non-zero.
func run(dir string, args []string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", &RunError{
				Args:   args,
				Dir:    dir,
				Stdout: strings.TrimSpace(string(out)),
				Stderr: strings.TrimSpace(string(exitErr.Stderr)),
				Code:   exitErr.ExitCode(),
			}
		}
		return "", fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
	}
	return strings.TrimRight(string(out), "\r\n"), nil
}

// runAllowFail runs a git command and returns stdout, or "" if it fails.
func runAllowFail(dir string, args []string) string {
	out, _ := run(dir, args)
	return out
}
