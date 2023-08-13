package buildrc

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
)

func (me *Package) UsesMap() map[string]string {
	m := make(map[string]string)
	for _, use := range me.Uses {
		m[use] = "1"
	}
	return m
}

func StringsToCSV[I ~string](ss []I) string {
	strs := make([]string, len(ss))
	for i, s := range ss {
		strs[i] = string(s)
	}
	return strings.Join(strs, ",")
}

func (me *Package) ArtifactFileNames() ([]string, error) {
	names := make([]string, 0)
	for _, s := range me.Platforms {
		tmp, err := s.OutputFile(me)
		if err != nil {
			return nil, err
		}
		names = append(names, tmp+".tar.gz", tmp+".sha256")
	}
	return names, nil
}

func (me *Package) ToArtifactCSV(ss []Platform) (string, error) {
	names, err := me.ArtifactFileNames()
	if err != nil {
		return "", err
	}

	return strings.Join(names, ","), nil
}

func (me *Package) CustomJSON() (string, error) {
	if me.Custom == nil {
		return "{}", nil
	}
	cust, err := json.Marshal(me.Custom)
	if err != nil {
		return "", err
	}

	return string(cust), nil
}

func (me *Package) TestArchiveFileName() string {
	return fmt.Sprintf("%s-test-output.tar.gz", me.Name)
}

func (me *Package) VerifyArchiveFileName() string {
	return fmt.Sprintf("%s-test-output.tar.gz", me.Name)
}

func (me DockerBuildArgs) CSV() (string, error) {

	return strings.Join(me.Array(), ","), nil
}

func (me DockerBuildArgs) JSONString() (string, error) {
	args := me.Array()

	joiner := strings.Join(args, "\n")

	res, err := json.Marshal(joiner)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func (me *Package) Dockerfile() string {
	return filepath.Join(me.Dir, "Dockerfile")
}

func (me *Package) DockerPlatformsCSV() string {
	return StringsToCSV(me.DockerPlatforms)
}
