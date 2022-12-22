package resolve

import (
	"testing"

	"github.com/d3mondev/puredns/v2/internal/app/ctx"
	"github.com/d3mondev/puredns/v2/pkg/filetest"
	"github.com/stretchr/testify/assert"
)

func TestSave(t *testing.T) {
	domainFile := filetest.CreateFile(t, "domain.com")
	massdnsFile := filetest.CreateFile(t, "domain.com. A 127.0.0.1")
	wildcardFile := filetest.CreateFile(t, "*.domain.com")

	domainOutputFile := filetest.CreateFile(t, "")
	massdnsOutputFile := filetest.CreateFile(t, "")
	wildcardOutputFile := filetest.CreateFile(t, "")

	tests := []struct {
		name                  string
		haveReadDomainFile    string
		haveWriteDomainFile   string
		haveReadMassdnsFile   string
		haveWriteMassdnsFile  string
		haveReadWildcardFile  string
		haveWriteWildcardFile string
		wantDomainContent     []string
		wantMassdnsContent    []string
		wantWildcardContent   []string
		wantErr               bool
	}{
		{
			name:    "don't save",
			wantErr: false,
		},
		{
			name:                  "save all",
			haveReadDomainFile:    domainFile.Name(),
			haveWriteDomainFile:   domainOutputFile.Name(),
			haveReadMassdnsFile:   massdnsFile.Name(),
			haveWriteMassdnsFile:  massdnsOutputFile.Name(),
			haveReadWildcardFile:  wildcardFile.Name(),
			haveWriteWildcardFile: wildcardOutputFile.Name(),
			wantDomainContent:     []string{"domain.com"},
			wantMassdnsContent:    []string{"domain.com. A 127.0.0.1"},
			wantWildcardContent:   []string{"*.domain.com"},
		},
		{
			name:                "domain file error handling",
			haveReadDomainFile:  "thisfiledoesntexist.txt",
			haveWriteDomainFile: domainOutputFile.Name(),
			wantErr:             true,
		},
		{
			name:                 "massdns answers file error handling",
			haveReadMassdnsFile:  "thisfiledoesntexist.txt",
			haveWriteMassdnsFile: massdnsOutputFile.Name(),
			wantErr:              true,
		},
		{
			name:                  "wildcard roots file error handling",
			haveReadWildcardFile:  "thisfiledoesntexist.txt",
			haveWriteWildcardFile: wildcardOutputFile.Name(),
			wantErr:               true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filetest.ClearFile(t, domainOutputFile)
			filetest.ClearFile(t, massdnsOutputFile)
			filetest.ClearFile(t, wildcardOutputFile)

			opt := &ctx.ResolveOptions{}
			opt.WriteDomainsFile = test.haveWriteDomainFile
			opt.WriteMassdnsFile = test.haveWriteMassdnsFile
			opt.WriteWildcardsFile = test.haveWriteWildcardFile

			workfiles := &Workfiles{}
			workfiles.Domains = test.haveReadDomainFile
			workfiles.MassdnsPublic = test.haveReadMassdnsFile
			workfiles.WildcardRoots = test.haveReadWildcardFile

			saver := NewResultFileSaver()

			gotErr := saver.Save(workfiles, opt)

			assert.Equal(t, test.wantErr, gotErr != nil)
			assert.ElementsMatch(t, test.wantDomainContent, filetest.ReadFile(t, test.haveWriteDomainFile))
			assert.ElementsMatch(t, test.wantMassdnsContent, filetest.ReadFile(t, test.haveWriteMassdnsFile))
			assert.ElementsMatch(t, test.wantWildcardContent, filetest.ReadFile(t, test.haveWriteWildcardFile))
		})
	}
}
