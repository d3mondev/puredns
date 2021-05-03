package resolve

import (
	"fmt"

	"github.com/d3mondev/puredns/v2/internal/app/ctx"
)

// Executor is a simple interface to execute shell commands.
type Executor interface {
	Shell(name string, arg ...string) error
}

// DefaultRequirementChecker checks that the required binaries are present.
type DefaultRequirementChecker struct {
	executor Executor
}

// NewDefaultRequirementChecker returns a new checker object used to validate whether the required binaries can be run.
func NewDefaultRequirementChecker(executor Executor) DefaultRequirementChecker {
	return DefaultRequirementChecker{executor: executor}
}

// Check makes sure that massdns can be executed on the system.
// If not, it displays a message to help the user fix the issue.
func (c DefaultRequirementChecker) Check(opt *ctx.ResolveOptions) error {
	if err := c.executor.Shell(opt.BinPath, "--help"); err != nil {
		fmt.Printf("Unable to execute massdns. Make sure it is present and that the\n")
		fmt.Printf("path to the binary is added to the PATH environment variable.\n\n")

		fmt.Printf("Alternatively, specify the path to massdns using --bin\n\n")

		return fmt.Errorf("unable to execute massdns: %w", err)
	}

	return nil
}
