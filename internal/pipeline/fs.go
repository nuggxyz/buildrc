package pipeline

import (
	"context"
	"path/filepath"

	"github.com/spf13/afero"
)

type BuildrcDir string

const (
	BuildrcArtifactsToReleaseAsTarGZDir  BuildrcDir = "buildrc_release_as_tar_gz"
	BuildrcArtifactsToReleaseAsSha256Dir BuildrcDir = "buildrc_release_as_sha256"
)

func SetupEnvDirs(ctx context.Context, pip Pipeline, fs afero.Fs) error {
	for _, d := range []BuildrcDir{
		BuildrcArtifactsToReleaseAsTarGZDir,
		BuildrcArtifactsToReleaseAsSha256Dir,
	} {
		dir, err := d.Path(ctx, pip)
		if err != nil {
			return err
		}

		if err := fs.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return nil
}

func (b BuildrcDir) String() string {
	return string(b)
}

func (b BuildrcDir) Path(ctx context.Context, pip Pipeline) (string, error) {
	s, err := pip.RootDir(ctx)
	if err != nil {
		return "", err
	}

	return filepath.Join(s, string(b)), nil
}
