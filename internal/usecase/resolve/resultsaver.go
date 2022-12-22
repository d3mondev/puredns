package resolve

import (
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
		if err := s.fileCopy(workfiles.MassdnsPublic, opt.WriteMassdnsFile); err != nil {
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
