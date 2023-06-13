package load

import (
	"context"

	"github.com/nuggxyz/buildrc/internal/buildrc"
	"github.com/nuggxyz/buildrc/internal/provider"
)

type Handler struct {
	File string `arg:"file" type:"file:" required:"true"`
}

func NewHandler(file string) *Handler {
	return &Handler{File: file}
}

func (me *Handler) Load(ctx context.Context, r provider.ContentProvider) (out *buildrc.BuildRC, err error) {

	out, err = buildrc.Parse(ctx, me.File)
	if err != nil {
		return nil, err
	}

	return
}

func (me *Handler) Express(ctx context.Context, out *buildrc.BuildRC) (map[string]string, error) {
	if len(out.Packages) == 0 {
		return map[string]string{
			"version": out.Version.String(),
		}, nil
	}
	return map[string]string{
		"version":          out.Version.String(),
		"dockerfile":       out.Packages[0].Dockerfile,
		"entry":            out.Packages[0].Entry,
		"platforms":        buildrc.StringsToCSV(out.Packages[0].Platforms),
		"docker_platforms": buildrc.StringsToCSV(out.Packages[0].DockerPlatforms),
		"artifacts":        (out.Packages[0].ToArtifactCSV(out.Packages[0].Platforms)),
	}, nil
}
