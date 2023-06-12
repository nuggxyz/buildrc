package next

import (
	"context"
	"fmt"

	buildrc "github.com/nuggxyz/buildrc/cmd/buildrc/load"
	"github.com/nuggxyz/buildrc/cmd/tag/list"
	"github.com/nuggxyz/buildrc/internal/env"
	"github.com/nuggxyz/buildrc/internal/github"
	"github.com/nuggxyz/buildrc/internal/provider"
	"github.com/rs/zerolog"
)

type Handler struct {
	Repo        string `flag:"repo" type:"repo:" default:""`
	BuildrcFile string `flag:"file" type:"file:" default:".buildrc"`
	AccessToken string `flag:"token" type:"access_token:" default:""`

	gettagsHandler *list.Handler
	buildrcHandler *buildrc.Handler
}

func (me *Handler) Init(ctx context.Context) (err error) {

	if me.AccessToken == "" {
		zerolog.Ctx(ctx).Debug().Msg("No access token provided, trying to get from env")
		// TODO: this should be a helper function, could grab from somewhere else
		me.AccessToken = env.GetOrEmpty("GITHUB_TOKEN")
		if me.AccessToken == "" {
			zerolog.Ctx(ctx).Debug().Msg("❌ No access token found in env")
		} else {
			zerolog.Ctx(ctx).Debug().Msg("✅ Access token found in env")
		}
	}

	if me.Repo == "" {

		zerolog.Ctx(ctx).Debug().Msg("No repo provided, trying to get from env")

		curr, err := github.GetCurrentRepo()
		if err != nil {
			return err
		}

		zerolog.Ctx(ctx).Debug().Msgf("✅ Repo found in env: %s", curr)

		me.Repo = curr
	}

	me.gettagsHandler, err = list.NewHandler(ctx, me.Repo, me.AccessToken)
	if err != nil {
		return err
	}

	me.buildrcHandler, err = buildrc.NewHandler(ctx, me.BuildrcFile)
	if err != nil {
		return err
	}
	return
}

func (me *Handler) Invoke(ctx context.Context, prv provider.ContentProvider) (out *TagNextOutput, err error) {

	prov, err := me.gettagsHandler.Helper().Run(ctx, prv)
	if err != nil {
		return nil, err
	}

	brc, err := me.buildrcHandler.Helper().Run(ctx, prv)
	if err != nil {
		return nil, err
	}

	// Increment patch version
	nextVersion := prov.HighestVersion.IncPatch()

	// If the buildrc version is higher than the highest version, use that instead
	if brc.Version.GreaterThan(&nextVersion) {
		nextVersion = *brc.Version
	}

	return &TagNextOutput{
		Major:           fmt.Sprintf("%d", nextVersion.Major()),
		Minor:           fmt.Sprintf("%d", nextVersion.Minor()),
		Patch:           fmt.Sprintf("%d", nextVersion.Patch()),
		MajorMinor:      fmt.Sprintf("%d.%d", nextVersion.Major(), nextVersion.Minor()),
		MajorMinorPatch: fmt.Sprintf("%d.%d.%d", nextVersion.Major(), nextVersion.Minor(), nextVersion.Patch()),
		Full:            nextVersion.String(),
	}, nil
}
