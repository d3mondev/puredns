package fileoperation

import (
	"errors"
	"testing"

	"github.com/d3mondev/puredns/v2/pkg/filetest"
	"github.com/stretchr/testify/assert"
)

func TestAppendWord(t *testing.T) {
	srcFile := filetest.CreateFile(t, "A\nB\nC\n")
	destFile := filetest.CreateFile(t, "")

	tests := []struct {
		name       string
		haveSource string
		haveDest   string
		wantErr    bool
	}{
		{name: "ok", haveSource: srcFile.Name(), haveDest: destFile.Name(), wantErr: false},
		{name: "source error handling", haveSource: "thisfiledoesnotexist.txt", haveDest: destFile.Name(), wantErr: true},
		{name: "dest error handling", haveSource: srcFile.Name(), haveDest: "", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotErr := AppendWord(test.haveSource, test.haveDest, ":", "word")

			assert.Equal(t, test.wantErr, gotErr != nil)
		})
	}
}

func TestAppendWordIO(t *testing.T) {
	tests := []struct {
		name       string
		haveReader *filetest.StubReader
		haveWriter *filetest.StubWriter
		haveSep    string
		haveWord   string
		wantBuffer []byte
		wantErr    bool
	}{
		{name: "ok", haveReader: filetest.NewStubReader([]byte("A\nB\nC\n"), nil), haveWriter: filetest.NewStubWriter(nil), haveSep: ":", haveWord: "word", wantBuffer: []byte("A:word\nB:word\nC:word\n"), wantErr: false},
		{name: "reader error handling", haveReader: filetest.NewStubReader(nil, errors.New("read error")), haveWriter: filetest.NewStubWriter(nil), haveSep: ":", haveWord: "word", wantErr: true},
		{name: "writer error handling", haveReader: filetest.NewStubReader([]byte("A\nB\nC\n"), nil), haveWriter: filetest.NewStubWriter(errors.New("write error")), haveSep: ":", haveWord: "word", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotErr := appendWordIO(test.haveReader, test.haveWriter, test.haveSep, test.haveWord)

			assert.Equal(t, test.wantErr, gotErr != nil)

			if gotErr == nil {
				assert.ElementsMatch(t, test.wantBuffer, test.haveWriter.Buffer)
			}
		})
	}
}
