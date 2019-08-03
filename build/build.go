package build

import (
	"fmt"
	"runtime"
)

var (
	// GitCommit git commit info
	GitCommit string
	// GitBranch git branch info
	GitBranch string
)

// Version version info
func Version() string {
	return fmt.Sprintf("commit: %s\n\tbranch: %s\n\tgo:%s", GitCommit, GitBranch, runtime.Version())
}
