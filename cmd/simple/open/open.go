package open

import (
	"context"

	"github.com/Masterminds/semver/v3"
	"github.com/nuggxyz/buildrc/internal/common"
	"github.com/nuggxyz/buildrc/internal/pipeline"
)

type Handler struct {
}

func (me *Handler) Run(ctx context.Context, prov common.Provider) (err error) {

	rid, err := prov.Pipeline().GetFromEnv(ctx, "BUILDRC_RELEASE_TAG", prov.FileSystem())
	if err != nil {
		return err
	}

	smvr, err := semver.NewVersion(rid)
	if err != nil {
		return err
	}

	next, err := prov.Release().TagRelease(ctx, prov.Git(), smvr)
	if err != nil {
		return err
	}

	err = pipeline.AddContentToEnv(ctx, prov.Pipeline(), prov.FileSystem(), "open", map[string]string{
		"BUILDRC_RELEASE_ID":  next.ID,
		"BUILDRC_RELEASE_TAG": smvr.String(),
	})
	if err != nil {
		return err
	}

	return err
}
