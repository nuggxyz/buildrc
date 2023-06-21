package git

import (
	"context"

	"github.com/Masterminds/semver/v3"
	"github.com/rs/zerolog"
	"github.com/spf13/afero"
)

type Release struct {
	ID         string
	CommitHash string
	Tag        string
	PR         *PullRequest
	Artifacts  []string
	Draft      bool
}

type ReleaseProvider interface {
	CreateRelease(ctx context.Context, g GitProvider, t *semver.Version) (*Release, error)
	UploadReleaseArtifact(ctx context.Context, r *Release, name string, file afero.File) error
	DownloadReleaseArtifact(ctx context.Context, r *Release, name string, filesystem afero.Fs) (afero.File, error)
	GetReleaseByTag(ctx context.Context, tag string) (*Release, error)
	TagRelease(ctx context.Context, r *Release, vers *semver.Version, commit string) (*Release, error)
	ListRecentReleases(ctx context.Context, limit int) ([]*Release, error)
}

func ReleaseAlreadyExists(ctx context.Context, prov ReleaseProvider, gitp GitProvider) (bool, string, error) {

	current, err := gitp.GetCurrentCommitHash(ctx)
	if err != nil {
		return false, "", err
	}

	releases, err := prov.ListRecentReleases(ctx, 100)
	if err != nil {
		return false, "", err
	}

	for _, rel := range releases {
		if current == rel.CommitHash && !rel.Draft {
			zerolog.Ctx(ctx).Info().Str("tag", rel.Tag).Any("release", rel).Msg("release already exists")
			return true, rel.Tag, nil
		}
	}

	return false, "", nil
}

func CopyReleaseArtifacts(ctx context.Context, fromprov, toprov ReleaseProvider, from, to *Release) error {

	files := afero.NewMemMapFs()

	for _, artifact := range from.Artifacts {

		osf, err := fromprov.DownloadReleaseArtifact(ctx, from, artifact, files)
		if err != nil {
			return err
		}

		err = toprov.UploadReleaseArtifact(ctx, to, artifact, osf)
		if err != nil {
			return err
		}

	}

	return nil
}

func (me *Release) Semver() (*semver.Version, error) {
	return semver.NewVersion(me.Tag)
}
