package runner

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/nuggxyz/buildrc/internal/cli"
	"github.com/nuggxyz/buildrc/internal/env"
	"github.com/nuggxyz/buildrc/internal/file"
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

var _ cli.ContentProvider = (*GHActionContentProvider)(nil)

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

func (me *GHActionContentProvider) tmpFileName(cmd cli.Identifiable) string {
	return me.RUNNER_TEMP + "/" + cmd.ID() + ".json"
}

func (me *GHActionContentProvider) Load(ctx context.Context, cmd cli.Identifiable) ([]byte, error) {
	// try to load from tmp folder
	f, err := me.fs.Get(ctx, me.tmpFileName(cmd))
	if err != nil {

		// if not found do nothing
		if errors.Is(err, os.ErrNotExist) {
			return []byte{}, nil
		}

		return nil, err
	}

	zerolog.Ctx(ctx).Debug().Str("data", string(f)).Msgf("loaded result from %s", me.tmpFileName(cmd))

	return f, nil
}

func (me *GHActionContentProvider) Save(ctx context.Context, cmd cli.Identifiable, result []byte) error {

	err := me.fs.AppendString(ctx, me.GITHUB_OUTPUT, fmt.Sprintf("result=%s", string(result)))
	if err != nil {
		return err
	}

	// save to tmp folder
	err = me.fs.Put(ctx, me.tmpFileName(cmd), result)
	if err != nil {
		return err
	}

	return nil
}
