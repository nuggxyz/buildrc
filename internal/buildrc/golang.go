package buildrc

import (
	"fmt"

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

	return fmt.Sprintf("-X main.Version=%s -X main.RawVersion=%s -X main.Revision=%s", ver, raw, commit), nil
}
