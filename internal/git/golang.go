package git

import (
	"context"
	"fmt"
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

var _ GitProvider = (*GitGoGitProvider)(nil)

type GitGoGitProvider struct {
}

func (me *GitGoGitProvider) GetContentHash(ctx context.Context, sha string) (string, error) {
	panic("implement me")
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

	tags, err := repo.TagObjects()
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %v", err)
	}

	var versions []*semver.Version

	if err = tags.ForEach(func(ref *object.Tag) error {

		// Attempt to parse each tag as a semver version
		ver, err := semver.NewVersion(ref.Name)
		if err != nil {
			// Skip this tag and move to the next one
			return nil
		}
		versions = append(versions, ver)

		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to iterate over tags: %v", err)
	}

	// Return error if no semver tags found
	if len(versions) == 0 {
		return nil, fmt.Errorf("no semver tags found from ref '%s'", ref)
	}

	// Sort the versions in descending order
	sort.Sort(sort.Reverse(semver.Collection(versions)))

	// Return the latest version
	return versions[0], nil
}
