package build

import (
	"fmt"
	"runtime"
)

var (
	// GitCommit git commit info
	Commit string
	// GitBranch git branch info
	Branch string
	// Date
	Date string
)

// Version version info
func Version() string {
	return fmt.Sprintf("build: %s\ncommit: %s\nbranch: %s\ngo: %s", Date, Commit, Branch, runtime.Version())
}
