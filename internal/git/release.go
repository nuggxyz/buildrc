package git

import (
	"context"

	"github.com/Masterminds/semver/v3"
	"github.com/spf13/afero"
)

type Release struct {
	ID         string
	CommitHash string
	Tag        string
	PR         *PullRequest
	Artifacts  []string
}

type ReleaseProvider interface {
	CreateRelease(ctx context.Context, g GitProvider, t *semver.Version) (*Release, error)
	UploadReleaseArtifact(ctx context.Context, r *Release, name string, file afero.File) error
	DownloadReleaseArtifact(ctx context.Context, r *Release, name string, filesystem afero.Fs) (afero.File, error)
	GetReleaseByCommit(ctx context.Context, ref string) (*Release, error)
	GetReleaseByTag(ctx context.Context, tag string) (*Release, error)
	TagRelease(ctx context.Context, r *Release, vers *semver.Version, commit string) (*Release, error)
}

func ReleaseAlreadyExists(ctx context.Context, prov ReleaseProvider, cmt GitProvider) (bool, error) {
	str, err := cmt.GetCurrentCommitHash(ctx)
	if err != nil {
		return false, err
	}

	rel, err := prov.GetReleaseByCommit(ctx, str)
	if err != nil {
		return false, err
	}

	return rel != nil && str == rel.CommitHash, nil
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

func GetCurrentRelease(ctx context.Context, prov ReleaseProvider, cmt GitProvider) (*Release, error) {
	str, err := cmt.GetCurrentCommitHash(ctx)
	if err != nil {
		return nil, err
	}

	rel, err := prov.GetReleaseByCommit(ctx, str)
	if err != nil {
		return nil, err
	}

	return rel, nil
}

func (me *Release) Semver() (*semver.Version, error) {
	return semver.NewVersion(me.Tag)
}
