package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v33/github"
)

type GithubAPI interface {
	ListTags(ctx context.Context, repo string) ([]*github.RepositoryTag, error)
	Pull(ctx context.Context, repo, branch string) error
	CreateRelease(ctx context.Context, repo, tag, body string) error
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
