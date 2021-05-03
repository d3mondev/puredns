package shellexecutor

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

var stubExitCode int

func stubExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)

	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{
		"GO_WANT_HELPER_PROCESS=1",
		fmt.Sprintf("EXIT_CODE=%d", stubExitCode),
	}

	return cmd
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	exitCode, err := strconv.Atoi(os.Getenv("EXIT_CODE"))

	if err != nil {
		t.Fatal(err)
	}

	os.Exit(exitCode)
}

func TestShell(t *testing.T) {
	tests := []struct {
		name         string
		haveExitCode int
		wantErr      bool
	}{
		{name: "success", haveExitCode: 0, wantErr: false},
		{name: "exit code 1", haveExitCode: 1, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stubExitCode = test.haveExitCode

			executor := NewShellExecutor()
			executor.execCommand = stubExecCommand

			gotErr := executor.Shell("dummy")

			assert.Equal(t, test.wantErr, gotErr != nil)
		})
	}
}
