package massdns

import (
	"io"
	"os"
)

// Runner is an interface that runs the commands required to execute the massdns binary.
type Runner interface {
	Run(reader io.Reader, output string, resolversFile string, qps int) error
}

// Resolver uses massdns to resolve a batch of domain names.
type Resolver struct {
	osOpen   func(file string) (*os.File, error)
	osCreate func(file string) (*os.File, error)

	runner Runner

	reader *LineReader
}

// NewResolver creates a new Resolver.
func NewResolver(binPath string) *Resolver {
	return &Resolver{
		runner: newDefaultRunner(binPath),

		osOpen:   os.Open,
		osCreate: os.Create,
	}
}

// Resolve reads domain names from the reader and saves the answers to a file.
// It uses the resolvers and queries per second limit specified.
func (r *Resolver) Resolve(reader io.Reader, output string, resolversFile string, qps int) error {
	r.reader = NewLineReader(reader, qps)

	if err := r.runner.Run(r.reader, output, resolversFile, qps); err != nil {
		return err
	}

	return nil
}

// Current returns the index of the last domain processed.
func (r *Resolver) Current() int {
	if r.reader == nil {
		return 0
	}

	return r.reader.Count()
}
