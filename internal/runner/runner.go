package runner

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/nuggxyz/buildrc/internal/env"
	"github.com/nuggxyz/buildrc/internal/file"
	"github.com/nuggxyz/buildrc/internal/provider"
	"github.com/rs/zerolog"
)

type GHActionContentProvider struct {
	RUNNER_TEMP    string
	GITHUB_ENV     string
	GITHUB_OUTPUT  string
	GITHUB_ACTIONS string
	CI             string
	fs             file.FileAPI
}

var _ provider.ContentProvider = (*GHActionContentProvider)(nil)

func IAmInAGithubAction(ctx context.Context) bool {
	ci1 := os.Getenv("CI")

	if ci1 != "true" {
		return false
	}

	ga1 := os.Getenv("GITHUB_ACTIONS")

	return ga1 == "true"
}

func NewGHActionContentProvider(ctx context.Context, api file.FileAPI) (*GHActionContentProvider, error) {

	obj := &GHActionContentProvider{
		fs: api,
	}
	var err error

	// check if we are in a github action

	ci1 := os.Getenv("CI")

	zerolog.Ctx(ctx).Debug().Str("CI", ci1).Msg("CI")

	if obj.CI, err = env.Get("CI"); err != nil {
		return nil, err
	}

	if obj.GITHUB_ACTIONS, err = env.Get("GITHUB_ACTIONS"); err != nil {
		return nil, err
	}

	if obj.CI != "true" || obj.GITHUB_ACTIONS != "true" {
		return nil, errors.New("not in a github action")
	}

	if obj.GITHUB_ENV, err = env.Get("GITHUB_ENV"); err != nil {
		return nil, err
	}

	if obj.GITHUB_OUTPUT, err = env.Get("GITHUB_OUTPUT"); err != nil {
		return nil, err
	}

	if obj.RUNNER_TEMP, err = env.Get("RUNNER_TEMP"); err != nil {
		return nil, err
	}

	zerolog.Ctx(ctx).Debug().Any("obj", obj).Msg("new ghaction content provider")

	return obj, nil
}

func (me *GHActionContentProvider) FileSystem() file.FileAPI {
	return me.fs
}

func (me *GHActionContentProvider) AddToEnv(ctx context.Context, id string, val string) error {
	return me.fs.AppendString(ctx, me.GITHUB_ENV, fmt.Sprintf("%s=%s\n", id, val))
}
