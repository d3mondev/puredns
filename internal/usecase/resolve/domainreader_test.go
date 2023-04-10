package resolve

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
)

func TestNewDomainReader(t *testing.T) {
	r := NewDomainReader(io.NopCloser(strings.NewReader("test")), nil, nil)
	assert.NotNil(t, r)
}

func TestDomainReaderRead(t *testing.T) {
	tests := []struct {
		name          string
		haveData      string
		haveDomains   []string
		haveSanitizer DomainSanitizer
		want          string
		wantErr       error
	}{
		{name: "domain list", haveData: "example.com\nwww.example.com\nftp.example.com", want: "example.com\nwww.example.com\nftp.example.com\n", wantErr: io.EOF},
		{name: "words", haveData: "www\nftp\nmail", haveDomains: []string{"example.com"}, want: "www.example.com\nftp.example.com\nmail.example.com\n", wantErr: io.EOF},
		{name: "wildcard", haveData: "www\nftp\nmail", haveDomains: []string{"www.*.example.com"}, want: "www.www.example.com\nwww.ftp.example.com\nwww.mail.example.com\n", wantErr: io.EOF},
		{name: "multiple wildcards", haveData: "word", haveDomains: []string{"www.*.*.example.com"}, want: "www.word.word.example.com\n", wantErr: io.EOF},
		{name: "words multiple domains", haveData: "www\nftp\nmail", haveDomains: []string{"example.com", "example.org"}, want: "www.example.com\nwww.example.org\nftp.example.com\nftp.example.org\nmail.example.com\nmail.example.org\n", wantErr: io.EOF},
		{name: "sanitize", haveData: "_", haveDomains: []string{"example.com"}, haveSanitizer: DefaultSanitizer, want: "\n", wantErr: io.EOF},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := NewDomainReader(io.NopCloser(strings.NewReader(test.haveData)), test.haveDomains, test.haveSanitizer)

			buf := make([]byte, 1024)
			n, err := r.Read(buf)

			assert.ErrorIs(t, err, test.wantErr)
			assert.Equal(t, test.want, string(buf[:n]))
		})
	}
}

func TestDomainReaderRead_ScannerError(t *testing.T) {
	wantErr := errors.New("error")

	r := NewDomainReader(io.NopCloser(iotest.ErrReader(wantErr)), nil, nil)
	buf := make([]byte, 1024)
	_, err := r.Read(buf)

	assert.ErrorIs(t, err, wantErr)
}
