package prebuild

import (
	"context"

	"github.com/nuggxyz/buildrc/internal/cli"
)

type output struct {
}

var _ cli.CommandRunner = (*Handler)(nil)
var _ cli.Command[output] = (*Handler)(nil)

func (me *Handler) ID() string {
	return "prebuild-hook"
}

func NewHandler(ctx context.Context, file string, pkg string) (*Handler, error) {
	h := &Handler{
		BuildrcFile: file,
		PackageName: pkg,
	}

	err := h.Init(ctx)

	return h, err

}

func (me *Handler) Helper() cli.CommandHelper[output] {
	return cli.NewHelper[output](me)
}

func (me *Handler) AnyHelper() cli.AnyHelper {
	return cli.NewHelper[output](me)
}
