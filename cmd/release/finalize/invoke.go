package finalize

import (
	"context"
	"fmt"

	"github.com/nuggxyz/buildrc/internal/common"
	"github.com/nuggxyz/buildrc/internal/git"
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
	return provider.Wrap(ctx, CommandID, prov, me.invoke)
}

func (me *Handler) invoke(ctx context.Context, prov common.Provider) (out *Output, err error) {

	curr, err := git.GetCurrentRelease(ctx, prov.Release(), prov.Git())
	if err != nil {
		return nil, err
	}

	vers, err := prov.Release().MakeReleaseLive(ctx, curr)
	if err != nil {
		return nil, err
	}

	zerolog.Ctx(ctx).Debug().Str("next-version", vers.String()).Int("buildrc-major-version", prov.Buildrc().Version).Msg("Calculated next version")

	return &Output{
		Major:           fmt.Sprintf("%d", vers.Major()),
		Minor:           fmt.Sprintf("%d", vers.Minor()),
		Patch:           fmt.Sprintf("%d", vers.Patch()),
		MajorMinor:      fmt.Sprintf("%d.%d", vers.Major(), vers.Minor()),
		MajorMinorPatch: fmt.Sprintf("%d.%d.%d", vers.Major(), vers.Minor(), vers.Patch()),
		Full:            vers.String(),
	}, nil
}
