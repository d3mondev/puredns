package massdns

import (
	"syscall"
	"testing"

	"github.com/d3mondev/puredns/v2/pkg/filetest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDefaultWriteCallback(t *testing.T) {
	domainFile := filetest.CreateFile(t, "")
	massdnsFile := filetest.CreateFile(t, "")
	badFile := filetest.CreateDir(t)

	tests := []struct {
		name            string
		haveMassdnsFile string
		haveDomainFile  string
		want            error
	}{
		{name: "ok", haveMassdnsFile: massdnsFile.Name(), haveDomainFile: domainFile.Name()},
		{name: "massdns file error handling", haveMassdnsFile: badFile, want: syscall.Errno(21)},
		{name: "domain file error handling", haveDomainFile: badFile, want: syscall.Errno(21)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewDefaultWriteCallback(test.haveMassdnsFile, test.haveDomainFile)
			assert.ErrorIs(t, err, test.want)
		})
	}
}

func TestDefaultWriteCallback(t *testing.T) {
	massdnsFile := filetest.CreateFile(t, "")
	domainFile := filetest.CreateFile(t, "")

	tests := []struct {
		name        string
		haveLines   []string
		wantMassdns []string
		wantDomain  []string
		wantErr     error
	}{
		{
			name: "single record",
			haveLines: []string{
				"example.com. A 127.0.0.1",
			},
			wantMassdns: []string{
				"example.com A 127.0.0.1",
			},
			wantDomain: []string{
				"example.com",
			},
		},
		{
			name: "multiple record",
			haveLines: []string{
				"www.example.com. CNAME example.com.",
				"example.com. A 127.0.0.1",
				"example.com. AAAA ::1",
			},
			wantMassdns: []string{
				"www.example.com CNAME example.com",
				"www.example.com A 127.0.0.1",
				"www.example.com AAAA ::1",
			},
			wantDomain: []string{
				"www.example.com",
			},
		},
		{
			name: "invalid record type",
			haveLines: []string{
				"example.com. NS ns.example.com.",
			},
			wantMassdns: []string{},
			wantDomain:  []string{},
		},
		{
			name: "save domain after valid record is found",
			haveLines: []string{
				"example.com. NS ns.example.com.",
				"example.com. AAAA ::1",
			},
			wantMassdns: []string{
				"example.com AAAA ::1",
			},
			wantDomain: []string{
				"example.com",
			},
		},
		{
			name: "multiple answer sections",
			haveLines: []string{
				"",
				"example.com. A 127.0.0.1",
				"",
				"",
				"www.test.com. CNAME test.com.",
				"test.com. A 127.0.0.1",
				"test.com. AAAA ::1",
				"",
			},
			wantMassdns: []string{
				"example.com A 127.0.0.1",
				"www.test.com CNAME test.com",
				"www.test.com A 127.0.0.1",
				"www.test.com AAAA ::1",
			},
			wantDomain: []string{
				"example.com",
				"www.test.com",
			},
		},
		{
			name: "skip answer section containing bad data",
			haveLines: []string{
				"garbage",
				"example.com. A 127.0.0.1",
			},
			wantMassdns: []string{},
			wantDomain:  []string{},
		},
		{
			name: "empty domain",
			haveLines: []string{
				". A 127.0.0.1",
			},
			wantMassdns: []string{},
			wantDomain:  []string{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cb, err := NewDefaultWriteCallback(massdnsFile.Name(), domainFile.Name())
			require.Nil(t, err)

			for _, line := range test.haveLines {
				err := cb.Callback(line)
				assert.ErrorIs(t, err, test.wantErr)

				if err != nil {
					break
				}
			}

			gotMassdns := filetest.ReadFile(t, massdnsFile.Name())
			gotDomain := filetest.ReadFile(t, domainFile.Name())

			assert.Equal(t, test.wantMassdns, gotMassdns)
			assert.Equal(t, test.wantDomain, gotDomain)
		})
	}
}

func TestDefaultWriteCallback_NoWriter(t *testing.T) {
	cb, err := NewDefaultWriteCallback("", "")
	require.Nil(t, err)

	gotErr := cb.Callback("")
	assert.Nil(t, gotErr)
}

func TestDefaultWriteCallbackClose(t *testing.T) {
	massdnsFile := filetest.CreateFile(t, "")
	domainFile := filetest.CreateFile(t, "")

	cb, err := NewDefaultWriteCallback(massdnsFile.Name(), domainFile.Name())
	assert.Nil(t, err)

	cb.Close()

	assert.Equal(t, uintptr(0xffffffffffffffff), cb.massdnsFile.Fd())
	assert.Equal(t, uintptr(0xffffffffffffffff), cb.domainFile.Fd())
}
