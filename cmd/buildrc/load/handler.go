package load

import (
	"context"

	"github.com/nuggxyz/buildrc/internal/buildrc"
	"github.com/nuggxyz/buildrc/internal/provider"
)

type output = buildrc.BuildRC

var _ provider.CommandRunner = (*Handler)(nil)
var _ provider.Command[buildrc.BuildRC] = (*Handler)(nil)

type Handler struct {
	File string `arg:"file" type:"file:" required:"true"`
}

func NewHandler(ctx context.Context, file string) (*Handler, error) {
	return &Handler{File: file}, nil
}

func (me *Handler) Init(ctx context.Context) error {
	return nil
}

func (me *Handler) ID() string {
	return "load"
}

func (me *Handler) Helper() provider.CommandHelper[output] {
	return provider.NewHelper[output](me)
}

func (me *Handler) AnyHelper() provider.AnyHelper {
	return provider.NewHelper[output](me)
}
