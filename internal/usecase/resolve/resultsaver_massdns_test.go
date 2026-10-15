package resolve

import (
	"io/ioutil"
	"testing"

	"github.com/d3mondev/puredns/v2/internal/app/ctx"
	"github.com/d3mondev/puredns/v2/pkg/filetest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveMassdnsFinalResults(t *testing.T) {
	const public = "junk.wild.synthetic.test. A 192.0.2.20\n\n" +
		"alias.synthetic.test. CNAME www.synthetic.test.\nwww.synthetic.test. A 192.0.2.10\n\n" +
		"www.synthetic.test. A 192.0.2.10\n\n" +
		"unvalidated.synthetic.test. A 192.0.2.30\n\n"
	const trusted = "www.synthetic.test. A 192.0.2.11\n\n" +
		"alias.synthetic.test. CNAME www.synthetic.test.\nwww.synthetic.test. A 192.0.2.11\n\n"
	const filteredPublic = "alias.synthetic.test. CNAME www.synthetic.test.\nwww.synthetic.test. A 192.0.2.10\n\n" +
		"www.synthetic.test. A 192.0.2.10\n\n"

	for _, test := range []struct {
		name           string
		skipValidation bool
		trustedOnly    bool
		want           string
	}{
		{name: "validated answers", want: trusted},
		{name: "skip validation", skipValidation: true, want: filteredPublic},
		{name: "trusted only", trustedOnly: true, want: filteredPublic},
	} {
		t.Run(test.name, func(t *testing.T) {
			output := filetest.CreateFile(t, "")
			files := &Workfiles{
				Domains:        filetest.CreateFile(t, "www.synthetic.test\nalias.synthetic.test\n").Name(),
				MassdnsPublic:  filetest.CreateFile(t, public).Name(),
				MassdnsTrusted: filetest.CreateFile(t, trusted).Name(),
			}
			opt := &ctx.ResolveOptions{
				WriteMassdnsFile: output.Name(),
				SkipValidation:   test.skipValidation,
				TrustedOnly:      test.trustedOnly,
			}
			require.NoError(t, NewResultFileSaver().Save(files, opt))
			got, err := ioutil.ReadFile(output.Name())
			require.NoError(t, err)
			assert.Equal(t, test.want, string(got))
		})
	}
}

func TestSaveMassdnsCNAMEAnswerSection(t *testing.T) {
	const answer = "alias.synthetic.test. CNAME target.synthetic.test.\ntarget.synthetic.test. A 192.0.2.10\n\n"
	output := filetest.CreateFile(t, "")
	files := &Workfiles{
		Domains:       filetest.CreateFile(t, "alias.synthetic.test\n").Name(),
		MassdnsPublic: filetest.CreateFile(t, answer).Name(),
	}
	opt := &ctx.ResolveOptions{WriteMassdnsFile: output.Name(), SkipValidation: true}

	require.NoError(t, NewResultFileSaver().Save(files, opt))
	got, err := ioutil.ReadFile(output.Name())
	require.NoError(t, err)
	assert.Equal(t, answer, string(got))
}

func TestSaveMassdnsNoValidDomains(t *testing.T) {
	output := filetest.CreateFile(t, "previous results\n")
	files := &Workfiles{
		Domains:       filetest.CreateFile(t, "").Name(),
		MassdnsPublic: filetest.CreateFile(t, "junk.wild.synthetic.test. A 192.0.2.20\n\n").Name(),
	}
	opt := &ctx.ResolveOptions{WriteMassdnsFile: output.Name(), SkipValidation: true}
	require.NoError(t, NewResultFileSaver().Save(files, opt))
	got, err := ioutil.ReadFile(output.Name())
	require.NoError(t, err)
	assert.Empty(t, got)
}
