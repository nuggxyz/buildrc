package packages

import (
	"context"

	"github.com/nuggxyz/buildrc/cmd/buildrc/load"
	"github.com/nuggxyz/buildrc/internal/buildrc"
	"github.com/nuggxyz/buildrc/internal/cli"
)

var _ cli.CommandRunner = (*Handler)(nil)
var _ cli.Command[output] = (*Handler)(nil)

type Handler struct {
	File string `arg:"file" type:"file:" required:"true"`

	buildrchandler *load.Handler
}

type packageByLanguage struct {
	Golang     []*buildrc.Package `json:"golang" yaml:"golang"`
	GolangAlt1 []*buildrc.Package `json:"go" yaml:"go"`
	Docker     []*buildrc.Package `json:"docker" yaml:"docker"`
}

type output = packageByLanguage

func (me *Handler) Init(ctx context.Context) error {

	brc, err := load.NewHandler(ctx, me.File)
	if err != nil {
		return err
	}

	me.buildrchandler = brc

	return nil
}

func (me *Handler) ID() string {
	return "buildrc-package-by-language"
}

func NewHandler(file string) *Handler {
	return &Handler{File: file}
}

func (me *Handler) Helper() cli.CommandHelper[output] {
	return cli.NewHelper[output](me)
}

func (me *Handler) AnyHelper() cli.AnyHelper {
	return cli.NewHelper[output](me)
}
