package restapi

import (
	"context"

	"github.com/google/go-github/v53/github"
	"github.com/nuggxyz/buildrc/internal/git"
	"github.com/rs/zerolog"
)

var _ git.PullRequestProvider = (*GithubClient)(nil)

func (me *GithubClient) ListRecentPullRequests(ctx context.Context, head string) ([]*git.PullRequest, error) {

	opts := &github.PullRequestListOptions{
		State:       "all",
		Head:        head,
		Base:        "main",
		ListOptions: github.ListOptions{PerPage: 100},
		Sort:        "updated",
		Direction:   "desc",
	}

	pulls, res, err := me.client.PullRequests.List(ctx, me.OrgName(), me.RepoName(), opts)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Any("response", res).Msg("failed to list pull requests")
		return nil, err
	}

	if len(pulls) == 0 {
		return []*git.PullRequest{}, nil
	}

	resp := make([]*git.PullRequest, len(pulls))
	for i, pr := range pulls {
		resp[i] = &git.PullRequest{
			Number: pr.GetNumber(),
			Head:   pr.GetHead().GetRef(),
			Closed: pr.GetState() == "closed",
			Open:   pr.GetState() == "open",
			Merged: pr.GetMerged(),
		}
	}

	zerolog.Ctx(ctx).Debug().Int("total_found", len(pulls)).Any("args", opts).Msg("pull requests loaded from github")

	// If there's more than one matching PR (which is unusual), return the first one.
	return resp, nil
}
