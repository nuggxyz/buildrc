package finalize

import (
	"context"

	"github.com/nuggxyz/buildrc/cmd/release/setup"
	"github.com/nuggxyz/buildrc/internal/common"
	"github.com/nuggxyz/buildrc/internal/pipeline"
	"github.com/rs/zerolog"
)

const (
	CommandID = "finalize"
)

type Output struct {
	Full      string `json:"full" express:"BUILDRC_RELEASE_FINALIZE_FULL"`
	ReleaseID string `json:"release_id" express:"BUILDRC_RELEASE_FINALIZE_RELEASE_ID"`
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

	su, err := setup.NewHandler("", "").Invoke(ctx, prov)
	if err != nil {
		return nil, err
	}

	next, err := prov.Release().TagRelease(ctx, prov.Git(), su.TagSemver)
	if err != nil {
		return nil, err
	}

	zerolog.Ctx(ctx).Debug().Any("next", next).Msg("tagged release")

	return &Output{
		// Major:           fmt.Sprintf("%d", vers.Major()),
		// Minor:           fmt.Sprintf("%d", vers.Minor()),
		// Patch:           fmt.Sprintf("%d", vers.Patch()),
		// MajorMinor:      fmt.Sprintf("%d.%d", vers.Major(), vers.Minor()),
		// MajorMinorPatch: fmt.Sprintf("%d.%d.%d", vers.Major(), vers.Minor(), vers.Patch()),
		Full:      su.TagSemver.String(),
		ReleaseID: next.ID,
	}, nil
}
