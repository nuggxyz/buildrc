package setup

import (
	"context"

	"github.com/Masterminds/semver/v3"
	"github.com/nuggxyz/buildrc/internal/common"
	"github.com/nuggxyz/buildrc/internal/git"
	"github.com/nuggxyz/buildrc/internal/pipeline"
)

const (
	CommandID = "setup"
)

type Handler struct {
	Repo        string `flag:"repo" type:"repo:" default:""`
	File        string `flag:"file" type:"file:" default:".buildrc"`
	AccessToken string `flag:"token" type:"access_token:" default:""`
}

func NewHandler(repo string, accessToken string) *Handler {
	h := &Handler{Repo: repo, AccessToken: accessToken}
	return h
}

func (me *Handler) Run(ctx context.Context, cp common.Provider) (err error) {
	_, err = me.Invoke(ctx, cp)
	return err
}

type Response struct {
	TagSemver *semver.Version
	Tag       string
	Full      string `json:"full" express:"BUILDRC_RELEASE_FINALIZE_FULL"`
	ReleaseID string `json:"release_id" express:"BUILDRC_RELEASE_FINALIZE_RELEASE_ID"`
}

func (me *Handler) Invoke(ctx context.Context, prov common.Provider) (out *Response, err error) {
	return pipeline.Cache(ctx, CommandID, prov, me.invoke)
}

func (me *Handler) invoke(ctx context.Context, prov common.Provider) (out *Response, err error) {

	targetSemver, err := git.CalculateNextPreReleaseTag(ctx, prov.Buildrc(), prov.Git(), prov.PR())
	if err != nil {
		return nil, err
	}

	next, err := prov.Release().TagRelease(ctx, prov.Git(), targetSemver)
	if err != nil {
		return nil, err
	}

	err = pipeline.AddContentToEnv(ctx, prov.Pipeline(), prov.FileSystem(), CommandID, map[string]string{
		"tag": targetSemver.String(),
	})

	return &Response{
		Tag: targetSemver.String(), TagSemver: targetSemver, Full: targetSemver.String(),
		ReleaseID: next.ID}, err
}
