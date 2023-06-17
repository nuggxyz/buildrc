package github

import (
	"context"
	"fmt"
	"net/http"
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
	client   *github.Client
	repo     string
	repoName string
	orgName  string
}

func (me *GithubClient) RepoName() string {
	return me.repoName
}

func (me *GithubClient) OrgName() string {
	return me.orgName
}

func (me *GithubClient) Client() *github.Client {
	return me.client
}

func NewGithubClient(ctx context.Context, token string, repo string) (*GithubClient, error) {

	var err error

	if token == "" {

		zerolog.Ctx(ctx).Debug().Msg("no token specified, trying to get from env")

		token, err = GetGithubTokenFromEnv(ctx)
		if err != nil {
			zerolog.Ctx(ctx).Debug().Msg("no token found in env, set one to the GITHUB_TOKEN env var")
			return nil, err
		}

		zerolog.Ctx(ctx).Debug().Msg("✅ Token found in env")
	}

	if repo == "" {

		zerolog.Ctx(ctx).Debug().Msg("no repo specified, trying to get from git config")

		repo, err = GetCurrentRepo()
		if err != nil {
			zerolog.Ctx(ctx).Debug().Msg("no repo found in git config, is this a git repo?")
			return nil, err
		}

		zerolog.Ctx(ctx).Debug().Msgf("✅ Repo found in env: %s", repo)

	}

	owner, name, err := ParseRepo(repo)
	if err != nil {
		return nil, err
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// check that the token is valid
	_, _, err = client.Zen(ctx)
	if err != nil {
		return nil, err
	}

	return &GithubClient{client: client, repo: repo, repoName: name, orgName: owner}, nil
}

func (me *GithubClient) ListTags(ctx context.Context) ([]*github.RepositoryTag, error) {

	tags := []*github.RepositoryTag{}

	opts := &github.ListOptions{PerPage: 100}
	for {
		t, resp, err := me.client.Repositories.ListTags(ctx, me.OrgName(), me.RepoName(), opts)
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

func (me *GithubClient) GetBranch(ctx context.Context, branch string) (*github.Branch, error) {

	b, _, err := me.client.Repositories.GetBranch(ctx, me.OrgName(), me.RepoName(), branch, true)
	if err != nil {
		return nil, err
	}

	zerolog.Ctx(ctx).Debug().Any("github_branch", b).Msg("branch loaded from github")

	return b, nil
}

func (me *GithubClient) EnsureRelease(ctx context.Context, majorRef *semver.Version) (*github.RepositoryRelease, error) {

	branch, err := GetCurrentBranch()
	if err != nil {
		return nil, err
	}

	var pr *github.PullRequest

	if branch != "main" {
		pr, err = me.EnsurePullRequest(ctx, branch)
		if err != nil {
			return nil, err
		}
	}

	cmt, err := me.GetLastCommit(ctx)
	if err != nil {
		return nil, err
	}

	vn, _, err := me.CalculateNextPreReleaseTag(ctx, majorRef, pr)
	if err != nil {
		return nil, err
	}

	isPrerelease := vn.Prerelease() != ""
	releaseName := ""

	if isPrerelease {
		releaseName = fmt.Sprintf("PR #%d", pr.GetNumber())
	} else {
		releaseName = fmt.Sprintf("v%s", vn.String())
	}

	rel := &github.RepositoryRelease{
		TagName:         github.String("v" + vn.String()),
		TargetCommitish: cmt.SHA,
		Name:            github.String(releaseName),
		Author:          cmt.Author,
		Prerelease:      github.Bool(vn.Prerelease() != ""),
		Draft:           github.Bool(true),
	}

	rel, _, err = me.client.Repositories.CreateRelease(ctx, me.OrgName(), me.RepoName(), rel)
	if err != nil {
		return nil, err
	}

	zerolog.Ctx(ctx).Debug().Any("github_release", rel).Msg("release created")

	// if prevId == 0 {
	// 	zerolog.Ctx(ctx).Debug().Any("release", rel).Msg("creating release")
	// 	rel, _, err = me.client.Repositories.CreateRelease(ctx, me.OrgName(), me.RepoName(), rel)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// } else {
	// 	zerolog.Ctx(ctx).Debug().Any("release", rel).Msg("updating release")
	// 	rel, _, err = me.client.Repositories.EditRelease(ctx, me.OrgName(), me.RepoName(), prevId, rel)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	// // only update assets if we are not on main or if we are on main and this is not a PR merge
	// shouldUpdateAssets := !(branch == "main" && pr != nil)

	// if shouldUpdateAssets {
	// 	rel, err = me.UpdateReleaseAssets(ctx, rel)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	return rel, nil
}

func (me *GithubClient) FinalizeRelease(ctx context.Context) (*github.RepositoryRelease, error) {

	t, err := GetCurrentTag()
	if err != nil {
		return nil, err
	}

	r, _, err := me.Client().Repositories.GetReleaseByTag(ctx, me.OrgName(), me.RepoName(), t)
	if err != nil {
		return nil, err
	}

	r.Draft = github.Bool(false)

	rel, _, err := me.client.Repositories.EditRelease(ctx, me.OrgName(), me.RepoName(), r.GetID(), r)
	if err != nil {
		return nil, err
	}

	return rel, nil
}

func (me *GithubClient) EnsureDraftRelease(ctx context.Context) (*github.RepositoryRelease, error) {

	r, _, err := me.client.Repositories.CreateRelease(ctx, me.OrgName(), me.RepoName(), &github.RepositoryRelease{})
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (me *GithubClient) TagCommit(ctx context.Context, tag string) error {

	sha, err := GetCurrentCommitSha()
	if err != nil {
		return err
	}

	_, _, err = me.client.Git.CreateRef(ctx, me.OrgName(), me.RepoName(), &github.Reference{
		Ref: github.String("refs/tags/" + tag),
		Object: &github.GitObject{
			SHA: github.String(sha),
		},
	})

	return err
}

func (me *GithubClient) UpdateReleaseAssets(ctx context.Context, rel *github.RepositoryRelease) (*github.RepositoryRelease, error) {

	zerolog.Ctx(ctx).Debug().Msg("uploading assets")

	files, err := FindFiles("./buildrc", ".tar.gz", ".sha256")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	zerolog.Ctx(ctx).Debug().Strs("files", files).Msg("found files")

	wrkg := sync.WaitGroup{}
	errchan := make(chan error)
	for _, asset := range files {
		wrkg.Add(1)
		go func(asset string) {
			defer wrkg.Done()

			for _, a := range rel.Assets {
				if a.GetName() == filepath.Base(asset) {
					_, err = me.client.Repositories.DeleteReleaseAsset(ctx, me.OrgName(), me.RepoName(), a.GetID())
					if err != nil {
						errchan <- err
						return
					}
				}
			}

			fle, err := os.Open(asset)
			if err != nil {
				errchan <- err
				return
			}

			_, _, err = me.client.Repositories.UploadReleaseAsset(ctx, me.OrgName(), me.RepoName(), rel.GetID(), &github.UploadOptions{
				Name: filepath.Base(asset),
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

func (me *GithubClient) CalculateNextPreReleaseTag(ctx context.Context, majorRef *semver.Version, pr *github.PullRequest) (*semver.Version, int64, error) {

	prevMain, err := me.ReduceTagVersions(ctx, func(prev *semver.Version, next *semver.Version) *semver.Version {
		if next.Prerelease() == "" && (prev == nil || next.GreaterThan(prev)) {
			return next
		}
		return prev
	})

	if err != nil {
		return nil, 0, err
	}

	if pr == nil {
		// if there is no pr, then this was a direct commit to main
		// so we just increment the patch version

		brnch, err := GetCurrentBranch()
		if err != nil {
			return nil, 0, err
		}

		if brnch != "main" {
			return nil, 0, fmt.Errorf("current branch is not main and no pull request was provided")
		}

		if prevMain == nil {
			prevMain = majorRef
		}

		res := prevMain.IncPatch()

		return &res, 0, nil

	}

	isFeature := strings.HasPrefix(pr.GetTitle(), "feat")

	prefixer := "beta"

	if pr.GetDraft() {
		prefixer = "alpha"
	}

	prefix := fmt.Sprintf("%s.%d.", prefixer, pr.GetNumber())

	isParent := func(v *semver.Version) bool {
		return strings.HasPrefix(v.Prerelease(), prefix)
	}

	cnt, err := me.CountTagVersions(ctx, isParent)
	if err != nil {
		return nil, 0, err
	}

	prev, err := me.ReduceTagVersions(ctx, func(prev *semver.Version, next *semver.Version) *semver.Version {
		if isParent(next) && (prev == nil || next.GreaterThan(prev)) {
			return next
		}
		return prev
	})
	if err != nil {
		return nil, 0, err
	}

	prefix += strconv.Itoa(cnt + 1)

	prevId := int64(0)

	zerolog.Ctx(ctx).Debug().Str("prefix", prefix).Any("prev", prev).Msg("release previon")

	if prev != nil {
		tag := "v" + prev.String()
		// check if the release already exists
		last, err := me.GetRelease(ctx, tag)
		if err != nil {
			return nil, 0, err
		}

		if last != nil {
			prevId = last.GetID()
		}

		zerolog.Ctx(ctx).Debug().Any("last", last).Msg("last release")
	} else {
		zerolog.Ctx(ctx).Debug().Msg("no previous release, looking for a tag on main")
		// check if there is a tag on main

		if err != nil {
			return nil, 0, err
		}
		if prev == nil {
			prev = majorRef
		}
	}

	if prevMain != nil && prevMain.GreaterThan(prev) {
		prev = prevMain
	}

	if majorRef.GreaterThan(prev) {
		prev = majorRef
	}

	shouldInc := !strings.HasPrefix(prev.Prerelease(), prefix)

	wrk := *prev

	if shouldInc {
		if isFeature {
			wrk = wrk.IncMinor()
		} else {
			wrk = wrk.IncPatch()
		}
	}

	vn, err := wrk.SetPrerelease(prefix)
	if err != nil {
		return nil, 0, err
	}

	zerolog.Ctx(ctx).Debug().Str("prefix", prefix).Any("prev", prev).Any("vn", vn).Msg("release version")

	return &vn, prevId, nil
}

func (me *GithubClient) GetRelease(ctx context.Context, tag string) (*github.RepositoryRelease, error) {

	r, _, err := me.client.Repositories.GetReleaseByTag(ctx, me.OrgName(), me.RepoName(), tag)
	if err != nil {
		if strings.Contains(err.Error(), tag+": 404 Not Found") {
			return nil, nil
		}

		return nil, err
	}

	return r, nil
}

func (me *GithubClient) GetOpenPullRequestForBranch(ctx context.Context, branch string) (*github.PullRequest, error) {

	opts := &github.PullRequestListOptions{
		State:       "open",
		Head:        branch,
		Base:        "main",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	pulls, res, err := me.client.PullRequests.List(ctx, me.OrgName(), me.RepoName(), opts)
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

func (me *GithubClient) ReduceTagVersions(ctx context.Context, compare Reducer[semver.Version]) (*semver.Version, error) {

	zerolog.Ctx(ctx).Trace().Msg("reducing tags")

	tags, err := me.ListTags(ctx)
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

func (me *GithubClient) CountTagVersions(ctx context.Context, compare Counter[semver.Version]) (int, error) {

	zerolog.Ctx(ctx).Trace().Msg("counting tags")

	tags, err := me.ListTags(ctx)
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

func (me *GithubClient) AddIssueToPullRequestBody(ctx context.Context, issue int, pr *github.PullRequest) error {

	_, _, err := me.client.PullRequests.Edit(ctx, me.OrgName(), me.RepoName(), pr.GetNumber(), &github.PullRequest{
		Body: github.String(fmt.Sprintf("%s\nresolves #%d", pr.GetBody(), issue)),
	})
	if err != nil {
		return err
	}

	return nil
}

func (me *GithubClient) GetReferencedIssueByLastCommit(ctx context.Context) ([]int, error) {
	issue := []int{}

	commit, err := me.GetLastCommit(ctx)
	if err != nil {
		return nil, err
	}

	mess := commit.GetCommit().GetMessage()

	if len(mess) > 0 {
		re := regexp.MustCompile(`#(.+?)`) // Change here
		matches := re.FindAllStringSubmatch(mess, -1)
		if len(matches) > 0 {
			for _, match := range matches {
				iss, err := strconv.Atoi(match[1]) // No change here
				if err == nil {
					issue = append(issue, iss)
				}
			}
		}
	}

	return issue, nil
}

func (me *GithubClient) EnsurePullRequest(ctx context.Context, branch string) (*github.PullRequest, error) {

	zerolog.Ctx(ctx).Trace().Str("repo", me.RepoName()).Str("branch", branch).Msg("ensuring pull request")

	req := &github.NewPullRequest{
		Title:               github.String(fmt.Sprintf("Release %s", branch)),
		Body:                github.String("Automatically generated release PR. Please update."),
		Base:                github.String("main"),
		Head:                github.String(branch),
		Draft:               github.Bool(true),
		MaintainerCanModify: github.Bool(true),
	}

	issue, err := me.GetReferencedIssueByLastCommit(ctx)
	if err != nil {
		return nil, err
	}

	if len(issue) > 0 {
		req.Issue = github.Int(issue[0])
	}

	// if commit message has "(issue:xxx)", add that to the PR body
	next, abc := me.GetOpenPullRequestForBranch(ctx, branch)
	if abc != nil {
		return nil, abc
	}

	for _, issue := range issue {
		if next == nil {
			req.Body = github.String(fmt.Sprintf("%s\nresolves #%d", req.GetBody(), issue))
		} else {
			err := me.AddIssueToPullRequestBody(ctx, issue, next)
			if err != nil {
				zerolog.Ctx(ctx).Error().Err(err).Int("issue", issue).Int("pr", next.GetNumber()).Msg("failed to add issue to PR")
				return nil, err
			}
			zerolog.Ctx(ctx).Debug().Int("issue", issue).Int("pr", next.GetNumber()).Msg("added issue to PR")
		}
	}

	if next != nil {
		zerolog.Ctx(ctx).Debug().Any("pr", next).Msg("PR already exists")

		return next, nil
	}

	zerolog.Ctx(ctx).Debug().Msgf("Creating PR for %s", branch)

	// create a new PR
	pr, res, err := me.client.PullRequests.Create(ctx, me.OrgName(), me.RepoName(), req)

	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msgf("Failed to create PR: %d", res.StatusCode)
		return nil, err
	}

	zerolog.Ctx(ctx).Trace().Int("pr", pr.GetNumber()).Msg("created PR")

	return pr, nil
}

func (me *GithubClient) GetLastCommit(ctx context.Context) (*github.RepositoryCommit, error) {

	commit, err := GetCurrentCommitSha()
	if err != nil {
		return nil, err
	}

	resp, _, err := me.client.Repositories.GetCommit(ctx, me.OrgName(), me.RepoName(), commit, &github.ListOptions{})
	if err != nil {
		return nil, err
	}

	return resp, nil
	// get the commit message
}

func (me *GithubClient) GetClosedPullRequestFromCommit(ctx context.Context, commit *github.RepositoryCommit) (*github.PullRequest, error) {

	// List all pull requests
	opts := &github.PullRequestListOptions{
		State: "closed",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
		Sort:      "updated",
		Direction: "desc",
		Base:      "main",
	}
	pulls, r, err := me.client.PullRequests.List(ctx, me.OrgName(), me.RepoName(), opts)
	if err != nil {
		if r != nil && r.StatusCode == http.StatusNotFound {
			return nil, nil
		}
		return nil, err
	}

	// Iterate over the PRs
	for _, pr := range pulls {
		res, err := me.isSameCode(ctx, commit.GetCommit(), pr.GetHead().GetSHA())
		if err != nil {
			return nil, err
		}
		if res {
			return pr, nil
		}
	}

	// Return nil if no matching PR was found
	return nil, nil
}

func (me *GithubClient) isSameCode(ctx context.Context, commit1 *github.Commit, commitSHA2 string) (bool, error) {

	commit2, _, err := me.client.Repositories.GetCommit(ctx, me.OrgName(), me.RepoName(), commitSHA2, &github.ListOptions{})
	if err != nil {
		return false, err
	}

	return commit1.GetTree().GetSHA() == commit2.GetCommit().GetTree().GetSHA(), nil
}
