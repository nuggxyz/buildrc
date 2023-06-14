package github

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/Masterminds/semver/v3"
	"github.com/google/go-github/v33/github"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
)

var _ GithubAPI = (*GithubClient)(nil)

type GithubClient struct {
	client *github.Client
}

func NewGithubClient(ctx context.Context, token string) (*GithubClient, error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// check that the token is valid
	_, _, err := client.Zen(ctx)
	if err != nil {
		return nil, err
	}

	return &GithubClient{client: client}, nil
}

func (me *GithubClient) ListTags(ctx context.Context, repository string) ([]*github.RepositoryTag, error) {
	owner, repo, err := ParseRepo(repository)
	if err != nil {
		return nil, err
	}

	tags := []*github.RepositoryTag{}

	opts := &github.ListOptions{PerPage: 100}
	for {
		t, resp, err := me.client.Repositories.ListTags(ctx, owner, repo, opts)
		if err != nil {
			return nil, err
		}

		tags = append(tags, t...)

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	zerolog.Ctx(ctx).Debug().Any("github_tags", tags).Msg("tags loaded from github")

	return tags, nil
}

func (me *GithubClient) GetBranch(ctx context.Context, repo, branch string) (*github.Branch, error) {
	owner, name, err := ParseRepo(repo)
	if err != nil {
		return nil, err
	}

	b, _, err := me.client.Repositories.GetBranch(ctx, owner, name, branch)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (me *GithubClient) PutRelease(ctx context.Context, repo, tag string, cb ReleaseCallback) (*github.RepositoryRelease, error) {
	owner, name, err := ParseRepo(repo)
	if err != nil {
		return nil, err
	}

	// check if the release already exists
	rel, err := me.GetRelease(ctx, repo, tag)
	if err != nil {
		return nil, err
	}

	prevId := int64(0)
	if rel != nil {
		prevId = rel.GetID()
	}

	rel = cb(ctx, rel)

	if prevId == 0 {
		rel, _, err = me.client.Repositories.CreateRelease(ctx, owner, name, rel)
		if err != nil {
			return nil, err
		}
	} else {
		rel, _, err = me.client.Repositories.EditRelease(ctx, owner, name, prevId, rel)
		if err != nil {
			return nil, err
		}
	}

	return rel, nil
}

func (me *GithubClient) GetCurrentPullRequest(ctx context.Context) (*github.PullRequest, error) {
	repository, err := GetCurrentRepo()
	if err != nil {
		return nil, err
	}

	branch, err := GetCurrentBranch()
	if err != nil {
		return nil, err
	}

	return me.GetPullRequest(ctx, repository, branch)
}

func (me *GithubClient) GetRelease(ctx context.Context, repo, tag string) (*github.RepositoryRelease, error) {
	owner, name, err := ParseRepo(repo)
	if err != nil {
		return nil, err
	}

	r, _, err := me.client.Repositories.GetReleaseByTag(ctx, owner, name, tag)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (me *GithubClient) GetPullRequest(ctx context.Context, repository, branch string) (*github.PullRequest, error) {
	owner, repo, err := ParseRepo(repository)
	if err != nil {
		return nil, err
	}

	opts := &github.PullRequestListOptions{
		State:       "open",
		Head:        branch,
		Base:        "main",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	pulls, res, err := me.client.PullRequests.List(ctx, owner, repo, opts)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Any("response", res).Msg("failed to list pull requests")
		return nil, err
	}

	zerolog.Ctx(ctx).Debug().Any("pull_requests", pulls).Any("response", res).Any("args", opts).Msg("pull requests loaded from github")

	if len(pulls) == 0 {
		return nil, nil
	}

	// If there's more than one matching PR (which is unusual), return the first one.
	return pulls[0], nil
}

func (me *GithubClient) ReduceTagVersions(ctx context.Context, repo string, compare Reducer[semver.Version]) (*semver.Version, error) {

	tags, err := me.ListTags(ctx, repo)
	if err != nil {
		return nil, err
	}

	wrk := semver.New(0, 0, 0, "", "")

	// compare all tags that are not of the format "vX.Y.Z"

	for _, t := range tags {
		ver, err := semver.StrictNewVersion(t.GetName())
		if err != nil {
			continue
		}

		wrk = compare(wrk, ver)

	}

	return wrk, nil

}

func (me *GithubClient) CountTagVersions(ctx context.Context, repo string, compare Counter[semver.Version]) (int, error) {
	tags, err := me.ListTags(ctx, repo)
	if err != nil {
		return 0, err
	}

	wrk := 0

	// compare all tags that are not of the format "vX.Y.Z"

	for _, t := range tags {
		ver, err := semver.StrictNewVersion(t.GetName())
		if err != nil {
			continue
		}

		if compare(ver) {
			wrk++
		}
	}

	return wrk, nil
}

func (me *GithubClient) EnsurePullRequest(ctx context.Context, repo, branch string) (*github.PullRequest, error) {
	owner, name, err := ParseRepo(repo)
	if err != nil {
		return nil, err
	}

	next, abc := me.GetPullRequest(ctx, repo, branch)
	if abc != nil {
		return nil, abc
	}

	if next != nil {

		return next, nil
	}

	zerolog.Ctx(ctx).Debug().Msgf("Creating PR for %s", branch)

	title := fmt.Sprintf("Release %s", branch)
	body := "Automatically generated release PR. Please update."
	base := "main"
	head := branch
	var issue *int

	// if commit message has "(issue:xxx)", add that to the PR body
	commit, err := GetCurrentCommit()
	if err == nil {
		if len(commit) > 0 {
			re := regexp.MustCompile(`\(issue:(.+)\)`)
			matches := re.FindAllStringSubmatch(commit, -1)
			if len(matches) > 0 {
				iss, err := strconv.Atoi(matches[0][1])
				if err == nil {
					body = fmt.Sprintf("%s\n\nIssue: #%s", body, matches[0][1])
					issue = github.Int(iss)
				}
			}
		}
	}

	// create a new PR
	pr, res, err := me.client.PullRequests.Create(ctx, owner, name, &github.NewPullRequest{
		Title: github.String(title),
		Body:  github.String(body),
		Base:  github.String(base),
		Head:  github.String(head),
		Issue: issue,
		Draft: github.Bool(true),
	})

	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Any("response", res).Msgf("Failed to create PR: %s", res.Status)
		return nil, err
	}

	return pr, nil
}
