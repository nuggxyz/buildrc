package github

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

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

	tagstrs := []string{}
	for _, t := range tags {
		tagstrs = append(tagstrs, t.GetName())
	}

	zerolog.Ctx(ctx).Debug().Any("github_tags", tagstrs).Msg("tags loaded from github")

	return tags, nil
}

func (me *GithubClient) GetBranch(ctx context.Context, repo, branch string) (*github.Branch, error) {
	owner, name, err := ParseRepo(repo)
	if err != nil {
		return nil, err
	}

	zerolog.Ctx(ctx).Debug().Str("owner", owner).Str("name", name).Str("branch", branch).Msg("get branch")

	b, _, err := me.client.Repositories.GetBranch(ctx, owner, name, branch)
	if err != nil {
		return nil, err
	}

	zerolog.Ctx(ctx).Debug().Any("github_branch", b).Msg("branch loaded from github")

	return b, nil
}

func (me *GithubClient) EnsureRelease(ctx context.Context, repo string, newtag *semver.Version, rel *github.RepositoryRelease) (*github.RepositoryRelease, error) {
	owner, name, err := ParseRepo(repo)
	if err != nil {
		return nil, err
	}

	pr, err := me.GetCurrentPullRequest(ctx)
	if err != nil {
		return nil, err
	}

	prefix := fmt.Sprintf("pr.%d.", pr.GetNumber())

	vers, err := me.ReduceTagVersions(ctx, repo, func(prev *semver.Version, next *semver.Version) *semver.Version {
		if strings.HasPrefix(next.Prerelease(), prefix) && prev.LessThan(next) {
			return next
		}
		return prev
	})
	if err != nil {
		return nil, err
	}

	zerolog.Ctx(ctx).Debug().Str("prefix", prefix).Any("pr", pr).Any("vers", vers).Msg("release version")

	// check if the release already exists
	last, err := me.GetRelease(ctx, repo, vers.String())
	if err != nil {
		return nil, err
	}

	prevId := int64(0)
	if last != nil {
		prevId = last.GetID()
	}

	if prevId == 0 {
		zerolog.Ctx(ctx).Debug().Any("release", rel).Msg("creating release")

		rel, _, err = me.client.Repositories.CreateRelease(ctx, owner, name, rel)
		if err != nil {
			return nil, err
		}
	} else {
		zerolog.Ctx(ctx).Debug().Any("release", rel).Msg("updating release")
		rel, _, err = me.client.Repositories.EditRelease(ctx, owner, name, prevId, rel)
		if err != nil {
			return nil, err
		}
	}

	if rel.Assets != nil {

		zerolog.Ctx(ctx).Debug().Any("assets", rel.Assets).Msg("uploading assets")

		wrkg := sync.WaitGroup{}
		errchan := make(chan error, len(rel.Assets))
		for _, asset := range rel.Assets {
			wrkg.Add(1)
			go func(asset *github.ReleaseAsset) {
				defer wrkg.Done()
				fle, err := os.OpenFile("./buildrc"+asset.GetName(), os.O_RDONLY, 0644)
				if err != nil {
					errchan <- err
					return
				}

				_, _, err = me.client.Repositories.UploadReleaseAsset(ctx, owner, name, rel.GetID(), &github.UploadOptions{
					Name:      asset.GetName(),
					Label:     strings.SplitN(asset.GetName(), "-", 1)[0],
					MediaType: strings.SplitN(asset.GetName(), ".", 1)[1],
				}, fle)
				if err != nil {
					errchan <- err
					return
				}
			}(asset)
		}

		wrkg.Wait()
		close(errchan)

		errs := []error{}

		for err := range errchan {
			errs = append(errs, err)
		}

		if len(errs) > 0 {
			return nil, fmt.Errorf("failed to upload assets: %s", errs)
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

	res, err := me.GetPullRequest(ctx, repository, branch)
	if err != nil {
		return nil, err
	}

	zerolog.Ctx(ctx).Debug().Any("pr", res).Msg("current pull request")

	return res, nil
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

	zerolog.Ctx(ctx).Trace().Str("repo", repo).Msg("reducing tags")

	tags, err := me.ListTags(ctx, repo)
	if err != nil {
		return nil, err
	}

	wrk := semver.New(0, 0, 0, "", "")

	// compare all tags that are not of the format "vX.Y.Z"

	for _, t := range tags {
		if strings.HasPrefix(t.GetName(), "v") {
			ver, err := semver.StrictNewVersion(t.GetName()[1:])
			if err != nil {
				zerolog.Ctx(ctx).Warn().Err(err).Str("tag", t.GetName()).Msg("failed to parse tag")
				continue
			}

			wrk = compare(wrk, ver)
		}

	}

	if wrk.String() == "0.0.0" {

		zerolog.Ctx(ctx).Warn().Any("tags", tags).Any("version", wrk).Msg("no tags found")
		return nil, fmt.Errorf("no tags found")
	}

	zerolog.Ctx(ctx).Trace().Any("tags", tags).Any("version", wrk).Msg("reduced tags")

	return wrk, nil

}

func (me *GithubClient) CountTagVersions(ctx context.Context, repo string, compare Counter[semver.Version]) (int, error) {

	zerolog.Ctx(ctx).Trace().Str("repo", repo).Msg("counting tags")

	tags, err := me.ListTags(ctx, repo)
	if err != nil {
		return 0, err
	}

	wrk := 0

	// compare all tags that are not of the format "vX.Y.Z"

	for _, t := range tags {
		if strings.HasPrefix(t.GetName(), "v") {
			ver, err := semver.StrictNewVersion(t.GetName()[1:])
			if err != nil {
				continue
			}

			if compare(ver) {
				wrk++
			}
		}
	}

	zerolog.Ctx(ctx).Trace().Any("tags", tags).Int("count", wrk).Msg("counted tags")

	return wrk, nil
}

func (me *GithubClient) EnsurePullRequest(ctx context.Context, repo, branch string) (*github.PullRequest, error) {

	zerolog.Ctx(ctx).Trace().Str("repo", repo).Str("branch", branch).Msg("ensuring pull request")

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
	commit, err := me.GetLastCommit(ctx, repo)
	if err == nil {

		mess := commit.GetCommit().GetMessage()

		if len(mess) > 0 {
			re := regexp.MustCompile(`\(issue:(.+)\)`)
			matches := re.FindAllStringSubmatch(mess, -1)
			if len(matches) > 1 {
				iss, err := strconv.Atoi(matches[1][0])
				if err == nil {
					body = fmt.Sprintf("%s\n\nIssue: #%s", body, matches[1][0])
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
		zerolog.Ctx(ctx).Error().Err(err).Any("response", res.Response).Msgf("Failed to create PR: %s", res.Status)
		return nil, err
	}

	zerolog.Ctx(ctx).Trace().Any("pr", pr).Msg("created PR")

	return pr, nil
}

func (me *GithubClient) GetLastCommit(ctx context.Context, repo string) (*github.RepositoryCommit, error) {

	commit, err := GetCurrentCommitSha()
	if err != nil {
		return nil, err
	}

	owner, name, err := ParseRepo(repo)
	if err != nil {
		return nil, err
	}

	resp, _, err := me.client.Repositories.GetCommit(ctx, owner, name, commit)
	if err != nil {
		return nil, err
	}

	return resp, nil
	// get the commit message

}
