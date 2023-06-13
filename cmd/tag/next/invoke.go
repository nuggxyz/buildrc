package next

import (
	"context"
	"fmt"

	"github.com/nuggxyz/buildrc/cmd/buildrc/load"
	"github.com/nuggxyz/buildrc/cmd/tag/list"
	"github.com/nuggxyz/buildrc/internal/docker"
	"github.com/nuggxyz/buildrc/internal/env"
	"github.com/nuggxyz/buildrc/internal/github"
	"github.com/nuggxyz/buildrc/internal/provider"
	"github.com/rs/zerolog"
)

const (
	CommandID = "tag_next"
)

type Output struct {
	Major           string `json:"major"`
	Minor           string `json:"minor"`
	Patch           string `json:"patch"`
	MajorMinor      string `json:"major_minor"`
	MajorMinorPatch string `json:"major_minor_patch"`
	Full            string `json:"full" express:"BUILDRC_TAG_NEXT_FULL"`
	BuildxTags      string `json:"buildx_tags" express:"BUILDRC_TAG_NEXT_BUILDX_TAGS"`
}

type Handler struct {
	Repo        string `flag:"repo" type:"repo:" default:""`
	File        string `flag:"file" type:"file:" default:".buildrc"`
	AccessToken string `flag:"token" type:"access_token:" default:""`
}

func NewHandler(repo string, accessToken string) *Handler {
	h := &Handler{Repo: repo, AccessToken: accessToken}
	return h
}

func (me *Handler) Run(ctx context.Context, cp provider.ContentProvider) (err error) {
	_, err = me.Next(ctx, cp)
	return err
}

func (me *Handler) Next(ctx context.Context, cp provider.ContentProvider) (out *Output, err error) {
	return provider.Wrap(CommandID, me.next)(ctx, cp)
}

func (me *Handler) next(ctx context.Context, prv provider.ContentProvider) (out *Output, err error) {

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
			return nil, err
		}

		zerolog.Ctx(ctx).Debug().Msgf("✅ Repo found in env: %s", curr)

		me.Repo = curr
	}

	brc, err := load.NewHandler(me.File).Load(ctx, prv)
	if err != nil {
		return nil, err
	}

	prov, err := list.NewHandler(me.Repo, me.AccessToken).Invoke(ctx, prv)
	if err != nil {
		return nil, err
	}

	// Increment patch version
	nextVersion := prov.HighestVersion.IncPatch()

	// If the buildrc version is higher than the highest version, use that instead
	if brc.Version.GreaterThan(&nextVersion) {
		nextVersion = *brc.Version
	}

	str, err := docker.BuildXTagString(ctx, me.Repo, nextVersion.String())
	if err != nil {
		return nil, err
	}

	return &Output{
		Major:           fmt.Sprintf("%d", nextVersion.Major()),
		Minor:           fmt.Sprintf("%d", nextVersion.Minor()),
		Patch:           fmt.Sprintf("%d", nextVersion.Patch()),
		MajorMinor:      fmt.Sprintf("%d.%d", nextVersion.Major(), nextVersion.Minor()),
		MajorMinorPatch: fmt.Sprintf("%d.%d.%d", nextVersion.Major(), nextVersion.Minor(), nextVersion.Patch()),
		Full:            nextVersion.String(),
		BuildxTags:      str,
	}, nil
}
