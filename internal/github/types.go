package github

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/google/go-github/v33/github"
)

type GithubAPI interface {
	ListTags(ctx context.Context, repo string) ([]*github.RepositoryTag, error)
	GetBranch(ctx context.Context, repo, branch string) (*github.Branch, error)
	PutRelease(ctx context.Context, repo, tag string, cb ReleaseCallback) (*github.RepositoryRelease, error)
	ReduceTagVersions(ctx context.Context, repo string, filter Reducer[semver.Version]) (*semver.Version, error)
	CountTagVersions(ctx context.Context, repo string, filter Counter[semver.Version]) (int, error)
}

type Reducer[T any] func(*T, *T) *T
type Counter[T any] func(*T) bool

type ReleaseCallback func(context.Context, *github.RepositoryRelease) *github.RepositoryRelease

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

func GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

func GetCurrentTag() (string, error) {
	cmd := exec.Command("git", "tag", "--points-at", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

func GetCurrentCommit() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}
