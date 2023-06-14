package next

import (
	"context"
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/nuggxyz/buildrc/cmd/buildrc/load"
	"github.com/nuggxyz/buildrc/internal/docker"
	"github.com/nuggxyz/buildrc/internal/env"
	"github.com/nuggxyz/buildrc/internal/github"
	"github.com/nuggxyz/buildrc/internal/provider"
	"github.com/rs/zerolog"
)

const (
	CommandID = "tag_next"
)

type Output struct {
	Major           string `json:"major"`
	Minor           string `json:"minor"`
	Patch           string `json:"patch"`
	MajorMinor      string `json:"major_minor"`
	MajorMinorPatch string `json:"major_minor_patch"`
	Full            string `json:"full" express:"BUILDRC_TAG_NEXT_FULL"`
	BuildxTags      string `json:"buildx_tags" express:"BUILDRC_TAG_NEXT_BUILDX_TAGS"`
}

type Handler struct {
	Repo        string `flag:"repo" type:"repo:" default:""`
	File        string `flag:"file" type:"file:" default:".buildrc"`
	AccessToken string `flag:"token" type:"access_token:" default:""`
}

func NewHandler(repo string, accessToken string) *Handler {
	h := &Handler{Repo: repo, AccessToken: accessToken}
	return h
}

func (me *Handler) Run(ctx context.Context, cp provider.ContentProvider) (err error) {
	_, err = me.Next(ctx, cp)
	return err
}

func (me *Handler) Next(ctx context.Context, cp provider.ContentProvider) (out *Output, err error) {
	return provider.Wrap(CommandID, me.next)(ctx, cp)
}

func (me *Handler) next(ctx context.Context, prv provider.ContentProvider) (out *Output, err error) {

	if me.AccessToken == "" {
		zerolog.Ctx(ctx).Debug().Msg("No access token provided, trying to get from env")
		// TODO: this should be a helper function, could grab from somewhere else
		me.AccessToken = env.GetOrEmpty("GITHUB_TOKEN")
		if me.AccessToken == "" {
			zerolog.Ctx(ctx).Debug().Msg("❌ No access token found in env")
		} else {
			zerolog.Ctx(ctx).Debug().Msg("✅ Access token found in env")
		}
	}

	if me.Repo == "" {

		zerolog.Ctx(ctx).Debug().Msg("No repo provided, trying to get from env")

		curr, err := github.GetCurrentRepo()
		if err != nil {
			return nil, err
		}

		zerolog.Ctx(ctx).Debug().Msgf("✅ Repo found in env: %s", curr)

		me.Repo = curr
	}

	brc, err := load.NewHandler(me.File).Load(ctx, prv)
	if err != nil {
		return nil, err
	}

	vers, err := calculateNextVersion(ctx, me.AccessToken, me.Repo)
	if err != nil {
		return nil, err
	}

	if brc.Version.GreaterThan(vers) {
		vers = semver.New(brc.Version.Major(), 0, 0, vers.Prerelease(), vers.Metadata())
	}

	str, err := docker.BuildXTagString(ctx, me.Repo, vers.String())
	if err != nil {
		return nil, err
	}

	return &Output{
		Major:           fmt.Sprintf("%d", vers.Major()),
		Minor:           fmt.Sprintf("%d", vers.Minor()),
		Patch:           fmt.Sprintf("%d", vers.Patch()),
		MajorMinor:      fmt.Sprintf("%d.%d", vers.Major(), vers.Minor()),
		MajorMinorPatch: fmt.Sprintf("%d.%d.%d", vers.Major(), vers.Minor(), vers.Patch()),
		Full:            vers.String(),
		BuildxTags:      str,
	}, nil
}

func calculateNextVersion(ctx context.Context, token, repo string) (out *semver.Version, err error) {
	// get the current main highest tag
	ghc, err := github.NewGithubClient(ctx, token)
	if err != nil {
		return nil, err
	}

	brnch, err := github.GetCurrentBranch()
	if err != nil {
		return nil, err
	}

	isMerge := brnch == "main"

	if isMerge {
		commit, err := github.GetCurrentCommit()
		if err != nil {
			return nil, err
		}

		last, err := ghc.ReduceTagVersions(ctx, repo, func(prev, next *semver.Version) *semver.Version {
			if prev.GreaterThan(next) && prev.Prerelease() == "" {
				return prev
			}
			return next
		})
		if err != nil {
			return nil, err
		}

		if strings.Contains(commit, "(feat)") {
			res := last.IncMinor()
			return &res, nil
		} else {
			res := last.IncPatch()
			return &res, nil
		}
	}

	pr, err := ghc.EnsurePullRequest(ctx, repo, brnch)
	if err != nil {
		return nil, err
	}

	if pr == nil {
		return nil, fmt.Errorf("no pull request found for branch %s", brnch)
	}

	prefix := fmt.Sprintf("pr.%d.", pr.Number)

	main, err := ghc.ReduceTagVersions(ctx, repo, func(prev, next *semver.Version) *semver.Version {
		shouldConsider := next.Prerelease() == "" || strings.HasPrefix(next.Prerelease(), prefix)

		if !shouldConsider {
			return prev
		}

		if prev.GreaterThan(next) {
			return prev
		}
		return next
	})
	if err != nil {
		return nil, err
	}

	isPrerelease := pr.GetState() == "open"
	isFeature := strings.HasPrefix(pr.GetTitle(), "feat")

	buildnum := "0"

	if isPrerelease {

		prc, err := ghc.CountTagVersions(ctx, repo, func(v *semver.Version) bool {
			return v.Prerelease() != "" && strings.HasPrefix(v.Prerelease(), prefix)
		})

		if err != nil {
			return nil, err
		}

		buildnum = fmt.Sprintf("%d", prc)
	}

	shouldInc := !strings.HasPrefix(main.Prerelease(), prefix)

	if shouldInc {
		if isFeature {
			main.IncMinor()
		} else {
			main.IncPatch()
		}
	}

	main.SetPrerelease(fmt.Sprintf("%s%s", prefix, buildnum))

	return main, nil

}
