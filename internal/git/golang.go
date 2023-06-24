package git

import (
	"context"
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/rs/zerolog"
)

var _ GitProvider = (*GitGoGitProvider)(nil)

type GitGoGitProvider struct {
}

func NewGitGoGitProvider() GitProvider {
	return &GitGoGitProvider{}
}

func (me *GitGoGitProvider) GetContentHashFromRef(ctx context.Context, ref string) (string, error) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return "", err
	}

	commit, _, err := me.getCommitFromRef(ctx, repo, ref)
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

func (me *GitGoGitProvider) GetCurrentCommitFromRef(ctx context.Context, ref string) (string, error) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return "", err
	}

	commit, _, err := me.getCommitFromRef(ctx, repo, ref)
	if err != nil {
		return "", err
	}

	return commit.Hash.String(), nil
}

func (me *GitGoGitProvider) GetCurrentBranchFromRef(ctx context.Context, ref string) (string, error) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return "", err
	}

	_, reffer, err := me.getCommitFromRef(ctx, repo, ref)
	if err != nil {
		return "", err
	}

	return reffer.Name().Short(), nil
}

func (me *GitGoGitProvider) getCommitFromRef(ctx context.Context, repo *git.Repository, ref string) (*object.Commit, *plumbing.Reference, error) {

	var refname plumbing.ReferenceName

	switch ref {
	case "HEAD":
		refname = plumbing.HEAD
	case "master":
		refname = plumbing.Master
	case "main":
		refname = plumbing.Main
	default:
		refname = plumbing.ReferenceName(ref)
	}

	zerolog.Ctx(ctx).Debug().Str("ref", ref).Str("refname", refname.String()).Msg("resolving ref")

	resolved, err := repo.Reference(refname, true)
	if err != nil {
		resolved, err = repo.Reference(plumbing.ReferenceName(strings.Replace(string(refname), "heads", "remotes/origin", 1)), true)
		if err != nil {
			return nil, nil, ErrRefNotFound
		}
	}

	commit, err := repo.CommitObject(resolved.Hash())
	if err != nil {
		return nil, nil, err
	}

	return commit, resolved, nil
}

func (me *GitGoGitProvider) GetLatestSemverTagFromRef(ctx context.Context, ref string) (*semver.Version, error) {

	repo, err := git.PlainOpen(".")
	if err != nil {
		return nil, err
	}

	commit, reffer, err := me.getCommitFromRef(ctx, repo, ref)
	if err != nil {
		return nil, err
	}

	zerolog.Ctx(ctx).Debug().Str("ref", ref).Msg("searching for semver logs")

	var latestSemver *semver.Version

	tagz := make(map[string]string)

	for commit != nil {

		tags, err := repo.Tags()
		if err != nil {
			break
		}
		defer tags.Close()
		err = tags.ForEach(func(refr *plumbing.Reference) error {
			tagCommit, err := repo.CommitObject(refr.Hash())
			if err != nil {
				return nil
			}

			if commit.Hash.String() == tagCommit.Hash.String() {
				tagz[refr.Name().Short()] = tagCommit.Hash.String()
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("failed to iterate over tags: %v", err)
		}

		if reffer.Name().IsTag() {
			break
		}

		if commit.NumParents() > 0 {
			commit, err = commit.Parents().Next()
			if err != nil {
				break
			}
		} else {
			break
		}
	}
	if err != nil {
		return nil, fmt.Errorf("failed to iterate over tags: %v", err)
	}

	for tag := range tagz {
		v, err := semver.NewVersion(tag)
		if err != nil {
			zerolog.Ctx(ctx).Warn().Err(err).Str("tag", tag).Msg("failed to parse semver tag")
			continue
		}

		if latestSemver == nil || v.GreaterThan(latestSemver) {
			latestSemver = v
		}
	}

	// Return error if no semver tags found
	if latestSemver == nil {
		zerolog.Ctx(ctx).Warn().Any("tags", tagz).Msgf("no semver tags found from ref '%s'", ref)
		return nil, fmt.Errorf("no semver tags found from ref '%s'", ref)
	}

	zerolog.Ctx(ctx).Debug().Str("semver", latestSemver.String()).Msgf("latest semver tag from ref '%s'", ref)
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

func (me *GitGoGitProvider) GetCurrentShortHashFromRef(ctx context.Context, ref string) (string, error) {
	commitHash, err := me.GetCurrentCommitFromRef(ctx, ref)
	if err != nil {
		return "", err
	}

	return commitHash[:7], nil
}
