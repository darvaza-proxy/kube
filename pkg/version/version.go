// Package version provides information about the build
package version

import (
	"fmt"
	"runtime"
	"strconv"
	"time"
)

const (
	undefined = "undefined"
)

// Values are set by compile time -ldflags
var (
	Version   = undefined
	Branch    = undefined
	Commit    = undefined
	BuildDate = undefined
)

// Info contains build information supplied during compile time
type Info struct {
	Version   string
	Branch    string
	Commit    string
	BuildDate time.Time
	GoVersion string
	Platform  string
}

// Get returns the built-in version and platform information
func Get() Info {
	var ts time.Time

	if BuildDate != undefined {
		n, err := strconv.ParseInt(BuildDate, 10, 64)
		if err != nil {
			panic(fmt.Errorf("invalid build time: %v", err))
		}
		ts = time.Unix(n, 0).UTC()
	}

	v := Info{
		Version:   Version,
		Branch:    Branch,
		Commit:    Commit,
		BuildDate: ts,
		GoVersion: runtime.Version(),
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}

	return v
}
