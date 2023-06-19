package github

import (
	"context"

	"github.com/nuggxyz/buildrc/internal/git"
)

var _ git.RepositoryMetadataProvider = (*GithubClient)(nil)

func (me *GithubClient) GetRepositoryMetadata(ctx context.Context) (*git.RepositoryMetadata, error) {

	repo, _, err := me.client.Repositories.Get(ctx, me.OrgName(), me.RepoName())
	if err != nil {
		return nil, err
	}

	return &git.RepositoryMetadata{
		Owner:       repo.GetOwner().GetLogin(),
		Name:        repo.GetName(),
		Description: repo.GetDescription(),
		Homepage:    repo.GetHTMLURL(),
		License:     repo.GetLicense().GetSPDXID(),
	}, nil
}
