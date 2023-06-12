package load

import (
	"context"

	"github.com/nuggxyz/buildrc/internal/buildrc"
	"github.com/nuggxyz/buildrc/internal/provider"
)

func (me *Handler) Invoke(ctx context.Context, r provider.ContentProvider) (out *buildrc.BuildRC, err error) {

	out, err = buildrc.Parse(ctx, me.File)
	if err != nil {
		return nil, err
	}

	return
}
