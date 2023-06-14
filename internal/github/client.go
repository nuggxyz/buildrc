package github

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/Masterminds/semver/v3"
	"github.com/google/go-github/v53/github"
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

	b, _, err := me.client.Repositories.GetBranch(ctx, owner, name, branch, true)
	if err != nil {
		return nil, err
	}

	zerolog.Ctx(ctx).Debug().Any("github_branch", b).Msg("branch loaded from github")

	return b, nil
}

func (me *GithubClient) EnsureRelease(ctx context.Context, repo string, upd *semver.Version, assets []string) (*github.RepositoryRelease, error) {
	owner, name, err := ParseRepo(repo)
	if err != nil {
		return nil, err
	}

	pr, err := me.GetCurrentPullRequest(ctx)
	if err != nil {
		return nil, err
	}

	prefixer := "beta"

	if pr.GetDraft() {
		prefixer = "alpha"
	}

	prefix := fmt.Sprintf("%s.%d.", prefixer, pr.GetNumber())

	isParent := func(v *semver.Version) bool {
		return strings.HasPrefix(v.Prerelease(), prefix)
	}

	cnt, err := me.CountTagVersions(ctx, repo, isParent)
	if err != nil {
		return nil, err
	}

	prev, err := me.ReduceTagVersions(ctx, repo, func(prev *semver.Version, next *semver.Version) *semver.Version {
		if isParent(next) && (prev == nil || next.GreaterThan(prev)) {
			return next
		}
		return prev
	})
	if err != nil {
		return nil, err
	}

	prefix += strconv.Itoa(cnt + 1)

	prevId := int64(0)

	zerolog.Ctx(ctx).Debug().Str("prefix", prefix).Any("prev", prev).Msg("release previon")

	if prev != nil {
		tag := "v" + prev.String()
		// check if the release already exists
		last, err := me.GetRelease(ctx, repo, tag)
		if err != nil {
			return nil, err
		}

		if last != nil {
			prevId = last.GetID()
		}

		zerolog.Ctx(ctx).Debug().Any("last", last).Msg("last release")
	}

	cmt, err := me.GetLastCommit(ctx, repo)
	if err != nil {
		return nil, err
	}

	vn, err := upd.SetPrerelease(prefix)
	if err != nil {
		return nil, err
	}

	rel := &github.RepositoryRelease{
		TagName:         github.String("v" + vn.String()),
		TargetCommitish: cmt.SHA,
		Name:            github.String(fmt.Sprintf("PR #%d", pr.GetNumber())),
		Author:          cmt.Author,
		Prerelease:      github.Bool(true),
		Draft:           github.Bool(false),
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

	if assets != nil {

		zerolog.Ctx(ctx).Debug().Any("assets", assets).Msg("uploading assets")

		wrkg := sync.WaitGroup{}
		errchan := make(chan error)
		for _, asset := range assets {
			wrkg.Add(1)
			go func(asset string) {
				defer wrkg.Done()
				fle, err := os.OpenFile(asset, os.O_RDONLY, 0644)
				if err != nil {
					errchan <- err
					return
				}

				defer fle.Close()

				for _, a := range rel.Assets {
					if a.GetName() == filepath.Base(asset) {
						_, err = me.client.Repositories.DeleteReleaseAsset(ctx, owner, name, a.GetID())
						if err != nil {
							errchan <- err
							return
						}
					}
				}

				_, _, err = me.client.Repositories.UploadReleaseAsset(ctx, owner, name, rel.GetID(), &github.UploadOptions{
					Name:  filepath.Base(asset),
					Label: strings.SplitN(filepath.Base(asset), "-", 1)[0],
				}, fle)
				if err != nil {
					errchan <- err
					return
				}
			}(asset)
		}

		ctx, cancel := context.WithCancel(ctx)
		go func() {
			defer cancel()
			wrkg.Wait()
		}()

		for {
			select {
			case <-ctx.Done():
				if ctx.Err() != context.Canceled {
					return nil, ctx.Err()
				}
				return rel, nil
			case err := <-errchan:
				return nil, err
			}
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

	return res, nil
}

func (me *GithubClient) GetRelease(ctx context.Context, repo, tag string) (*github.RepositoryRelease, error) {
	owner, name, err := ParseRepo(repo)
	if err != nil {
		return nil, err
	}

	r, _, err := me.client.Repositories.GetReleaseByTag(ctx, owner, name, tag)
	if err != nil {
		if strings.Contains(err.Error(), tag+": 404 Not Found") {
			return nil, nil
		}

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

	if len(pulls) == 0 {
		return nil, nil
	}

	prselection := pulls[0]

	zerolog.Ctx(ctx).Debug().Int("total_found", len(pulls)).Any("args", opts).Int("selected", prselection.GetNumber()).Msg("pull requests loaded from github")

	// If there's more than one matching PR (which is unusual), return the first one.
	return prselection, nil
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
		zerolog.Ctx(ctx).Trace().Str("tag", t.GetName()).Msg("checking tag")
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
		return nil, nil
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

	resp, _, err := me.client.Repositories.GetCommit(ctx, owner, name, commit, &github.ListOptions{})
	if err != nil {
		return nil, err
	}

	return resp, nil
	// get the commit message

}
