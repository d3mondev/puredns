package fileoperation

import (
	"os"
	"syscall"
	"testing"

	"github.com/d3mondev/puredns/v2/pkg/filetest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppendLines(t *testing.T) {
	tests := []struct {
		name            string
		haveFileContent string
		haveLines       []string
		wantContent     []string
		wantErr         bool
	}{
		{name: "empty file", haveLines: []string{"foo", "bar"}, wantContent: []string{"foo", "bar"}},
		{name: "file with content", haveFileContent: "one\ntwo\n", haveLines: []string{"foo", "bar"}, wantContent: []string{"one", "two", "foo", "bar"}},
		{name: "no lines", haveFileContent: "one\ntwo\n", wantContent: []string{"one", "two"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dir := filetest.CreateDir(t)
			filename := dir + "/testfile.txt"

			if test.haveFileContent != "" {
				file, err := os.Create(filename)
				require.Nil(t, err)
				_, err = file.WriteString(test.haveFileContent)
				require.Nil(t, err)
				require.Nil(t, file.Sync())
				require.Nil(t, file.Close())
			}

			gotErr := AppendLines(test.haveLines, filename)
			gotContent := filetest.ReadFile(t, filename)

			assert.Equal(t, test.wantErr, gotErr != nil)
			assert.Equal(t, test.wantContent, gotContent)
		})
	}
}

func TestAppendLines_OpenError(t *testing.T) {
	dir := filetest.CreateDir(t)
	err := AppendLines([]string{}, dir)
	assert.ErrorIs(t, err, syscall.Errno(21))
}
