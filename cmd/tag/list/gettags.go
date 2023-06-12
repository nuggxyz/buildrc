package list

import (
	"github.com/nuggxyz/buildrc/internal/cli"
	"github.com/nuggxyz/buildrc/internal/github"

	"context"
	"errors"

	"github.com/Masterminds/semver/v3"
)

type Output struct {
	VersionTags    []semver.Version `json:"tags" yaml:"tags"`
	HighestVersion semver.Version   `json:"highest_version" yaml:"highest_version"`
	NonVersionTags []string         `json:"all_tags" yaml:"all_tags"`
}

type output = Output

type Handler struct {
	Repo        string `arg:"repo" type:"repo:" required:"true"`
	AccessToken string `arg:"GITHUB_TOKEN" type:"access_token:" required:"true"`

	githubClient github.GithubAPI
}

var _ cli.CommandRunner = (*Handler)(nil)
var _ cli.Command[output] = (*Handler)(nil)

func (me *Handler) Help() string {
	return ``
}

func (me *Handler) ID() string {
	return "next-tag"
}

func (me *Handler) Init(ctx context.Context) error {
	ghc, err := github.NewGithubClient(ctx, me.AccessToken)
	if err != nil {
		return err
	}

	me.githubClient = ghc

	return nil
}

func NewHandler(ctx context.Context, repo string, accessToken string) (*Handler, error) {
	h := &Handler{Repo: repo, AccessToken: accessToken}
	err := h.Init(ctx)
	if err != nil {
		return nil, err
	}
	return h, nil
}

func (me *Handler) Helper() cli.CommandHelper[output] {
	return cli.NewHelper[output](me)
}

func (me *Handler) AnyHelper() cli.AnyHelper {
	return cli.NewHelper[output](me)
}

func (me *Handler) Invoke(ctx context.Context, _ cli.ContentProvider) (out *output, err error) {

	// Get tags
	tags, err := me.githubClient.ListTags(ctx, me.Repo)
	if err != nil {
		return nil, err
	}

	if len(tags) == 0 {
		return nil, errors.New("no tags found in the repository")
	}

	out = &output{
		VersionTags:    make([]semver.Version, 0),
		NonVersionTags: make([]string, 0),
	}

	// Filter tags based on semver and find the highest version
	var highest *semver.Version
	for _, tag := range tags {
		v, err := semver.NewVersion(tag.GetName())
		if err == nil {
			out.VersionTags = append(out.VersionTags, *v)
			if highest == nil || v.GreaterThan(highest) {
				highest = v
			}
		} else {
			out.NonVersionTags = append(out.NonVersionTags, tag.GetName())
		}

	}

	if highest == nil {
		return nil, errors.New("no valid semver tags found in the repository")
	}

	out.HighestVersion = *highest

	return out, nil

}
