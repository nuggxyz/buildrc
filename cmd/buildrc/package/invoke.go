package packagecmd

import (
	"context"
	"fmt"

	"github.com/nuggxyz/buildrc/internal/buildrc"
	"github.com/nuggxyz/buildrc/internal/common"
	"github.com/nuggxyz/buildrc/internal/pipeline"
)

const (
	CommandID = "package"
)

type Handler struct {
	File string `flag:"file" type:"file:" default:".buildrc"`
	Name string `arg:"name" help:"The name of the package to load."`
}

func NewHandler(file string, name string) *Handler {
	return &Handler{File: file, Name: name}
}

func (me *Handler) Run(ctx context.Context, prov common.Provider) (err error) {
	_, err = me.CachedLoad(ctx, prov)
	return err
}

func (me *Handler) CachedLoad(ctx context.Context, prov common.Provider) (out *buildrc.Package, err error) {
	return pipeline.Cache(ctx, CommandID, prov, me.load)
}

func (me *Handler) load(ctx context.Context, prov common.Provider) (out *buildrc.Package, err error) {

	pkg, ok := prov.Buildrc().PackageByName()[me.Name]
	if !ok {
		return nil, fmt.Errorf("package %s not found", me.Name)
	}

	artifacts, err := pkg.ToArtifactCSV(pkg.Platforms)
	if err != nil {
		return nil, err
	}

	export := map[string]string{
		"docker_platforms_csv":   buildrc.StringsToCSV(pkg.DockerPlatforms),
		"platforms_csv":          buildrc.StringsToCSV(pkg.Platforms),
		"platform_artifacts_csv": artifacts,
	}

	err = pipeline.AddContentToEnv(ctx, prov.Pipeline(), prov.FileSystem(), CommandID, export)
	if err != nil {
		return nil, err
	}

	return pkg, nil

}
