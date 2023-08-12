package actions

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/nuggxyz/buildrc/internal/file"
	"github.com/nuggxyz/buildrc/internal/pipeline"
	"github.com/rs/zerolog"
	"github.com/spf13/afero"
)

func (me *GithubActionPipeline) UploadArtifact(ctx context.Context, fls afero.Fs, name string, fle afero.File) error {

	res, err := pipeline.ArtifactsToUplaodDir(ctx, me, fls)
	if err != nil {
		return err
	}

	fileName := filepath.Base(name)

	a, err := fls.Create(filepath.Join(res, fileName))
	if err != nil {
		return err
	}

	defer a.Close()

	_, err = io.Copy(a, fle)
	if err != nil {
		return err
	}

	zerolog.Ctx(ctx).Debug().Str("artifact", fileName).Str("location", filepath.Join(res, fileName)).Msg("artifact added to output dir to be picked up by github actions")

	return nil
}

func (me *GithubActionPipeline) DownloadArtifactLegacy(ctx context.Context, fls afero.Fs, name string) (afero.File, error) {
	fle := pipeline.GetNamedCacheFile(ctx, me, fls, name)

	// check if the file exists
	fi, err := fls.Open(fle.String())
	if err != nil {
		zerolog.Ctx(ctx).Debug().Str("artifact", name).Str("location", fle.String()).Msg("artifact not found in cache")
		// file exists
		// zerolog.Ctx(ctx).Debug().Str("artifact", name).Str("location", fle.String()).Msg("artifact already downloaded")
		return nil, err
	}

	return fi, nil

}

func (me *GithubActionPipeline) DownloadArtifact(ctx context.Context, fls afero.Fs, name string) (afero.File, error) {

	// fle := pipeline.GetNamedCacheFile(ctx, me, fls, name)

	// // check if the file exists
	// fi, err := fls.Open(fle.String())
	// if err != nil {
	// 	zerolog.Ctx(ctx).Debug().Str("artifact", name).Str("location", fle.String()).Msg("artifact not found in cache")
	// 	// file exists
	// 	// zerolog.Ctx(ctx).Debug().Str("artifact", name).Str("location", fle.String()).Msg("artifact already downloaded")
	// 	return nil, err
	// }

	dir, err := pipeline.ArtifactsDir(ctx, me, fls)
	if err != nil {
		return nil, err
	}

	f, err := afero.ReadDir(fls, dir)
	if err != nil {
		return nil, err
	}

	if len(f) == 0 {

		zerolog.Ctx(ctx).Debug().Str("artifact", name).Str("location", dir).Msg("artifact not found in cache")

		runid, err := me.RunId(ctx)
		if err != nil {
			return nil, err
		}

		// use gh run download to download the artifact
		ex := exec.CommandContext(ctx, "gh", "run", "download", fmt.Sprintf("%d", runid), "-d", dir)

		err = ex.Run()
		if err != nil {
			return nil, err
		}

		err = afero.Walk(fls, dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			zerolog.Ctx(ctx).Debug().Str("path", path).Msg("found file")

			if info.IsDir() {
				return nil
			}

			if filepath.Ext(path) != ".tar.gz" {
				return nil
			}

			_, err = file.Untargz(ctx, fls, path)
			if err != nil {
				return err
			}

			zerolog.Ctx(ctx).Debug().Str("path", path).Msg("untared file")

			return nil
		})

		if err != nil {
			return nil, err
		}

	}

	downloaded, err := fls.Open(filepath.Join(dir, name))
	if err != nil {
		return nil, err
	}

	return downloaded, nil
}
