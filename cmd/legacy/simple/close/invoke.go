package close

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/walteh/buildrc/internal/common"
)

const ()

type Handler struct {
}

func (me *Handler) Run(ctx context.Context, prov common.Provider) (err error) {

	rid, err := prov.Pipeline().GetFromEnv(ctx, "BUILDRC_RELEASE_ID", prov.FileSystem())
	if err != nil {
		return err
	}

	rel, err := prov.Release().GetReleaseByID(ctx, rid)
	if err != nil {
		return err
	}

	err = prov.Release().TakeReleaseOutOfDraft(ctx, rel)
	if err != nil {
		return err
	}

	zerolog.Ctx(ctx).Debug().Any("next", rel.ID).Msg("tagged release")

	return nil
}
