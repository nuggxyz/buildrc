package load

import (
	"context"

	"github.com/nuggxyz/buildrc/internal/buildrc"
	"github.com/nuggxyz/buildrc/internal/provider"
)

const (
	CommandID = "load"
)

type Handler struct {
	File string `flag:"file" type:"file:" default:".buildrc"`
}

type Output struct {
	*buildrc.BuildRC
}

func NewHandler(file string) *Handler {
	return &Handler{File: file}
}

func (me *Handler) Run(ctx context.Context, cp provider.ContentProvider) (err error) {
	_, err = me.Load(ctx, cp)
	return err
}

func (me *Handler) Load(ctx context.Context, cp provider.ContentProvider) (out *buildrc.BuildRC, err error) {
	return provider.Wrap(CommandID, me.load)(ctx, cp)
}

func (me *Handler) load(ctx context.Context, r provider.ContentProvider) (out *buildrc.BuildRC, err error) {

	out, err = buildrc.Parse(ctx, me.File)
	if err != nil {
		return nil, err
	}

	r.Express(ctx, CommandID, map[string]string{
		"package_names_array": out.PackagesNamesArrayJSON(),
	})

	return
}

func (me *Output) Express() interface{} {
	return me.BuildRC
}
