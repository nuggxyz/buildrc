package github

import (
	"context"
	"fmt"
	"os/exec"
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

func GetCurrentRepo() (string, error) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	url := strings.TrimSpace(string(output))
	url = strings.TrimSuffix(url, ".git")

	if strings.HasPrefix(url, "https://") {
		parts := strings.Split(url, "/")
		return parts[len(parts)-2] + "/" + parts[len(parts)-1], nil
	} else if strings.HasPrefix(url, "git@") {
		parts := strings.Split(url, ":")
		return parts[1], nil
	}

	return "", fmt.Errorf("unrecognized URL format: %s", url)
}
