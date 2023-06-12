package next

import (
	"context"

	"github.com/nuggxyz/buildrc/internal/provider"
)

type output struct {
	Major           string `json:"major"`
	Minor           string `json:"minor"`
	Patch           string `json:"patch"`
	MajorMinor      string `json:"major_minor"`
	MajorMinorPatch string `json:"major_minor_patch"`
	Full            string `json:"full"`
}

var _ provider.CommandRunner = (*Handler)(nil)
var _ provider.Command[output] = (*Handler)(nil)

func (me *Handler) ID() string {
	return "next-tag"
}

func NewHandler(ctx context.Context, repo string, accessToken string) (*Handler, error) {
	h := &Handler{Repo: repo, AccessToken: accessToken}

	err := h.Init(ctx)

	return h, err

}

func (me *Handler) Helper() provider.CommandHelper[output] {
	return provider.NewHelper[output](me)
}

func (me *Handler) AnyHelper() provider.AnyHelper {
	return provider.NewHelper[output](me)
}
