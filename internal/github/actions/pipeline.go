package actions

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/nuggxyz/buildrc/internal/pipeline"
	"github.com/rs/zerolog"
	"github.com/spf13/afero"
)

type GithubActionPipeline struct {
	// RUNNER_TEMP    string
	// GITHUB_ENV     string
	// GITHUB_OUTPUT  string
	// GITHUB_ACTIONS string
	// CI             string
	// fs             file.FileAPI
}

var _ pipeline.Pipeline = (*GithubActionPipeline)(nil)

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
