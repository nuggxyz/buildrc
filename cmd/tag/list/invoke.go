package list

import (
	"github.com/nuggxyz/buildrc/internal/github"
	"github.com/nuggxyz/buildrc/internal/provider"

	"context"
	"errors"

	"github.com/Masterminds/semver/v3"
)

const (
	CommandID = "tag_list"
)

type Output struct {
	VersionTags    []semver.Version `json:"tags" yaml:"tags"`
	HighestVersion semver.Version   `json:"highest_version" yaml:"highest_version"`
	NonVersionTags []string         `json:"all_tags" yaml:"all_tags"`
}

type Handler struct {
	Repo        string `flag:"repo" type:"repo:" default:""`
	File        string `flag:"file" type:"file:" default:".buildrc"`
	AccessToken string `flag:"token" type:"access_token:" default:""`
}

func (me *Handler) Help() string {
	return ``
}

func NewHandler(repo string, accessToken string) *Handler {
	h := &Handler{Repo: repo, AccessToken: accessToken}

	return h
}

func (me *Handler) Run(ctx context.Context, cp provider.ContentProvider) (err error) {
	_, err = me.Invoke(ctx, cp)
	return err
}

func (me *Handler) Invoke(ctx context.Context, cp provider.ContentProvider) (out *Output, err error) {
	return provider.Wrap(CommandID, me.invoke)(ctx, cp)
}

func (me *Handler) invoke(ctx context.Context, _ provider.ContentProvider) (out *Output, err error) {
	ghc, err := github.NewGithubClient(ctx, me.AccessToken, me.Repo)
	if err != nil {
		return nil, err
	}

	// Get tags
	tags, err := ghc.ListTags(ctx)
	if err != nil {
		return nil, err
	}

	if len(tags) == 0 {
		return nil, errors.New("no tags found in the repository")
	}

	out = &Output{
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
