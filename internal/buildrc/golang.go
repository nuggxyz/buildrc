package buildrc

import (
	"fmt"
	"time"

	"github.com/Masterminds/semver/v3"
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

	return fmt.Sprintf("-X main.Version=%s -X main.RawVersion=%s -X main.Revision=%s -X main.Time=%d", ver, raw, commit, time.Now().Unix()), nil
}
