package zkevm

import (
	"fmt"
	"io"
	"runtime"
)

// Populated during build, don't touch!
var (
	Version   = "v0.1.0"
	GitRev    = "undefined"
	GitBranch = "undefined"
	BuildDate = "Fri, 17 Jun 1988 01:58:00 +0200"
)

// PrintVersion prints version info into the provided io.Writer.
func PrintVersion(w io.Writer) {
	fmt.Fprintf(w, "Version:      %s\n", Version)
	fmt.Fprintf(w, "Git revision: %s\n", GitRev)
	fmt.Fprintf(w, "Git branch:   %s\n", GitBranch)
	fmt.Fprintf(w, "Go version:   %s\n", runtime.Version())
	fmt.Fprintf(w, "Built:        %s\n", BuildDate)
	fmt.Fprintf(w, "OS/Arch:      %s/%s\n", runtime.GOOS, runtime.GOARCH)
}
