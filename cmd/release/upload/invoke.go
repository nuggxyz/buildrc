package upload

import (
	"context"
	"fmt"
	"time"

	"github.com/nuggxyz/buildrc/cmd/release/setup"
	"github.com/nuggxyz/buildrc/internal/buildrc"
	"github.com/nuggxyz/buildrc/internal/common"
	"github.com/nuggxyz/buildrc/internal/git"
	"github.com/nuggxyz/buildrc/internal/pipeline"
	"github.com/spf13/afero"

	"github.com/rs/zerolog"
)

const (
	CommandID = "upload"
)

type Handler struct {
	File string `flag:"file" type:"file:" default:".buildrc"`
}

func (me *Handler) Run(ctx context.Context, prov common.Provider) (err error) {
	_, err = me.Build(ctx, prov)
	return err
}

func (me *Handler) Build(ctx context.Context, prov common.Provider) (out *any, err error) {
	return pipeline.Cache(ctx, CommandID, prov, me.build)
}

func (me *Handler) build(ctx context.Context, prov common.Provider) (out *any, err error) {

	yes, _, err := git.ReleaseAlreadyExists(ctx, prov.Release(), prov.Git())
	if err != nil {
		return nil, err
	}
	if yes {
		zerolog.Ctx(ctx).Info().Msg("build not required - release already exists")
		return nil, nil
	}

	err = me.run(ctx, prov)
	if err != nil {
		return nil, err
	}

	return nil, nil

}

func (me *Handler) run(ctx context.Context, prov common.Provider) error {

	su, err := setup.NewHandler("", "").Invoke(ctx, prov)
	if err != nil {
		return err
	}

	return buildrc.RunAllPackages(ctx, prov.Buildrc(), 10*time.Minute, func(ctx context.Context, pkg *buildrc.Package, arc buildrc.Platform) error {

		file, err := arc.OutputFile(pkg)
		if err != nil {
			return fmt.Errorf("error running upload with [%s:%s]: %v", arc.OS(), arc.Arch(), err)
		}

		cacher := pipeline.GetCacheFile(ctx, prov.Pipeline(), prov.FileSystem(), file)

		rel, err := prov.Release().GetReleaseByTag(ctx, su.UniqueReleaseTag)
		if err != nil {
			return fmt.Errorf("error getting current release: %v", err)
		}

		fs := afero.NewOsFs()

		r1, err := fs.Open(cacher.String() + ".tar.gz")
		if err != nil {
			return fmt.Errorf("error opening archive: %v", err)
		}

		err = prov.Release().UploadReleaseArtifact(ctx, rel, file, r1)
		if err != nil {
			return fmt.Errorf("error uploading archive: %v", err)
		}

		r2, err := fs.Open(cacher.String() + ".sha256")
		if err != nil {
			return fmt.Errorf("error opening checksum: %v", err)
		}

		err = prov.Release().UploadReleaseArtifact(ctx, rel, file+".sha256", r2)
		if err != nil {
			return fmt.Errorf("error uploading checksum: %v", err)
		}

		zerolog.Ctx(ctx).Debug().Msgf("uploaded checksum %s.sha256", file)

		return nil
	})

}
