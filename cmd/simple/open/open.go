package open

import (
	"context"

	"github.com/nuggxyz/buildrc/internal/common"
	"github.com/nuggxyz/buildrc/internal/git"
	"github.com/nuggxyz/buildrc/internal/pipeline"
)

type Handler struct {
}

func (me *Handler) Run(ctx context.Context, prov common.Provider) (err error) {

	targetSemver, err := git.CalculateNextPreReleaseTag(ctx, prov.Buildrc(), prov.Git(), prov.PR())
	if err != nil {
		return err
	}

	next, err := prov.Release().TagRelease(ctx, prov.Git(), targetSemver)
	if err != nil {
		return err
	}

	err = pipeline.AddContentToEnv(ctx, prov.Pipeline(), prov.FileSystem(), "open", map[string]string{
		"BUILDRC_RELEASE_ID":  next.ID,
		"BUILDRC_RELEASE_TAG": targetSemver.String(),
	})
	if err != nil {
		return err
	}

	return err
}
