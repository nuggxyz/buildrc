package finalize

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/nuggxyz/buildrc/cmd/buildrc/load"
	"github.com/nuggxyz/buildrc/internal/docker"
	"github.com/nuggxyz/buildrc/internal/github"
	"github.com/nuggxyz/buildrc/internal/provider"
	"github.com/rs/zerolog"
)

const (
	CommandID = "finalize"
)

type Output struct {
	Major           string `json:"major"`
	Minor           string `json:"minor"`
	Patch           string `json:"patch"`
	MajorMinor      string `json:"major_minor"`
	MajorMinorPatch string `json:"major_minor_patch"`
	Full            string `json:"full" express:"BUILDRC_RELEASE_FINALIZE_FULL"`
	BuildxTags      string `json:"buildx_tags" express:"BUILDRC_RELEASE_FINALIZE_BUILDX_TAGS"`
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

	brc, err := load.NewHandler(me.File).Load(ctx, prv)
	if err != nil {
		return nil, err
	}

	ghc, err := github.NewGithubClient(ctx, me.AccessToken, me.Repo)
	if err != nil {
		return nil, err
	}

	vers, err := ghc.Finalize(ctx)
	if err != nil {
		return nil, err
	}

	zerolog.Ctx(ctx).Debug().Str("next-version", vers.String()).Str("buildrc-version", brc.Version.String()).Msg("Calculated next version")

	if brc.Version.GreaterThan(vers) {
		vers = semver.New(brc.Version.Major(), 0, 0, vers.Prerelease(), vers.Metadata())
	}

	str, err := docker.BuildXTagString(ctx, me.Repo, vers.String())
	if err != nil {
		return nil, err
	}

	return &Output{
		Major:           fmt.Sprintf("%d", vers.Major()),
		Minor:           fmt.Sprintf("%d", vers.Minor()),
		Patch:           fmt.Sprintf("%d", vers.Patch()),
		MajorMinor:      fmt.Sprintf("%d.%d", vers.Major(), vers.Minor()),
		MajorMinorPatch: fmt.Sprintf("%d.%d.%d", vers.Major(), vers.Minor(), vers.Patch()),
		Full:            vers.String(),
		BuildxTags:      str,
	}, nil
}
