package prebuild

import (
	"context"

	"github.com/nuggxyz/buildrc/internal/provider"
)

type output struct {
}

var _ provider.CommandRunner = (*Handler)(nil)
var _ provider.Command[output] = (*Handler)(nil)

func (me *Handler) ID() string {
	return "prebuild-hook"
}

func NewHandler(ctx context.Context, file string) (*Handler, error) {
	h := &Handler{
		BuildrcFile: file,
	}

	err := h.Init(ctx)

	return h, err

}

func (me *Handler) Helper() provider.CommandHelper[output] {
	return provider.NewHelper[output](me)
}

func (me *Handler) AnyHelper() provider.AnyHelper {
	return provider.NewHelper[output](me)
}
