package actions

import (
	"context"
	"io"
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

	tmp, err := pipeline.NewTempDir(ctx, me, fls)
	if err != nil {
		return nil, err
	}

	// use gh run download to download the artifact
	ex := exec.CommandContext(ctx, "gh", "run", "download", "-n", name, "-d", tmp.String())

	err = ex.Run()
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
