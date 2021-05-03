package shellexecutor

import "os/exec"

// ShellExecutor is a shell executor object.
type ShellExecutor struct {
	execCommand func(name string, arg ...string) *exec.Cmd
}

// NewShellExecutor returns a new ShellExecutor object.
func NewShellExecutor() *ShellExecutor {
	return &ShellExecutor{
		execCommand: exec.Command,
	}
}

// Shell executes a program with the specified arguments.
// The execution is silent, and an error is returned if the execution ends with an error code.
func (e ShellExecutor) Shell(name string, arg ...string) error {
	cmd := e.execCommand(name, arg...)

	return cmd.Run()
}
