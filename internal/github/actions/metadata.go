package actions

import (
	"errors"
	"strings"

	"github.com/walteh/buildrc/pkg/git"
)

func (me *GithubActionPipeline) LocalRepositoryMetadata() (*git.LocalRepositoryMetadata, error) {

	owner := EnvVarGithubRepositoryOwner.Load()
	if owner == "" {
		return nil, errors.New("env var not set or empty: " + string(EnvVarGithubRepositoryOwner))
	}

	repo := EnvVarGithubRepository.Load()
	if repo == "" {
		return nil, errors.New("env var not set or empty: " + string(EnvVarGithubRepository))
	}

	repo = strings.TrimPrefix(repo, owner+"/")

	remote := "github.com/" + owner + "/" + repo

	return &git.LocalRepositoryMetadata{
		Owner:  owner,
		Name:   repo,
		Remote: remote,
	}, nil

}
