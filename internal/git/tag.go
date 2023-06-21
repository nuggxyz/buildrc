package git

import (
	"context"
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/nuggxyz/buildrc/internal/buildrc"
	"github.com/rs/zerolog"
)

func CalculateNextPreReleaseTag(ctx context.Context, brc *buildrc.Buildrc, git GitProvider, prp PullRequestProvider) (*semver.Version, error) {

	latestHead, err := git.GetLatestSemverTagFromRef(ctx, "HEAD")
	if err != nil {
		return nil, err
	}

	latestMain, err := git.GetLatestSemverTagFromRef(ctx, "main")
	if err != nil {
		return nil, err
	}

	latestMajor := semver.New(uint64(brc.Version), 0, 0, "", "")

	pr, err := getLatestPullRequest(ctx, prp, "HEAD")
	if err != nil {
		return nil, err
	}

	if pr == nil {
		// if there is no pr, then this was a direct commit to main
		// so we just increment the patch version

		brnch, err := git.GetCurrentBranch(ctx)
		if err != nil {
			return nil, err
		}

		if brnch != "main" {
			return nil, fmt.Errorf("current branch is not main and no pull request was provided")
		}

		if latestMain == nil {
			latestMain = latestMajor
		}

		res := latestMain.IncPatch()

		return &res, nil
	}

	// isFeature := strings.HasPrefix(pr.GetTitle(), "feat")

	// prefixer := "beta"

	prefix := pr.PreReleaseTag()

	if latestMain != nil && latestMain.GreaterThan(latestHead) {
		latestHead = latestMain
	}

	if latestMajor.GreaterThan(latestHead) {
		latestHead = latestMajor
	}

	shouldInc := !strings.Contains(latestHead.Prerelease(), prefix)

	var result semver.Version

	if shouldInc {
		result = latestHead.IncMinor()
	} else {
		result = *latestHead
	}

	result, err = result.SetPrerelease(prefix)
	if err != nil {
		return nil, err
	}

	zerolog.Ctx(ctx).Debug().Str("prefix", prefix).Any("latestHead", latestHead).Any("result", result).Msg("release version")

	return &result, nil
}
