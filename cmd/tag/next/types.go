package next

import (
	"context"

	"github.com/nuggxyz/buildrc/internal/provider"
)

type TagNextOutput struct {
	Major           string `json:"major"`
	Minor           string `json:"minor"`
	Patch           string `json:"patch"`
	MajorMinor      string `json:"major_minor"`
	MajorMinorPatch string `json:"major_minor_patch"`
	Full            string `json:"full"`
	BuildxTags      string `json:"buildx_tags"`
}

var _ provider.CommandRunner = (*Handler)(nil)
var _ provider.Command[TagNextOutput] = (*Handler)(nil)

func (me *Handler) ID() string {
	return "tag_next"
}

func NewHandler(ctx context.Context, repo string, accessToken string) (*Handler, error) {
	h := &Handler{Repo: repo, AccessToken: accessToken}

	err := h.Init(ctx)

	return h, err

}

func (me *Handler) Helper() provider.CommandHelper[TagNextOutput] {
	return provider.NewHelper[TagNextOutput](me)
}

func (me *Handler) AnyHelper() provider.AnyHelper {
	return provider.NewHelper[TagNextOutput](me)
}

var _ provider.Expressable = (*TagNextOutput)(nil)

func (me *TagNextOutput) Express() map[string]string {
	return map[string]string{
		"major":             me.Major,
		"minor":             me.Minor,
		"patch":             me.Patch,
		"major_minor":       me.MajorMinor,
		"major_minor_patch": me.MajorMinorPatch,
		"full":              me.Full,
		"buildx_tags":       me.BuildxTags,
	}
}
