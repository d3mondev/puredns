package fileoperation

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/d3mondev/puredns/v2/pkg/filetest"
	"github.com/stretchr/testify/assert"
)

func TestCountLines(t *testing.T) {
	file, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatal()
	}
	defer func() { file.Close(); os.Remove(file.Name()) }()

	_, err = file.WriteString("line1\nline2\nline3\n")
	if err != nil {
		t.Fatal()
	}

	tests := []struct {
		name         string
		haveFilename string
		wantCount    int
		wantErr      bool
	}{
		{name: "existing file", haveFilename: file.Name(), wantCount: 3},
		{name: "file error handling", haveFilename: "thisfiledoesnotexist.txt", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotCount, gotErr := CountLines(test.haveFilename)

			assert.Equal(t, test.wantCount, gotCount)
			assert.Equal(t, test.wantErr, gotErr != nil)
		})
	}
}

func TestCount(t *testing.T) {
	lines := "line1\nline2\n"
	readerError := errors.New("reader error")

	tests := []struct {
		name       string
		haveReader io.ReadCloser
		wantCount  int
		wantErr    error
	}{
		{name: "success", haveReader: filetest.NewStubReader([]byte(lines), nil), wantCount: 2, wantErr: nil},
		{name: "empty reader", haveReader: filetest.NewStubReader(nil, nil), wantCount: 0, wantErr: nil},
		{name: "reader error handling", haveReader: filetest.NewStubReader(nil, readerError), wantCount: 0, wantErr: readerError},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotCount, gotErr := count(test.haveReader)

			assert.ErrorIs(t, gotErr, test.wantErr)
			assert.Equal(t, test.wantCount, gotCount)
		})
	}
}
