package load

import (
	"context"

	"github.com/nuggxyz/buildrc/internal/buildrc"
	"github.com/nuggxyz/buildrc/internal/cli"
)

type output = buildrc.BuildRC

var _ cli.CommandRunner = (*Handler)(nil)
var _ cli.Command[buildrc.BuildRC] = (*Handler)(nil)

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
	return "buildrc"
}

func (me *Handler) Helper() cli.CommandHelper[output] {
	return cli.NewHelper[output](me)
}

func (me *Handler) AnyHelper() cli.AnyHelper {
	return cli.NewHelper[output](me)
}
