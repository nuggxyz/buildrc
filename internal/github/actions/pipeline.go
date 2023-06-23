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

	// var err error

	// check if we are in a github action

	// ci1 := os.Getenv("CI")

	// zerolog.Ctx(ctx).Debug().Str("CI", ci1).Msg("CI")

	// if obj.CI, err = env.Get("CI"); err != nil {
	// 	return nil, err
	// }

	// if obj.GITHUB_ACTIONS, err = env.Get("GITHUB_ACTIONS"); err != nil {
	// 	return nil, err
	// }

	// if obj.CI != "true" || obj.GITHUB_ACTIONS != "true" {
	// 	return nil, errors.New("not in a github action")
	// }

	// if obj.GITHUB_ENV, err = env.Get("GITHUB_ENV"); err != nil {
	// 	return nil, err
	// }

	// if obj.GITHUB_OUTPUT, err = env.Get("GITHUB_OUTPUT"); err != nil {
	// 	return nil, err
	// }

	// if obj.RUNNER_TEMP, err = env.Get("RUNNER_TEMP"); err != nil {
	// 	return nil, err
	// }

	// zerolog.Ctx(ctx).Debug().Any("obj", obj).Msg("new ghaction content provider")

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
