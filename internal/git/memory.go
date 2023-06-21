package git

import (
	"context"

	"github.com/Masterminds/semver/v3"
	"github.com/spf13/afero"
)

/* /////////////////////////////////////////
	REPOMETADATA PROVIDER
///////////////////////////////////////// */

type memoryRepoMetadataProvider struct {
	cmd *RemoteRepositoryMetadata
}

func NewMemoryRepoMetadataProvider(cmd *RemoteRepositoryMetadata) RemoteRepositoryMetadataProvider {
	return &memoryRepoMetadataProvider{cmd: cmd}
}

func (me *memoryRepoMetadataProvider) GetRemoteRepositoryMetadata(ctx context.Context) (*RemoteRepositoryMetadata, error) {
	return me.cmd, nil
}

/* /////////////////////////////////////////
	PULLREQUEST PROVIDER
///////////////////////////////////////// */

type memoryPullRequestProvider struct {
	prs []*PullRequest
}

func NewMemoryPullRequestProvider(prs []*PullRequest) PullRequestProvider {
	return &memoryPullRequestProvider{prs: prs}
}

func (me *memoryPullRequestProvider) ListRecentPullRequests(ctx context.Context, head string) ([]*PullRequest, error) {
	return me.prs, nil
}

/* /////////////////////////////////////////
	RELEASE PROVIDER
///////////////////////////////////////// */

type memoryReleaseProvider struct {
	rels []*Release
}

func NewMemoryReleaseProvider(rels []*Release) ReleaseProvider {
	return &memoryReleaseProvider{rels: rels}
}

func (me *memoryReleaseProvider) CreateRelease(ctx context.Context, g GitProvider) (*Release, error) {
	r, err := g.GetCurrentCommitHash(ctx)
	if err != nil {
		return nil, err
	}

	return &Release{
		ID:         r,
		CommitHash: r,
		Artifacts:  []string{},
		PR:         nil,
		Tag:        "",
	}, nil
}

func (me *memoryReleaseProvider) UploadReleaseArtifact(ctx context.Context, r *Release, name string, file afero.File) error {
	r.Artifacts = append(r.Artifacts, name)
	return nil
}

func (me *memoryReleaseProvider) DownloadReleaseArtifact(ctx context.Context, r *Release, name string, filesystem afero.Fs) (afero.File, error) {
	return filesystem.Create(name)
}

func (me *memoryReleaseProvider) GetReleaseByCommit(ctx context.Context, ref string) (*Release, error) {
	for _, r := range me.rels {
		if r.CommitHash == ref {
			return r, nil
		}
	}
	return nil, nil
}

func (me *memoryReleaseProvider) GetReleaseByTag(ctx context.Context, tag string) (*Release, error) {
	for _, r := range me.rels {
		if r.Tag == tag {
			return r, nil
		}
	}
	return nil, nil
}

func (me *memoryReleaseProvider) TagRelease(ctx context.Context, r *Release, vers *semver.Version, commit string) (*Release, error) {
	r.Tag = vers.String()
	r.CommitHash = commit
	return r, nil
}
