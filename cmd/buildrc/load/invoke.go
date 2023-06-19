package load

import (
	"context"

	"github.com/nuggxyz/buildrc/internal/buildrc"
	"github.com/nuggxyz/buildrc/internal/common"
	"github.com/nuggxyz/buildrc/internal/pipeline"
)

const (
	CommandID = "load"
)

type Handler struct {
	File string `flag:"file" type:"file:" default:".buildrc"`
}

type Output struct {
	*buildrc.Buildrc
}

func NewHandler(file string) *Handler {
	return &Handler{File: file}
}

func (me *Handler) Run(ctx context.Context, prov common.Provider) (err error) {
	_, err = me.Load(ctx, prov)
	return err
}

func (me *Handler) Load(ctx context.Context, prov common.Provider) (out *buildrc.Buildrc, err error) {
	return pipeline.Cache(ctx, CommandID, prov, me.load)
}

func (me *Handler) load(ctx context.Context, prov common.Provider) (out *buildrc.Buildrc, err error) {

	out, err = buildrc.Parse(ctx, me.File)
	if err != nil {
		return nil, err
	}

	err = pipeline.AddContentToEnv(ctx, prov.Pipeline(), CommandID, map[string]string{
		"package_names_array": out.PackagesNamesArrayJSON(),
	})

	if err != nil {
		return nil, err
	}

	return out, nil
}

// func (me *Output) Express() interface{} {
// 	return me.BuildRC
// }
