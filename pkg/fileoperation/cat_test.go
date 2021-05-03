package fileoperation

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/d3mondev/puredns/v2/pkg/filetest"
	"github.com/stretchr/testify/assert"
)

func TestCat(t *testing.T) {
	testFileA := filetest.CreateFile(t, "contentA")
	testFileB := filetest.CreateFile(t, "contentB")

	tests := []struct {
		name          string
		haveFilenames []string
		wantBuffer    string
		wantErr       bool
	}{
		{name: "cat single file", haveFilenames: []string{testFileA.Name()}, wantBuffer: "contentA\n", wantErr: false},
		{name: "cat two files", haveFilenames: []string{testFileA.Name(), testFileB.Name()}, wantBuffer: "contentA\ncontentB\n", wantErr: false},
		{name: "file error handling", haveFilenames: []string{"thisfiledoesnotexist.txt"}, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stubWriter := filetest.NewStubWriter(nil)

			err := Cat(test.haveFilenames, stubWriter)

			assert.Equal(t, test.wantErr, err != nil)
			assert.Equal(t, test.wantBuffer, string(stubWriter.Buffer))
		})
	}
}

func TestCatIO(t *testing.T) {
	tests := []struct {
		name       string
		haveReader io.Reader
		haveWriter *filetest.StubWriter
		wantBuffer []byte
		wantErr    bool
	}{
		{name: "ok", haveReader: strings.NewReader("test\nfile\n"), haveWriter: filetest.NewStubWriter(nil), wantBuffer: []byte("test\nfile\n"), wantErr: false},
		{name: "read error handling", haveReader: iotest.ErrReader(errors.New("reader error")), haveWriter: filetest.NewStubWriter(nil), wantBuffer: []byte{}, wantErr: true},
		{name: "write error handling", haveReader: strings.NewReader("test\nfile\n"), haveWriter: filetest.NewStubWriter(errors.New("write error")), wantBuffer: []byte{}, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotErr := CatIO([]io.Reader{test.haveReader}, test.haveWriter)

			assert.ElementsMatch(t, test.wantBuffer, test.haveWriter.Buffer)
			assert.Equal(t, test.wantErr, gotErr != nil)
		})
	}
}
