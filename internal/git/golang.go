package git

import (
	"context"
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/rs/zerolog"
)

var _ GitProvider = (*GitGoGitProvider)(nil)

type GitGoGitProvider struct {
}

func NewGitGoGitProvider() GitProvider {
	return &GitGoGitProvider{}
}

func (me *GitGoGitProvider) GetContentHash(ctx context.Context, sha string) (string, error) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return "", err
	}

	// Parse the provided SHA to a Hash
	hash := plumbing.NewHash(sha)

	// Get the commit object from the hash
	commit, err := repo.CommitObject(hash)
	if err != nil {
		return "", err
	}

	// Get the tree from the commit
	tree, err := commit.Tree()
	if err != nil {
		return "", err
	}

	// The hash of the tree can be used as a hash of the contents
	return tree.Hash.String(), nil
}

func (me *GitGoGitProvider) GetCurrentCommitHash(ctx context.Context) (string, error) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return "", err
	}

	headRef, err := repo.Head()
	if err != nil {
		return "", err
	}

	return headRef.Hash().String(), nil
}

func (me *GitGoGitProvider) GetCurrentBranch(ctx context.Context) (string, error) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return "", err
	}

	headRef, err := repo.Head()
	if err != nil {
		return "", err
	}

	return headRef.Name().Short(), nil
}

func (me *GitGoGitProvider) GetLatestSemverTagFromRef(ctx context.Context, ref string) (*semver.Version, error) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return nil, err
	}

	tags, err := repo.Tags()
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %v", err)
	}

	var latestSemver *semver.Version
	err = tags.ForEach(func(refObj *plumbing.Reference) error {
		// resolve tag commit
		tagObj, err := repo.TagObject(refObj.Hash())
		if err != nil {
			// could be lightweight tag, ignore and continue
			return nil
		}

		// check if the commit SHA from the tag matches the provided ref
		if tagObj.Target.String() != ref {
			// not the tag for the provided ref
			return nil
		}

		// Attempt to parse each tag as a semver version
		ver, err := semver.NewVersion(refObj.Name().Short())
		if err != nil {
			// Skip this tag and move to the next one
			return nil
		}

		// keep track of the maximum version encountered
		if latestSemver == nil || ver.GreaterThan(latestSemver) {
			latestSemver = ver
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to iterate over tags: %v", err)
	}

	// Return error if no semver tags found
	if latestSemver == nil {
		zerolog.Ctx(ctx).Warn().Any("tags", tags).Msgf("no semver tags found from ref '%s'", ref)
		return nil, fmt.Errorf("no semver tags found from ref '%s'", ref)
	}

	// Return the latest version
	return latestSemver, nil
}

func (me *GitGoGitProvider) GetLocalRepositoryMetadata(ctx context.Context) (*LocalRepositoryMetadata, error) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	remotes, err := repo.Remotes()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if len(remotes) == 0 {
		fmt.Println("No remotes found")
		return nil, fmt.Errorf("no remotes found")
	}

	remoteURL := remotes[0].Config().URLs[0]
	splitURL := strings.Split(remoteURL, "/")
	repoNameWithGit := splitURL[len(splitURL)-1]

	// Remove .git from repo name
	repoName := strings.TrimSuffix(repoNameWithGit, ".git")

	return &LocalRepositoryMetadata{
		Owner:  strings.Join(splitURL[len(splitURL)-2:len(splitURL)-1], "/"),
		Name:   repoName,
		Remote: remoteURL,
	}, nil
}
