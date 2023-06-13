package packages

import (
	"context"

	"github.com/nuggxyz/buildrc/cmd/buildrc/load"
	"github.com/nuggxyz/buildrc/internal/buildrc"
	"github.com/nuggxyz/buildrc/internal/provider"
)

const (
	CommandID = "packages"
)

type Handler struct {
	File string `flag:"file" type:"file:" default:".buildrc"`
}

type Output struct {
	Golang     []*buildrc.Package `json:"golang" yaml:"golang"`
	GolangAlt1 []*buildrc.Package `json:"go" yaml:"go"`
	Docker     []*buildrc.Package `json:"docker" yaml:"docker"`
}

func (me *Handler) Run(ctx context.Context, cp provider.ContentProvider) (err error) {
	_, err = me.Load(ctx, cp)
	return err
}

func (me *Handler) Load(ctx context.Context, cp provider.ContentProvider) (out *Output, err error) {
	return provider.Wrap(CommandID, me.invoke)(ctx, cp)
}

func (me *Handler) invoke(ctx context.Context, r provider.ContentProvider) (out *Output, err error) {

	brc, err := load.NewHandler(me.File).Load(ctx, r)
	if err != nil {
		return nil, err
	}

	gopkgs := make([]*buildrc.Package, 0)
	docker := make([]*buildrc.Package, 0)

	for _, pkg := range brc.Packages {
		switch pkg.Language {
		case buildrc.PackageLanguageGo:
			gopkgs = append(gopkgs, pkg)
		case buildrc.PackageLanguageDocker:
			docker = append(docker, pkg)
		}

	}

	return &Output{
		Golang:     gopkgs,
		GolangAlt1: gopkgs,
		Docker:     docker,
	}, nil

}
