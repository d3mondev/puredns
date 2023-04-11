package app

const (
	// AppName is the name of the application.
	AppName string = "puredns"

	// AppDesc is a short description of the application.
	AppDesc string = "Very accurate massdns resolving and bruteforcing."

	// AppVersion is the program version.
	AppVersion string = "v2.1.1"

	// AppSponsorsURL is the text file containing sponsors information.
	AppSponsorsURL string = "https://gist.githubusercontent.com/d3mondev/0bfff529a4dad627bdb684ad1ef2506d/raw/sponsors.txt"
)

// GitBranch is the current git branch.
var GitBranch string

// GitRevision is the current git commit.
var GitRevision string
