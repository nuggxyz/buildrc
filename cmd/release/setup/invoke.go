package setup

import (
	"context"

	"github.com/Masterminds/semver/v3"
	"github.com/nuggxyz/buildrc/internal/common"
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
	Tag              string
	UniqueReleaseTag string
}

func (me *Handler) Invoke(ctx context.Context, prov common.Provider) (out *Response, err error) {
	return pipeline.Cache(ctx, CommandID, prov, me.invoke)
}

func (me *Handler) invoke(ctx context.Context, prov common.Provider) (out *Response, err error) {

	crt, err := prov.Release().CreateRelease(ctx, prov.Git(), &semver.Version{})
	if err != nil {
		return nil, err
	}

	smvr, err := crt.Semver()
	if err != nil {
		return nil, err
	}

	err = pipeline.AddContentToEnv(ctx, prov.Pipeline(), CommandID, map[string]string{
		"tag":                smvr.String(),
		"unique_release_tag": crt.UntaggedTag,
	})

	return &Response{Tag: smvr.String(), UniqueReleaseTag: crt.UntaggedTag}, err
}
