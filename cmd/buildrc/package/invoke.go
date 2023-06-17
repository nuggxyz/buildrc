package packagecmd

import (
	"context"
	"fmt"

	"github.com/nuggxyz/buildrc/cmd/buildrc/load"
	"github.com/nuggxyz/buildrc/internal/buildrc"
	"github.com/nuggxyz/buildrc/internal/provider"
)

const (
	CommandID = "package"
)

type Handler struct {
	File string `flag:"file" type:"file:" default:".buildrc"`
	Name string `arg:"name" help:"The name of the package to load."`
}

func (me *Handler) Run(ctx context.Context, cp provider.ContentProvider) (err error) {
	_, err = me.Load(ctx, cp)
	return err
}

func NewHandler(file string, name string) *Handler {
	return &Handler{File: file, Name: name}
}

func (me *Handler) Load(ctx context.Context, cp provider.ContentProvider) (out *buildrc.Package, err error) {
	return provider.Wrap(CommandID, me.load)(ctx, cp)
}

func (me *Handler) load(ctx context.Context, r provider.ContentProvider) (out *buildrc.Package, err error) {

	brc, err := load.NewHandler(me.File).Load(ctx, r)
	if err != nil {
		return nil, err
	}

	pkg, ok := brc.PackageByName()[me.Name]
	if !ok {
		return nil, fmt.Errorf("package %s not found", me.Name)
	}

	err = r.Express(ctx, CommandID, pkg.UsesMap())
	if err != nil {
		return nil, err
	}

	artifacts, err := pkg.ToArtifactCSV(pkg.Platforms)
	if err != nil {
		return nil, err
	}

	err = r.Express(ctx, CommandID, map[string]string{
		"docker_platforms_csv":   buildrc.StringsToCSV(pkg.DockerPlatforms),
		"platforms_csv":          buildrc.StringsToCSV(pkg.Platforms),
		"platform_artifacts_csv": artifacts,
	})
	if err != nil {
		return nil, err
	}

	return pkg, nil

}
