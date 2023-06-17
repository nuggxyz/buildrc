package upload

import (
	"context"
	"fmt"
	"time"

	"github.com/nuggxyz/buildrc/cmd/buildrc/load"
	"github.com/nuggxyz/buildrc/internal/buildrc"
	"github.com/nuggxyz/buildrc/internal/github"
	"github.com/nuggxyz/buildrc/internal/provider"

	"github.com/rs/zerolog"
)

const (
	CommandID = "upload"
)

type Handler struct {
	File string `flag:"file" type:"file:" default:".buildrc"`
}

func (me *Handler) Run(ctx context.Context, cp provider.ContentProvider) (err error) {
	_, err = me.Build(ctx, cp)
	return err
}

func (me *Handler) Build(ctx context.Context, cp provider.ContentProvider) (out *any, err error) {

	return provider.Wrap(CommandID, me.build)(ctx, cp)
}

func (me *Handler) build(ctx context.Context, prv provider.ContentProvider) (out *any, err error) {

	brc, err := load.NewHandler(me.File).Load(ctx, prv)
	if err != nil {
		return nil, err
	}

	ghclient, err := github.NewGithubClient(ctx, "", "")
	if err != nil {
		return nil, err
	}

	ok, reason, err := ghclient.ShouldBuild(ctx)
	if err != nil {
		return nil, err
	}

	if !ok {
		zerolog.Ctx(ctx).Info().Str("reason", reason).Msg("build not required")
		return nil, nil
	}

	err = me.run(ctx, ghclient, brc)
	if err != nil {
		return nil, err
	}

	return nil, nil

}

func (me *Handler) run(ctx context.Context, clnt *github.GithubClient, brc *buildrc.BuildRC) error {
	return buildrc.RunAllPackages(ctx, brc, 10*time.Minute, func(ctx context.Context, pkg *buildrc.Package, arc buildrc.Platform) error {

		file, err := arc.OutputFile(pkg)
		if err != nil {
			return fmt.Errorf("error running upload with [%s:%s]: %v", arc.OS(), arc.Arch(), err)
		}

		zerolog.Ctx(ctx).Debug().Msgf("wrote SHA-256 checksum to %s.sha256", file)

		err = clnt.Upload(ctx, file+".tar.gz")
		if err != nil {
			return fmt.Errorf("error uploading archive: %v", err)
		}

		err = clnt.Upload(ctx, file+".sha256")
		if err != nil {
			return fmt.Errorf("error uploading checksum: %v", err)
		}

		zerolog.Ctx(ctx).Debug().Msgf("uploaded checksum %s.sha256", file)

		return nil
	})

}
