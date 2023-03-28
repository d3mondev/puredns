package massdns

import (
	"io"
	"os/exec"
	"strconv"
)

type defaultRunner struct {
	binPath     string
	execCommand func(name string, arg ...string) *exec.Cmd
}

func newDefaultRunner(binPath string) *defaultRunner {
	return &defaultRunner{
		binPath:     binPath,
		execCommand: exec.Command,
	}
}

// Run executes massdns on the specified domains files and saves the results to the output file.
func (runner *defaultRunner) Run(r io.Reader, output string, resolvers string, qps int) error {
	// Create massdns program arguments
	massdnsArgs := runner.createMassdnsArgs(output, resolvers, qps)

	// Create a new exec.Cmd and set Stdin and Stdout to our custom handlers to avoid file operations
	massdns := runner.execCommand(runner.binPath, massdnsArgs...)
	massdns.Stdin = r

	// Run massdns and block until it's done
	if err := massdns.Run(); err != nil {
		return err
	}

	return nil
}

// createMassdnsArgs creates the command line arguments for massdns.
func (runner *defaultRunner) createMassdnsArgs(output string, resolvers string, qps int) []string {
	// Default command line
	args := []string{"-q", "-r", resolvers, "-o", "Snl", "-t", "A", "--root", "--retry", "REFUSED", "--retry", "SERVFAIL", "-w", output}

	// Set the massdns hashmap size manually to prevent it from accumulating DNS query on start
	if qps > 0 {
		args = append(args, "-s", strconv.Itoa(qps))
	}

	return args
}
