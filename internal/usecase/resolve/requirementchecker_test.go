package resolve

import (
	"errors"
	"testing"

	"github.com/d3mondev/puredns/v2/internal/app/ctx"
	"github.com/stretchr/testify/assert"
)

type stubExecutor struct {
	index        int
	returnValues []error
}

func (s *stubExecutor) Shell(name string, arg ...string) error {
	ret := s.returnValues[s.index]
	s.index++

	return ret
}

func TestCheck(t *testing.T) {
	wantErr := errors.New("error")

	tests := []struct {
		name      string
		haveError []error
		wantErr   error
	}{
		{name: "ok", haveError: []error{nil, nil}, wantErr: nil},
		{name: "massdns error handling", haveError: []error{wantErr, nil}, wantErr: wantErr},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			executor := &stubExecutor{}
			executor.returnValues = test.haveError
			checker := NewDefaultRequirementChecker(executor)

			err := checker.Check(&ctx.ResolveOptions{})

			assert.ErrorIs(t, err, test.wantErr)
		})
	}
}
