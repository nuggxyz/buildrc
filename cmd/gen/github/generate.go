package github

import (
	"context"

	"github.com/nuggxyz/buildrc/internal/cli"
)

type OUTPUT = any

var _ cli.Command[OUTPUT] = (*Handler)(nil)
var _ cli.CommandRunner = (*Handler)(nil)

func (me *Handler) ID() string {
	return "generate"
}

func NewHandler(ctx context.Context, outdir string) (*Handler, error) {
	h := &Handler{OutDir: outdir}

	err := h.Init(ctx)

	return h, err
}

func (me *Handler) Init(ctx context.Context) error {
	return nil
}

func (me *Handler) Helper() cli.CommandHelper[OUTPUT] {
	return cli.NewHelper[OUTPUT](me)
}

func (me *Handler) AnyHelper() cli.AnyHelper {
	return cli.NewHelper[OUTPUT](me)
}
