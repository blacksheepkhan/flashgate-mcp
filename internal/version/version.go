package version

import "fmt"

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

// Info contains build and release metadata.
type Info struct {
	Version string
	Commit  string
	Date    string
}

// Get returns the build metadata embedded in the binary.
func Get() Info {
	return Info{
		Version: version,
		Commit:  commit,
		Date:    date,
	}
}

// String returns a human-readable version string.
func (i Info) String() string {
	return fmt.Sprintf("fileserver-mcp\nversion: %s\ncommit: %s\ndate: %s", i.Version, i.Commit, i.Date)
}
