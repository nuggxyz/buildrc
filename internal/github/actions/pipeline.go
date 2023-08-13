package actions

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/nuggxyz/buildrc/internal/pipeline"
	"github.com/rs/zerolog"
	"github.com/spf13/afero"
)

type GithubActionPipeline struct {
}

var _ pipeline.Pipeline = (*GithubActionPipeline)(nil)

func (me *GithubActionPipeline) SupportsDocker() bool {
	return runtime.GOOS == "linux"
}

func IAmInAGithubAction(ctx context.Context) bool {
	return EnvVarCI.Load() != "" && EnvVarGithubActions.Load() != ""
}

func NewGithubActionPipeline(ctx context.Context) (*GithubActionPipeline, error) {

	obj := &GithubActionPipeline{}

	if IAmInAGithubAction(ctx) {
		zerolog.Ctx(ctx).Debug().Msg("github action detected")
		return obj, nil
	} else {
		return nil, errors.New("not in a github action")
	}

}

func (me *GithubActionPipeline) AddToEnv(ctx context.Context, id string, val string, fs afero.Fs) error {

	zerolog.Ctx(ctx).Debug().Str("id", id).Str("val", val).Msg("AddToEnv")

	envfile := EnvVarGithubEnv.Load()

	if envfile == "" {
		return fmt.Errorf("env var %s not set", EnvVarGithubEnv)
	}

	fle, err := fs.OpenFile(envfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer fle.Close()

	_, err = fle.WriteString(fmt.Sprintf("%s=%s\n", id, val))
	if err != nil {
		return err
	}

	return nil
}

func (me *GithubActionPipeline) GetFromEnv(ctx context.Context, id string, fs afero.Fs) (string, error) {
	res := os.Getenv(id)
	if res == "" {
		return "", fmt.Errorf("env var %s not set", id)
	}

	return res, nil
}

func (me *GithubActionPipeline) RunId(ctx context.Context) (int64, error) {

	res := EnvVarGithubRunID.Load()
	if res == "" {
		return 0, fmt.Errorf("env var %s not set", EnvVarGithubRunID)
	}

	resp, err := strconv.ParseInt(res, 10, 64)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Any("res", res).Msg("failed to parse run id")
		return 0, err
	}

	return resp, nil

}

func (me *GithubActionPipeline) RootDir(ctx context.Context) (string, error) {
	res := EnvVarGithubWorkspace.Load()
	if res == "" {
		return "", fmt.Errorf("env var %s not set", EnvVarGithubWorkspace)
	}

	return res, nil
}

func (me *GithubActionPipeline) RunsOnResolution(osp pipeline.PipelineRunsOn) (string, error) {
	switch osp {
	case pipeline.Linux:
		return "ubuntu-latest", nil
	case pipeline.Windows:
		return "windows-latest", nil
	case pipeline.MacOS:
		return "macos-latest", nil
	case pipeline.Custom:
		return "self-hosted", nil
	}

	return "", fmt.Errorf("unknown os %v", osp)
}
