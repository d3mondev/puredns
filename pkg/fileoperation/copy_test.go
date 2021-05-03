package fileoperation

import (
	"errors"
	"testing"

	"github.com/d3mondev/puredns/v2/pkg/filetest"
	"github.com/stretchr/testify/assert"
)

func TestCopy(t *testing.T) {
	srcFile := filetest.CreateFile(t, "file content")
	destFile := filetest.CreateFile(t, "")

	tests := []struct {
		name     string
		haveSrc  string
		haveDest string
		wantErr  bool
	}{
		{name: "ok", haveSrc: srcFile.Name(), haveDest: destFile.Name(), wantErr: false},
		{name: "source file error handling", haveSrc: "thisfiledoesnotexist.txt", haveDest: destFile.Name(), wantErr: true},
		{name: "dest file error handling", haveSrc: srcFile.Name(), haveDest: "", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotErr := Copy(test.haveSrc, test.haveDest)

			assert.Equal(t, test.wantErr, gotErr != nil)
		})
	}
}

func TestCopyByBuffer(t *testing.T) {
	data := "this is test data"

	tests := []struct {
		name       string
		haveReader *filetest.StubReader
		haveWriter *filetest.StubWriter
		haveBuffer int
		wantErr    bool
	}{
		{name: "ok", haveReader: filetest.NewStubReader([]byte(data), nil), haveWriter: filetest.NewStubWriter(nil), haveBuffer: 32768, wantErr: false},
		{name: "small buffer", haveReader: filetest.NewStubReader([]byte(data), nil), haveWriter: filetest.NewStubWriter(nil), haveBuffer: 1, wantErr: false},
		{name: "read error handling", haveReader: filetest.NewStubReader(nil, errors.New("read error")), haveWriter: filetest.NewStubWriter(nil), haveBuffer: 32768, wantErr: true},
		{name: "write error handling", haveReader: filetest.NewStubReader([]byte(data), nil), haveWriter: filetest.NewStubWriter(errors.New("write error")), haveBuffer: 1, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotErr := copyByBuffer(test.haveReader, test.haveWriter, test.haveBuffer)

			assert.Equal(t, test.wantErr, gotErr != nil)

			if gotErr == nil {
				assert.ElementsMatch(t, test.haveReader.Buffer, test.haveWriter.Buffer)
			}
		})
	}
}
