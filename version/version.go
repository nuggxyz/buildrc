package version

import (
	"fmt"
	"time"

	"github.com/Masterminds/semver/v3"
)

var (
	// Package is filled at linking time
	Package = "github.com/nuggxyz/buildrc"

	// Version holds the complete version number. Filled in at linking time.
	Version = ""

	// Revision is filled with the VCS (e.g. git) revision being used to build
	// the program at linking time.
	Revision = ""

	Time = 0

	RawVersion = ""

	Temporary = false
)

func GenerateGoLdflags(version string, commit string) (string, error) {
	vers, err := semver.NewVersion(version)
	if err != nil {
		return "", err
	}

	raw := vers.String()

	var ver string

	if string(vers.Prerelease()) == "" && string(vers.Metadata()) == "" {
		ver = vers.String()
	} else {
		ver = vers.IncPatch().String()
	}

	return fmt.Sprintf("-X version.Version=%s -X version.RawVersion=%s -X version.Revision=%s -X version.Time=%d", ver, raw, commit, time.Now().Unix()), nil
}
