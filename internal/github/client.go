package github

import (
	"context"

	"github.com/google/go-github/v33/github"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
)

var _ GithubAPI = (*GithubClient)(nil)

type GithubClient struct {
	client *github.Client
}

func NewGithubClient(ctx context.Context, token string) (*GithubClient, error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// check that the token is valid
	_, _, err := client.Zen(ctx)
	if err != nil {
		return nil, err
	}

	return &GithubClient{client: client}, nil
}

func (me *GithubClient) ListTags(ctx context.Context, repo string) ([]*github.RepositoryTag, error) {
	owner, name, err := ParseRepo(repo)
	if err != nil {
		return nil, err
	}

	tags, _, err := me.client.Repositories.ListTags(ctx, owner, name, &github.ListOptions{
		PerPage: 64,
		Page:    0,
	})
	if err != nil {
		return nil, err
	}

	zerolog.Ctx(ctx).Debug().Any("github_tags", tags).Msg("tags loaded from github")

	return tags, nil
}

func (me *GithubClient) Pull(ctx context.Context, repo, branch string) error {
	owner, name, err := ParseRepo(repo)
	if err != nil {
		return err
	}

	_, _, err = me.client.Repositories.GetBranch(ctx, owner, name, branch)
	if err != nil {
		return err
	}

	return nil
}

func (me *GithubClient) CreateRelease(ctx context.Context, repo, tag, body string) error {
	owner, name, err := ParseRepo(repo)
	if err != nil {
		return err
	}

	_, _, err = me.client.Repositories.CreateRelease(ctx, owner, name, &github.RepositoryRelease{
		TagName:    &tag,
		Body:       &body,
		Prerelease: github.Bool(false),
	})
	if err != nil {
		return err
	}

	return nil
}
