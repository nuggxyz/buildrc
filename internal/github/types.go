package github

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/google/go-github/v53/github"
	"github.com/nuggxyz/buildrc/internal/env"
)

type GithubAPI interface {
	ListTags(ctx context.Context) ([]*github.RepositoryTag, error)
	GetBranch(ctx context.Context, branch string) (*github.Branch, error)
	ReduceTagVersions(ctx context.Context, filter Reducer[semver.Version]) (*semver.Version, error)
	CountTagVersions(ctx context.Context, filter Counter[semver.Version]) (int, error)
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

func GetGithubTokenFromEnv(ctx context.Context) (string, error) {
	tkn := env.GetOrEmpty("GITHUB_TOKEN")
	if tkn == "" {
		return "", fmt.Errorf("no access token found in env")
	} else {

		return tkn, nil
	}

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

func GetCurrentCommitSha() (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

func GetCurrentShortCommitSha() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

func GetCurrentCommitTags() ([]string, error) {
	cmd := exec.Command("git", "tag", "--points-at", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return []string{}, err
	}

	return strings.Split(strings.TrimSpace(string(output)), "\n"), nil
}

func GetNameForThisBuildrcCommitTagPrefix() (string, error) {
	sha, err := GetCurrentShortCommitSha()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("buildrc-%s", sha), nil
}

func GetNameForThisBuildrcCommitTag() (string, error) {
	sha, err := GetNameForThisBuildrcCommitTagPrefix()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s-%s", sha, time.Now().Format("20060102")), nil
}

func IsAlreadyTaggedByBuildRc() (bool, error) {
	tags, err := GetCurrentCommitTags()
	if err != nil {
		return false, err
	}

	brc, err := GetNameForThisBuildrcCommitTagPrefix()
	if err != nil {
		return false, err
	}

	for _, tag := range tags {
		if strings.HasPrefix(tag, brc) {
			return true, nil
		}
	}

	return false, nil
}

func ArtifactListFromFileNames(cmt *github.Commit, names []string) []*github.ReleaseAsset {
	assets := []*github.ReleaseAsset{}
	for _, name := range names {
		assets = append(assets, &github.ReleaseAsset{
			Name: &name,
			Uploader: &github.User{
				Email: github.String(cmt.GetAuthor().GetEmail()),
				Name:  github.String(cmt.GetAuthor().GetName()),
			},
		})
	}

	return assets
}

// FindFiles walks the path p and returns a list of files that have one of the provided extensions.
// It returns an error if it was unable to access the path p or any of the subdirectories.
func FindFiles(p string, exts ...string) ([]string, error) {
	var files []string
	err := filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			for _, targetExt := range exts {
				if ext == targetExt {
					files = append(files, path)
					break
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func CommitHasBuildrcReleaseTag(cmt *github.Commit) bool {

	return false
}
