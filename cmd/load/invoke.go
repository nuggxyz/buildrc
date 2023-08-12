package load

import (
	"context"

	"github.com/nuggxyz/buildrc/internal/buildrc"
	"github.com/nuggxyz/buildrc/internal/common"
	"github.com/nuggxyz/buildrc/internal/git"
	"github.com/nuggxyz/buildrc/internal/pipeline"
)

const (
	CommandID = "load"
)

type Handler struct {
	File string `flag:"file" type:"file:" default:".buildrc"`
}

func NewHandler(file string) *Handler {
	return &Handler{File: file}
}

func (me *Handler) Run(ctx context.Context, prov common.Provider) (err error) {

	out, err := buildrc.Parse(ctx, me.File)
	if err != nil {
		return err
	}

	err = pipeline.SetupEnvDirs(ctx, prov.Pipeline(), prov.FileSystem())
	if err != nil {
		return err
	}

	targetSemver, err := git.CalculateNextPreReleaseTag(ctx, prov.Buildrc(), prov.Git(), prov.PR())
	if err != nil {
		return err
	}

	export := map[string]string{
		"BUILDRC_PACKAGES_ARRAY_JSON":    out.PackagesNamesArrayJSON(),
		"BUILDRC_PACKAGES_ON_ARRAY_JSON": out.PackagesOnArrayJSON(),
		"BUILDRC_TAG":                    targetSemver.String(),
	}

	err = pipeline.AddContentToEnv(ctx, prov.Pipeline(), prov.FileSystem(), CommandID, export)

	if err != nil {
		return err
	}

	return nil
}
