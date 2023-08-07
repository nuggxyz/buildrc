package version

import (
	"fmt"
	"time"

	"github.com/Masterminds/semver/v3"
)

var (
	// Package is filled at linking time
	Package = "github.com/walteh/buildrc"

	// Version holds the complete version number. Filled in at linking time.
	Version = ""

	// Revision is filled with the VCS (e.g. git) revision being used to build
	// the program at linking time.
	Revision = ""

	Time = 0

	RawVersion = ""

	Temporary = false
)

func GenerateGoLdflags(pkg string, version string, commit string) (string, error) {
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

	return fmt.Sprintf("-X %s/version.Version=%s -X %s/version.RawVersion=%s -X %s/version.Revision=%s -X %s/version.Time=%d", pkg, ver, pkg, raw, pkg, commit, pkg, time.Now().Unix()), nil
}
