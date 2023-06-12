package packages

import (
	"context"

	_buildrc "github.com/nuggxyz/buildrc/internal/buildrc"
	"github.com/nuggxyz/buildrc/internal/cli"
)

func (me *Handler) Invoke(ctx context.Context, r cli.ContentProvider) (out *output, err error) {

	buildrc, err := me.buildrchandler.Helper().Run(ctx, r)
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

	return &output{
		Golang:     gopkgs,
		GolangAlt1: gopkgs,
		Docker:     docker,
	}, nil

}
