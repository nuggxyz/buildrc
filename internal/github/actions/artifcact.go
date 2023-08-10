package actions

import (
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/nuggxyz/buildrc/internal/pipeline"
	"github.com/rs/zerolog"
	"github.com/spf13/afero"
)

func (me *GithubActionPipeline) UploadArtifact(ctx context.Context, fls afero.Fs, name string, fle afero.File) error {

	res, err := pipeline.ArtifactsDir(ctx, me, fls)
	if err != nil {
		return err
	}

	a, err := fls.Create(filepath.Join(res, name))
	if err != nil {
		return err
	}

	defer a.Close()

	_, err = io.Copy(a, fle)
	if err != nil {
		return err
	}

	zerolog.Ctx(ctx).Debug().Str("artifact", name).Str("location", filepath.Join(res, name)).Msg("artifact added to output dir to be picked up by github actions")

	return nil
}

func (me *GithubActionPipeline) DownloadArtifact(ctx context.Context, fls afero.Fs, name string) (afero.File, error) {

	// create a temp dir to download the artifact to
	tmp, err := pipeline.NewNamedTempDir(ctx, me, fls, "temporary-artifacts")
	if err != nil {
		return nil, err
	}

	// check if the proccessed file exists in the temp dir
	_, err = fls.Stat(filepath.Join(tmp.String(), name))
	if err == nil {
		// file exists
		zerolog.Ctx(ctx).Debug().Str("artifact", name).Str("location", filepath.Join(tmp.String(), name)).Msg("artifact already downloaded")
		return fls.Open(filepath.Join(tmp.String(), name))
	}

	// use gh run download to download the artifact
	ex := exec.CommandContext(ctx, "gh", "run", "download", "-n", name, "-d", tmp.String())

	err = ex.Run()
	if err != nil {
		return nil, err
	}

	// loop through the files in the temp dir and untar them
	// this is because gh run download will download the artifact as a tar file
	// and we need to untar it
	err = afero.Walk(fls, tmp.String(), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		zerolog.Ctx(ctx).Debug().Str("file", path).Msg("untaring file")

		if info.IsDir() {
			return nil
		}

		// untar the file
		exg := exec.CommandContext(ctx, "tar", "-xvf", path, "-C", tmp.String())

		err = exg.Run()
		if err != nil {
			return err
		}

		zerolog.Ctx(ctx).Debug().Str("file", path).Msg("untared file")

		// remove the tar file
		err = fls.Remove(path)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	zerolog.Ctx(ctx).Debug().Str("cmd", ex.String()).Msg("DownloadArtifact")

	downloaded, err := fls.Open(filepath.Join(tmp.String(), name))
	if err != nil {
		return nil, err
	}

	return downloaded, nil
}
