package fileoperation

import (
	"errors"
	"testing"

	"github.com/d3mondev/puredns/v2/pkg/filetest"
	"github.com/stretchr/testify/assert"
)

func TestWriteLines(t *testing.T) {
	tests := []struct {
		name       string
		haveLines  []string
		haveWriter *filetest.StubWriter
		haveBuffer int
		wantBuffer []byte
		wantErr    bool
	}{
		{name: "single line", haveLines: []string{"foo"}, haveWriter: filetest.NewStubWriter(nil), haveBuffer: 1024, wantBuffer: []byte("foo\n"), wantErr: false},
		{name: "multiple lines", haveLines: []string{"foo", "bar"}, haveWriter: filetest.NewStubWriter(nil), haveBuffer: 1024, wantBuffer: []byte("foo\nbar\n"), wantErr: false},
		{name: "no lines", haveLines: []string{}, haveWriter: filetest.NewStubWriter(nil), haveBuffer: 1024, wantBuffer: nil, wantErr: false},
		{name: "write error handling", haveLines: []string{"foo"}, haveWriter: filetest.NewStubWriter(errors.New("write error")), haveBuffer: 1, wantBuffer: nil, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotErr := writelines(test.haveLines, test.haveWriter, test.haveBuffer)

			assert.Equal(t, test.wantErr, gotErr != nil)
			assert.Equal(t, test.wantBuffer, test.haveWriter.Buffer)
		})
	}
}

func TestWriteLinesFileError(t *testing.T) {
	file := filetest.CreateFile(t, "")

	tests := []struct {
		name         string
		haveFilename string
		haveLines    []string
		wantErr      bool
	}{
		{name: "valid output file", haveFilename: file.Name(), wantErr: false},
		{name: "file error handling", haveFilename: "", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotErr := WriteLines([]string{"foo", "bar"}, test.haveFilename)

			assert.Equal(t, test.wantErr, gotErr != nil)
		})
	}
}
