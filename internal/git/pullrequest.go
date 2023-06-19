package git

import (
	"context"
	"fmt"
	"sort"
)

type PullRequest struct {
	Number int
	Head   string
	Closed bool
	Open   bool
	Merged bool
}

type PullRequestProvider interface {
	ListRecentPullRequests(ctx context.Context, head string) ([]*PullRequest, error)
}

func (me *PullRequest) PreReleaseTag() string {
	return fmt.Sprintf("pr.%d", me.Number)
}

func getLatestPullRequest(ctx context.Context, prprov PullRequestProvider, head string) (*PullRequest, error) {

	prs, err := prprov.ListRecentPullRequests(ctx, head)
	if err != nil {
		return nil, err
	}

	if len(prs) == 0 {
		return nil, nil
	}

	sort.Slice(prs, func(i, j int) bool {
		return prs[i].Number > prs[j].Number
	})

	for _, pr := range prs {
		if pr.Open || pr.Merged {
			return pr, nil
		}
	}

	return nil, fmt.Errorf("no open or merged PRs found")
}
