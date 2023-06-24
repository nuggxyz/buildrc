package git

import (
	"context"
	"fmt"
	"sort"

	"github.com/k0kubun/pp/v3"
	"github.com/rs/zerolog"
)

type PullRequest struct {
	Number int
	Head   string
	Closed bool
	Open   bool
}

type PullRequestProvider interface {
	ListRecentPullRequests(ctx context.Context, head string) ([]*PullRequest, error)
}

func (me *PullRequest) PreReleaseTag() string {
	return fmt.Sprintf("pr.%d", me.Number)
}

func getLatestPullRequestForRef(ctx context.Context, prprov PullRequestProvider, head string) (*PullRequest, error) {

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
		if pr.Open || pr.Closed {
			return pr, nil
		}
	}

	return nil, fmt.Errorf("no open or merged PRs found")
}

func getLatestMergedPullRequestThatHasAMatchingContentHash(ctx context.Context, prprov PullRequestProvider, git GitProvider) (*PullRequest, error) {

	mycontenthash, err := git.GetContentHashFromRef(ctx, "HEAD")
	if err != nil {
		return nil, err
	}

	// make sure we are on main

	branch, err := git.GetCurrentBranchFromRef(ctx, "HEAD")
	if err != nil {
		return nil, err
	}

	if branch != "main" {
		return nil, fmt.Errorf("not on main branch")
	}

	prs, err := prprov.ListRecentPullRequests(ctx, "main")
	if err != nil {
		return nil, err
	}

	pp.Println(prs)

	for _, pr := range prs {
		if !pr.Closed {
			continue
		}

		prcontenthash, err := git.GetContentHashFromRef(ctx, pr.Head)
		if err != nil {
			return nil, err
		}

		zerolog.Ctx(ctx).Debug().Str("prcontenthash", prcontenthash).Str("mycontenthash", mycontenthash).Msg("comparing content hashes")

		if prcontenthash == mycontenthash {
			return pr, nil
		}
	}

	return nil, ErrNoMatchingPR
}
