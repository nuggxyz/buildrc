package setup

import (
	"context"

	"github.com/nuggxyz/buildrc/cmd/buildrc/load"
	"github.com/nuggxyz/buildrc/internal/common"
	"github.com/nuggxyz/buildrc/internal/github"
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

func (me *Handler) Invoke(ctx context.Context, cp common.Provider) (out *Response, err error) {
	return pipeline.Cache(ctx, CommandID, cp, me.invoke)
}

func (me *Handler) invoke(ctx context.Context, r common.Provider) (out *Response, err error) {

	brc, err := load.NewHandler(me.File).Load(ctx, r)
	if err != nil {
		return nil, err
	}

	ghc, err := github.NewGithubClient(ctx, "", "")
	if err != nil {
		return nil, err
	}

	t, rid, err := ghc.Setup(ctx, brc.Version)

	if err != nil {
		return nil, err
	}

	err = pipeline.AddContentToEnv(ctx, r.Pipeline(), CommandID, map[string]string{
		"tag":                t,
		"unique_release_tag": rid,
	})

	return &Response{Tag: t, UniqueReleaseTag: rid}, err
}
