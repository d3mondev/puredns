package ctx

import (
	"os"

	"github.com/d3mondev/puredns/v2/internal/app"
)

// Ctx is the program context. It contains the necessary parameters for a command to run.
type Ctx struct {
	ProgramName    string
	ProgramVersion string
	ProgramTagline string
	GitBranch      string
	GitRevision    string

	Options *GlobalOptions
	Stdin   *os.File
}

// NewCtx creates a new context.
func NewCtx() *Ctx {
	return &Ctx{
		ProgramName:    app.AppName,
		ProgramVersion: app.AppVersion,
		ProgramTagline: app.AppDesc,

		GitBranch:   app.GitBranch,
		GitRevision: app.GitRevision,

		Options: DefaultGlobalOptions(),
	}
}
