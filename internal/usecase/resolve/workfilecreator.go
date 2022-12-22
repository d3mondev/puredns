package resolve

import (
	"fmt"
	"io/ioutil"
	"os"
)

// Workfiles are temporary files used during the program execution.
type Workfiles struct {
	TempDirectory string

	Domains        string
	MassdnsPublic  string
	MassdnsTrusted string
	Temporary      string

	PublicResolvers  string
	TrustedResolvers string

	WildcardRoots string
}

// Close deletes all the temporary files that were created.
func (w *Workfiles) Close() {
	if w.TempDirectory != "" {
		os.RemoveAll(w.TempDirectory)
	}
}

// DefaultWorkfileCreator is a service that creates a set of workfiles on disk.
type DefaultWorkfileCreator struct {
	osMkdirTemp func(dir string, pattern string) (string, error)
	osCreate    func(name string) (*os.File, error)
}

// NewDefaultWorkfileCreator creates a new set of temporary files.
// Call Close() to cleanup the files once they are no longer needed.
func NewDefaultWorkfileCreator() *DefaultWorkfileCreator {
	return &DefaultWorkfileCreator{
		osMkdirTemp: ioutil.TempDir,
		osCreate:    os.Create,
	}
}

// Create creates a new set of workfiles.
func (w *DefaultWorkfileCreator) Create() (*Workfiles, error) {
	files := &Workfiles{}

	dir, err := w.osMkdirTemp("", "puredns.")
	if err != nil {
		return nil, fmt.Errorf("unable to create temporary work directory: %w", err)
	}

	files.TempDirectory = dir

	if files.Domains, err = w.createFile(files.TempDirectory + "/" + "domains.txt"); err != nil {
		return nil, err
	}

	if files.MassdnsPublic, err = w.createFile(files.TempDirectory + "/" + "massdns_public.txt"); err != nil {
		return nil, err
	}

	if files.MassdnsTrusted, err = w.createFile(files.TempDirectory + "/" + "massdns_trusted.txt"); err != nil {
		return nil, err
	}

	if files.Temporary, err = w.createFile(files.TempDirectory + "/" + "temporary.txt"); err != nil {
		return nil, err
	}

	if files.PublicResolvers, err = w.createFile(files.TempDirectory + "/" + "resolvers.txt"); err != nil {
		return nil, err
	}

	if files.TrustedResolvers, err = w.createFile(files.TempDirectory + "/" + "trusted.txt"); err != nil {
		return nil, err
	}

	if files.WildcardRoots, err = w.createFile(files.TempDirectory + "/" + "wildcards.txt"); err != nil {
		return nil, err
	}

	return files, nil
}

func (w *DefaultWorkfileCreator) createFile(filepath string) (string, error) {
	file, err := w.osCreate(filepath)
	if err != nil {
		return "", fmt.Errorf("unable to create temporary file %s: %w", filepath, err)
	}
	defer file.Close()

	return file.Name(), nil
}
