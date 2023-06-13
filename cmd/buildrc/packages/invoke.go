package packages

import (
	"context"

	"github.com/nuggxyz/buildrc/cmd/buildrc/load"
	"github.com/nuggxyz/buildrc/internal/buildrc"
	_buildrc "github.com/nuggxyz/buildrc/internal/buildrc"
	"github.com/nuggxyz/buildrc/internal/provider"
)

type Handler struct {
	File string `flag:"file" type:"file:" default:".buildrc"`
}

type packageByLanguage struct {
	Golang     []*buildrc.Package `json:"golang" yaml:"golang"`
	GolangAlt1 []*buildrc.Package `json:"go" yaml:"go"`
	Docker     []*buildrc.Package `json:"docker" yaml:"docker"`
}

func (me *Handler) Invoke(ctx context.Context, r provider.ContentProvider) (out *packageByLanguage, err error) {

	brc := load.NewHandler(me.File)

	buildrc, err := brc.Load(ctx, r)
	if err != nil {
		return nil, err
	}

	gopkgs := make([]*_buildrc.Package, 0)
	docker := make([]*_buildrc.Package, 0)

	for _, pkg := range buildrc.Packages {
		switch pkg.Language {
		case _buildrc.PackageLanguageGo:
			gopkgs = append(gopkgs, pkg)
		case _buildrc.PackageLanguageDocker:
			docker = append(docker, pkg)
		}

	}

	return &packageByLanguage{
		Golang:     gopkgs,
		GolangAlt1: gopkgs,
		Docker:     docker,
	}, nil

}
