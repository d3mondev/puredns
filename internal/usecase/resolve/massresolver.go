package resolve

import (
	"io"

	"github.com/d3mondev/puredns/v2/internal/pkg/console"
	"github.com/d3mondev/puredns/v2/pkg/massdns"
	"github.com/d3mondev/puredns/v2/pkg/progressbar"
)

// DefaultMassResolver implements the MassResolver interface.
type DefaultMassResolver struct {
	massdns *massdns.Resolver
}

// NewDefaultMassResolver creates a new DefaultMassResolver.
func NewDefaultMassResolver(binPath string) *DefaultMassResolver {
	return &DefaultMassResolver{
		massdns: massdns.NewResolver(binPath),
	}
}

// Resolve calls massdns to resolve the domains contained in the input file.
func (m *DefaultMassResolver) Resolve(r io.Reader, output string, total int, resolversFilename string, qps int) error {
	var template string

	if total == 0 {
		template = "Processed: {{ current }} Rate: {{ rate }} Elapsed: {{ time }}"
	} else {
		template = "[ETA {{ eta }}] {{ bar }} {{ current }}/{{ total }} rate: {{ rate }} qps (time: {{ time }})"
	}

	bar := progressbar.New(
		m.updateProgressBar,
		int64(total),
		progressbar.WithTemplate(template),
		progressbar.WithWriter(console.Output),
	)

	bar.Start()
	err := m.massdns.Resolve(r, output, resolversFilename, qps)
	bar.Stop()

	return err
}

// updateProgressBar is the progress bar update callback.
func (m *DefaultMassResolver) updateProgressBar(bar *progressbar.ProgressBar) {
	current := m.massdns.Current()
	bar.SetCurrent(int64(current))
}
