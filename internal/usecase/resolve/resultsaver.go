package resolve

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/d3mondev/puredns/v2/internal/app/ctx"
	"github.com/d3mondev/puredns/v2/pkg/fileoperation"
)

// ResultFileSaver is responsible for saving the results of the resolve operation to files.
type ResultFileSaver struct {
	fileCopy func(src string, dest string) error
}

// NewResultFileSaver creates a new ResultSaver object.
func NewResultFileSaver() *ResultFileSaver {
	return &ResultFileSaver{
		fileCopy: fileoperation.Copy,
	}
}

// Save saves the results contained in the working files according to the specified options.
func (s *ResultFileSaver) Save(workfiles *Workfiles, opt *ctx.ResolveOptions) error {
	if opt.WriteDomainsFile != "" {
		if err := s.fileCopy(workfiles.Domains, opt.WriteDomainsFile); err != nil {
			return err
		}
	}

	if opt.WriteMassdnsFile != "" {
		if err := s.saveMassdns(workfiles, opt); err != nil {
			return err
		}
	}

	if opt.WriteWildcardsFile != "" {
		if err := s.fileCopy(workfiles.WildcardRoots, opt.WriteWildcardsFile); err != nil {
			return err
		}
	}

	return nil
}

// saveMassdns preserves complete answer sections for the final valid domains.
func (s *ResultFileSaver) saveMassdns(workfiles *Workfiles, opt *ctx.ResolveOptions) error {
	domains, err := fileoperation.ReadLines(workfiles.Domains)
	if err != nil {
		return err
	}
	valid := make(map[string]struct{}, len(domains))
	for _, domain := range domains {
		valid[domain] = struct{}{}
	}

	filename := workfiles.MassdnsPublic
	if !opt.SkipValidation && !opt.TrustedOnly {
		filename = workfiles.MassdnsTrusted
	}
	input, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer input.Close()

	output, err := os.Create(opt.WriteMassdnsFile)
	if err != nil {
		return err
	}
	defer output.Close()
	writer := bufio.NewWriter(output)

	scanner := bufio.NewScanner(input)
	firstLine, keep := true, false
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" && firstLine {
			domain := strings.TrimSuffix(strings.SplitN(line, " ", 2)[0], ".")
			_, keep = valid[domain]
			firstLine = false
		}
		if keep {
			if _, err := fmt.Fprintln(writer, line); err != nil {
				return err
			}
		}
		if line == "" {
			firstLine = true
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return writer.Flush()
}
