package upload

import (
	"context"
	"fmt"
	"time"

	"github.com/nuggxyz/buildrc/cmd/release/finalize"
	"github.com/nuggxyz/buildrc/internal/buildrc"
	"github.com/nuggxyz/buildrc/internal/common"
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

	// yes, tagg, err := git.ReleaseAlreadyExists(ctx, prov.Release(), prov.Git())
	// if err != nil {
	// 	return nil, err
	// }
	// if yes {
	// 	zerolog.Ctx(ctx).Info().Bool("release_aleady_exists", yes).Str("tag", tagg).Msg("build not required")
	// 	return nil, nil
	// }

	err = me.run(ctx, prov)
	if err != nil {
		return nil, err
	}

	return nil, nil

}

func (me *Handler) run(ctx context.Context, prov common.Provider) error {

	su, err := finalize.NewHandler().Invoke(ctx, prov)
	if err != nil {
		return err
	}

	rel, err := prov.Release().GetReleaseByID(ctx, su.ReleaseID)
	if err != nil {
		return err
	}

	return buildrc.RunAllPackages(ctx, prov.Buildrc(), 10*time.Minute, func(ctx context.Context, pkg *buildrc.Package, arc buildrc.Platform) error {

		file, err := arc.OutputFile(pkg)
		if err != nil {
			return fmt.Errorf("error running upload with [%s:%s]: %v", arc.OS(), arc.Arch(), err)
		}

		cacher := pipeline.GetCacheFile(ctx, prov.Pipeline(), prov.FileSystem(), file)

		fs := afero.NewOsFs()

		for _, arc := range []string{".tar.gz", ".sha256"} {
			yes, err := prov.Release().HasReleaseArtifact(ctx, rel, file+arc)
			if err != nil {
				return fmt.Errorf("error getting current release: %v", err)
			}

			if yes {
				err = prov.Release().DeleteReleaseArtifact(ctx, rel, file+arc)
				if err != nil {
					return fmt.Errorf("error deleting current release: %v", err)
				}
			}
			r1, err := fs.Open(cacher.String() + arc)
			if err != nil {
				return fmt.Errorf("error openifile + archive: %v", err)
			}

			err = prov.Release().UploadReleaseArtifact(ctx, rel, file+arc, r1)
			if err != nil {
				return fmt.Errorf("error uploading archive: %v", err)
			}

			zerolog.Ctx(ctx).Debug().Msgf("uploaded file %s", file+arc)
		}

		return nil
	})

}
