package github

import (
	"context"

	"github.com/nuggxyz/buildrc/internal/provider"
)

type OUTPUT = any

var _ provider.Command[OUTPUT] = (*Handler)(nil)
var _ provider.CommandRunner = (*Handler)(nil)

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

func (me *Handler) Helper() provider.CommandHelper[OUTPUT] {
	return provider.NewHelper[OUTPUT](me)
}

func (me *Handler) AnyHelper() provider.AnyHelper {
	return provider.NewHelper[OUTPUT](me)
}
