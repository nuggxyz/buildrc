package restapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v53/github"
	"github.com/walteh/buildrc/pkg/git"
	"golang.org/x/oauth2"
)

type GithubClient struct {
	client   *github.Client
	repoName string
	orgName  string
}

func (me *GithubClient) RepoName() string {
	return me.repoName
}

func (me *GithubClient) OrgName() string {
	return me.orgName
}

func (me *GithubClient) Client() *github.Client {
	return me.client
}

type GithubRestApiTokenProvider interface {
	GithubRestApiToken() (string, error)
}

func NewGithubClient(ctx context.Context, gitp git.GitProvider, token GithubRestApiTokenProvider) (*GithubClient, error) {

	var err error

	cmd, err := gitp.GetLocalRepositoryMetadata(ctx)
	if err != nil {
		return nil, err
	}

	t, err := token.GithubRestApiToken()
	if err != nil {
		return nil, err
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: t})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// check that the token is valid
	_, _, err = client.Zen(ctx)
	if err != nil {
		return nil, err
	}

	return &GithubClient{client: client, repoName: cmd.Name, orgName: cmd.Owner}, nil
}

func ParseRepo(input string) (owner string, name string, err error) {
	if input == "" {
		return "", "", fmt.Errorf("repo is empty")
	}

	parts := strings.Split(input, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("repo is not in the format owner/repo")
	}

	return parts[0], parts[1], nil
}
