package git

import (
	"context"

	"github.com/Masterminds/semver/v3"
)

type Release struct {
	CommitHash string
	Tag        string
	PR         *PullRequest
	Artifacts  []string
}

type ReleaseProvider interface {
	CreateRelease(ctx context.Context) error
	UploadReleaseArtifact(ctx context.Context, r *Release, artifactPath string) error
	DownloadReleaseArtifact(ctx context.Context, r *Release, artifactPath string) error
	GetReleaseByCommit(ctx context.Context, ref string) (*Release, error)
	MakeReleaseLive(ctx context.Context, r *Release) (*semver.Version, error)
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

	return str == rel.CommitHash, nil
}

func CopyReleaseArtifacts(ctx context.Context, fromprov, toprov ReleaseProvider, from, to *Release) error {
	for _, artifact := range from.Artifacts {
		err := fromprov.DownloadReleaseArtifact(ctx, from, artifact)
		if err != nil {
			return err
		}

		err = toprov.UploadReleaseArtifact(ctx, to, artifact)
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
