package fileoperation

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileExists(t *testing.T) {
	file, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatal()
	}
	defer func() { file.Close(); os.Remove(file.Name()) }()

	tests := []struct {
		name         string
		haveFilename string
		want         bool
	}{
		{name: "existing file", haveFilename: file.Name(), want: true},
		{name: "non-existing file", haveFilename: "thisfiledoesnotexist.txt", want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := FileExists(test.haveFilename)

			assert.Equal(t, test.want, got)
		})
	}
}
