package resolve

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type stubFileCreator struct {
	successes        int
	returnsOnFailure error
}

func (s *stubFileCreator) Create(filepath string) (*os.File, error) {
	s.successes--

	if s.successes < 0 {
		return nil, s.returnsOnFailure
	}

	return os.Create(filepath)
}

type stubDirCreator struct {
	returns error
}

func (s *stubDirCreator) MkdirTemp(dir string, pattern string) (string, error) {
	if s.returns != nil {
		return "", s.returns
	}

	return os.MkdirTemp("", "")
}

func TestCreate(t *testing.T) {
	createError := errors.New("create failed")
	mkDirTempError := errors.New("mkdirtemp failed")

	tests := []struct {
		name                string
		haveMkdirTempError  error
		haveCreateSuccesses int
		haveCreateError     error
		wantErr             error
	}{
		{name: "success", haveCreateSuccesses: 100},
		{name: "mkdirtemp error handling", haveMkdirTempError: mkDirTempError, wantErr: mkDirTempError},
		{name: "first create error handling", haveCreateSuccesses: 0, haveCreateError: createError, wantErr: createError},
		{name: "second create error handling", haveCreateSuccesses: 1, haveCreateError: createError, wantErr: createError},
		{name: "third create error handling", haveCreateSuccesses: 2, haveCreateError: createError, wantErr: createError},
		{name: "fourth create error handling", haveCreateSuccesses: 3, haveCreateError: createError, wantErr: createError},
		{name: "fifth create error handling", haveCreateSuccesses: 4, haveCreateError: createError, wantErr: createError},
		{name: "sixth create error handling", haveCreateSuccesses: 5, haveCreateError: createError, wantErr: createError},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			createFile := stubFileCreator{successes: test.haveCreateSuccesses, returnsOnFailure: test.haveCreateError}
			createDir := stubDirCreator{returns: test.haveMkdirTempError}

			creator := NewDefaultWorkfileCreator()
			creator.osCreate = createFile.Create
			creator.osMkdirTemp = createDir.MkdirTemp

			gotFiles, gotErr := creator.Create()

			if gotFiles != nil {
				defer gotFiles.Close()
			}

			assert.ErrorIs(t, gotErr, test.wantErr)
		})
	}
}
