package upload

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nuggxyz/buildrc/internal/common"
	"github.com/nuggxyz/buildrc/internal/file"
	"github.com/nuggxyz/buildrc/internal/pipeline"
	"github.com/rs/zerolog"
	"github.com/spf13/afero"
)

type Handler struct {
	Package string `arg:"" help:"package to close"`
}

func (me *Handler) Run(ctx context.Context, prov common.Provider) (err error) {

	gz, err := pipeline.BuildrcArtifactsToReleaseAsTarGZDir.Path(ctx, prov.Pipeline())
	if err != nil {
		return err
	}

	rid, err := prov.Pipeline().GetFromEnv(ctx, "BUILDRC_RELEASE_ID", prov.FileSystem())
	if err != nil {
		return err
	}

	rel, err := prov.Release().GetReleaseByID(ctx, rid)
	if err != nil {
		return err
	}

	if err = afero.Walk(prov.FileSystem(), gz, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == gz {
			return nil
		}

		fle, err := file.Targz(ctx, prov.FileSystem(), path)
		if err != nil {
			return err
		}

		defer fle.Close()

		err = prov.Release().UploadReleaseArtifact(ctx, rel, filepath.Base(fle.Name()), fle)
		if err != nil {
			return fmt.Errorf("error uploading archive: %v", err)
		}

		zerolog.Ctx(ctx).Debug().Msgf("uploaded file %s", filepath.Base(fle.Name()))

		return nil
	}); err != nil {
		return err
	}

	sha, err := pipeline.BuildrcArtifactsToReleaseAsSha256Dir.Path(ctx, prov.Pipeline())
	if err != nil {
		return err
	}

	if err = afero.Walk(prov.FileSystem(), sha, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == sha {
			return nil
		}

		fle, err := file.Sha256(ctx, prov.FileSystem(), path)
		if err != nil {
			return err
		}

		defer fle.Close()

		err = prov.Release().UploadReleaseArtifact(ctx, rel, filepath.Base(fle.Name()), fle)
		if err != nil {
			return fmt.Errorf("error uploading archive: %v", err)
		}

		zerolog.Ctx(ctx).Debug().Msgf("uploaded file %s", filepath.Base(fle.Name()))

		return nil
	}); err != nil {
		return err
	}

	zerolog.Ctx(ctx).Debug().Any("next", rel).Msg("tagged release")

	return nil
}
