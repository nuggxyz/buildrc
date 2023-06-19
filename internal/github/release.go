package github

import (
	"context"

	"github.com/Masterminds/semver/v3"
	"github.com/nuggxyz/buildrc/internal/git"
)

var _ git.ReleaseProvider = (*GithubClient)(nil)

func (me *GithubClient) GetLatestRelease(ctx context.Context) (string, error) {
	return "", nil
}

func (me *GithubClient) CreateRelease(ctx context.Context) error {
	return nil
}

func (me *GithubClient) UploadReleaseArtifact(ctx context.Context, r *git.Release, artifactPath string) error {
	return nil
}

func (me *GithubClient) DownloadReleaseArtifact(ctx context.Context, r *git.Release, artifactPath string) error {
	return nil
}

func (me *GithubClient) GetReleaseByCommit(ctx context.Context, ref string) (*git.Release, error) {
	return nil, nil
}

func (me *GithubClient) MakeReleaseLive(ctx context.Context, r *git.Release) (*semver.Version, error) {
	return nil, nil
}
