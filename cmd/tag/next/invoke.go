package next

import (
	"context"
	"fmt"

	buildrc "github.com/nuggxyz/buildrc/cmd/buildrc/load"
	"github.com/nuggxyz/buildrc/cmd/tag/list"
	"github.com/nuggxyz/buildrc/internal/provider"
)

type Handler struct {
	Repo        string `arg:"repo" type:"repo:" required:"true"`
	BuildrcFile string `arg:"buildrc_file" type:"file:" required:"true"`
	AccessToken string `arg:"access_token" type:"access_token:" required:"true"`

	gettagsHandler *list.Handler
	buildrcHandler *buildrc.Handler
}

func (me *Handler) Init(ctx context.Context) (err error) {
	me.gettagsHandler, err = list.NewHandler(ctx, me.Repo, me.AccessToken)
	if err != nil {
		return err
	}

	me.buildrcHandler, err = buildrc.NewHandler(ctx, me.BuildrcFile)
	if err != nil {
		return err
	}
	return
}

func (me *Handler) Invoke(ctx context.Context, prv provider.ContentProvider) (out *output, err error) {

	prov, err := me.gettagsHandler.Helper().Run(ctx, prv)
	if err != nil {
		return nil, err
	}

	brc, err := me.buildrcHandler.Helper().Run(ctx, prv)
	if err != nil {
		return nil, err
	}

	// Increment patch version
	nextVersion := prov.HighestVersion.IncPatch()

	// If the buildrc version is higher than the highest version, use that instead
	if brc.Version.GreaterThan(&nextVersion) {
		nextVersion = *brc.Version
	}

	return &output{
		Major:           fmt.Sprintf("%d", nextVersion.Major()),
		Minor:           fmt.Sprintf("%d", nextVersion.Minor()),
		Patch:           fmt.Sprintf("%d", nextVersion.Patch()),
		MajorMinor:      fmt.Sprintf("%d.%d", nextVersion.Major(), nextVersion.Minor()),
		MajorMinorPatch: fmt.Sprintf("%d.%d.%d", nextVersion.Major(), nextVersion.Minor(), nextVersion.Patch()),
		Full:            nextVersion.String(),
	}, nil
}
