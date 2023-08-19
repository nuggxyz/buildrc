package restapi

import (
	"context"

	"github.com/walteh/buildrc/pkg/git"
)

var _ git.RemoteRepositoryMetadataProvider = (*GithubClient)(nil)

func (me *GithubClient) GetRemoteRepositoryMetadata(ctx context.Context) (*git.RemoteRepositoryMetadata, error) {

	repo, _, err := me.client.Repositories.Get(ctx, me.OrgName(), me.RepoName())
	if err != nil {
		return nil, err
	}

	return &git.RemoteRepositoryMetadata{
		// Owner:       repo.GetOwner().GetLogin(),
		// Name:        repo.GetName(),
		Description: repo.GetDescription(),
		Homepage:    repo.GetHTMLURL(),
		License:     repo.GetLicense().GetSPDXID(),
	}, nil
}
