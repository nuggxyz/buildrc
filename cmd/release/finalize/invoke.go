package finalize

import (
	"context"

	"github.com/nuggxyz/buildrc/internal/common"
	"github.com/nuggxyz/buildrc/internal/git"
	"github.com/nuggxyz/buildrc/internal/pipeline"
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
	File string `flag:"file" type:"file:" default:".buildrc"`
}

func NewHandler() *Handler {
	h := &Handler{}
	return h
}

func (me *Handler) Run(ctx context.Context, prov common.Provider) (err error) {
	_, err = me.Invoke(ctx, prov)
	return err
}

func (me *Handler) Invoke(ctx context.Context, prov common.Provider) (out *Output, err error) {
	return pipeline.Cache(ctx, CommandID, prov, me.invoke)
}

func (me *Handler) invoke(ctx context.Context, prov common.Provider) (out *Output, err error) {

	curr, err := git.GetCurrentRelease(ctx, prov.Release(), prov.Git())
	if err != nil {
		return nil, err
	}

	vers, err := git.CalculateNextPreReleaseTag(ctx, prov.Buildrc(), prov.Git(), prov.PR())
	if err != nil {
		return nil, err
	}

	commit, err := prov.Git().GetCurrentCommitHash(ctx)
	if err != nil {
		return nil, err
	}

	next, err := prov.Release().TagRelease(ctx, curr, vers, commit)
	if err != nil {
		return nil, err
	}

	zerolog.Ctx(ctx).Debug().Any("next", next).Any("prev", curr).Msg("tagged release")

	return &Output{
		// Major:           fmt.Sprintf("%d", vers.Major()),
		// Minor:           fmt.Sprintf("%d", vers.Minor()),
		// Patch:           fmt.Sprintf("%d", vers.Patch()),
		// MajorMinor:      fmt.Sprintf("%d.%d", vers.Major(), vers.Minor()),
		// MajorMinorPatch: fmt.Sprintf("%d.%d.%d", vers.Major(), vers.Minor(), vers.Patch()),
		Full: vers.String(),
	}, nil
}
