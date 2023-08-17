package pipeline

import (
	"encoding/json"

	"github.com/walteh/buildrc/internal/buildrc"
)

type PipelineRunsOn string

const (
	MacOS   PipelineRunsOn = "macos"
	Linux   PipelineRunsOn = "linux"
	Windows PipelineRunsOn = "windows"
	Custom  PipelineRunsOn = "custom"
)

func (p PipelineRunsOn) String() string {
	return string(p)
}

func ResolveRunsOnMap(brc *buildrc.Buildrc, pip Pipeline) (map[string]string, error) {
	mapper := map[string]string{}

	for _, os := range brc.Packages {
		a, err := pip.RunsOnResolution(PipelineRunsOn(os.On))
		if err != nil {
			panic(err)
		}
		mapper[os.Name] = a
	}

	return mapper, nil
}

func ResolveRunsOnMapJSON(brc *buildrc.Buildrc, pip Pipeline) (string, error) {
	m, err := ResolveRunsOnMap(brc, pip)
	if err != nil {
		return "", err
	}

	by, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	return string(by), nil
}
