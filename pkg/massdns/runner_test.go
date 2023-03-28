package massdns

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var stubMassdnsExitCode int

func stubExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)

	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{
		"GO_WANT_HELPER_PROCESS=1",
		fmt.Sprintf("MASSDNS_EXIT_CODE=%d", stubMassdnsExitCode),
	}

	return cmd
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	if os.Args[3] == "massdns" {
		fmt.Fprintf(os.Stderr, "massdns: %v\n", os.Args)
		exitCode, _ := strconv.Atoi(os.Getenv("MASSDNS_EXIT_CODE"))
		os.Exit(exitCode)
	}

	fmt.Fprintf(os.Stderr, "%v\n", os.Args)
}

func TestDefaultRunnerRun(t *testing.T) {
	tests := []struct {
		name                 string
		haveMassdnsExitCode  int
		haveCommandAutoStart bool
		wantErr              bool
	}{
		{name: "success"},
		{name: "massdns error handling", haveMassdnsExitCode: 1, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stubMassdnsExitCode = test.haveMassdnsExitCode

			runner := newDefaultRunner("massdns")
			runner.execCommand = stubExecCommand

			gotErr := runner.Run(strings.NewReader(""), "", "", 10)

			assert.Equal(t, test.wantErr, gotErr != nil, gotErr)
		})
	}
}

func TestCreateMassdnsArgs_DefaultRateLimit(t *testing.T) {
	runner := defaultRunner{}
	gotArgs := runner.createMassdnsArgs("output.txt", "resolvers.txt", 0)
	assert.ElementsMatch(t, []string{"-q", "-r", "resolvers.txt", "-o", "Snl", "-t", "A", "--root", "--retry", "REFUSED", "--retry", "SERVFAIL", "-w", "output.txt"}, gotArgs)
}

func TestCreateMassdnsArgs_CustomRateLimit(t *testing.T) {
	runner := defaultRunner{}
	gotArgs := runner.createMassdnsArgs("output.txt", "resolvers.txt", 100)
	assert.ElementsMatch(t, []string{"-q", "-r", "resolvers.txt", "-o", "Snl", "-t", "A", "--root", "--retry", "REFUSED", "--retry", "SERVFAIL", "-w", "output.txt", "-s", "100"}, gotArgs)
}
