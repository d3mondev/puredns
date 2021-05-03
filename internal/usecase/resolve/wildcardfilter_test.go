package resolve

import (
	"testing"

	"github.com/d3mondev/puredns/v2/pkg/filetest"
	"github.com/stretchr/testify/assert"
)

func TestWildcardFilter(t *testing.T) {
	cacheFile := filetest.CreateFile(t, "")
	domainFile := filetest.CreateFile(t, "")
	rootFile := filetest.CreateFile(t, "")

	wc := NewDefaultWildcardFilter()

	opt := WildcardFilterOptions{
		CacheFilename:        cacheFile.Name(),
		DomainOutputFilename: domainFile.Name(),
		RootOutputFilename:   rootFile.Name(),
		Resolvers:            []string{},
		QueriesPerSecond:     10,
		ThreadCount:          1,
	}
	_, _, err := wc.Filter(opt, 1)

	assert.Nil(t, err)
}

func TestWildcardFilter_Files(t *testing.T) {
	badFile := filetest.CreateDir(t)
	cacheFile := filetest.CreateFile(t, "example.com. A 127.0.0.1")
	domainFile := filetest.CreateFile(t, "")
	rootFile := filetest.CreateFile(t, "")

	tests := []struct {
		name                 string
		haveCacheFile        string
		haveOutputDomainFile string
		haveOutputRootFile   string
		wantErr              bool
	}{
		{
			name:                 "valid files",
			haveCacheFile:        cacheFile.Name(),
			haveOutputDomainFile: domainFile.Name(),
			haveOutputRootFile:   rootFile.Name(),
			wantErr:              false,
		},
		{
			name:                 "cache file error handling",
			haveCacheFile:        "",
			haveOutputDomainFile: domainFile.Name(),
			haveOutputRootFile:   rootFile.Name(),
			wantErr:              true,
		},
		{
			name:                 "output file error handling",
			haveCacheFile:        cacheFile.Name(),
			haveOutputDomainFile: badFile,
			haveOutputRootFile:   rootFile.Name(),
			wantErr:              true,
		},
		{
			name:                 "root file error handling",
			haveCacheFile:        cacheFile.Name(),
			haveOutputDomainFile: domainFile.Name(),
			haveOutputRootFile:   badFile,
			wantErr:              true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filter := NewDefaultWildcardFilter()

			opt := WildcardFilterOptions{
				CacheFilename:        test.haveCacheFile,
				DomainOutputFilename: test.haveOutputDomainFile,
				RootOutputFilename:   test.haveOutputRootFile,
				Resolvers:            []string{},
				QueriesPerSecond:     10,
				ThreadCount:          1,
			}

			_, _, err := filter.Filter(opt, 0)

			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestQPSPerResolver(t *testing.T) {
	tests := []struct {
		name              string
		haveResolverCount int
		haveGlobalQPS     int
		want              int
	}{
		{name: "no resolvers", haveResolverCount: 0, haveGlobalQPS: 10, want: 0},
		{name: "two resolvers", haveResolverCount: 2, haveGlobalQPS: 10, want: 5},
		{name: "many resolvers", haveResolverCount: 10, haveGlobalQPS: 1, want: 1},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := qpsPerResolver(test.haveResolverCount, test.haveGlobalQPS)

			assert.Equal(t, test.want, got)
		})
	}
}
