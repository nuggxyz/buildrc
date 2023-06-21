package main

import (
	"context"
	"log"
	"os"

	"github.com/nuggxyz/buildrc/cmd/buildrc/load"
	packagecmd "github.com/nuggxyz/buildrc/cmd/buildrc/package"
	"github.com/nuggxyz/buildrc/cmd/release/build"
	"github.com/nuggxyz/buildrc/cmd/release/container"
	"github.com/nuggxyz/buildrc/cmd/release/finalize"
	"github.com/nuggxyz/buildrc/cmd/release/setup"
	"github.com/nuggxyz/buildrc/cmd/release/upload"
	"github.com/nuggxyz/buildrc/internal/buildrc"
	"github.com/nuggxyz/buildrc/internal/common"
	"github.com/nuggxyz/buildrc/internal/github/actions"
	"github.com/nuggxyz/buildrc/internal/github/restapi"
	"github.com/nuggxyz/buildrc/internal/pipeline"
	"github.com/spf13/afero"

	"github.com/nuggxyz/buildrc/internal/git"
	"github.com/nuggxyz/buildrc/internal/logging"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

type CLI struct {
	Load    *load.Handler       `cmd:""`
	Package *packagecmd.Handler `cmd:""`
	Release struct {
		Build     *build.Handler     `cmd:""`
		Setup     *setup.Handler     `cmd:""`
		Finalize  *finalize.Handler  `cmd:""`
		Upload    *upload.Handler    `cmd:""`
		Container *container.Handler `cmd:""`
	} `cmd:"" help:"release related commands"`
	Version *VersionHandler `cmd:"" help:"show version"`
	Quiet   bool            `flag:"" help:"enable quiet logging" short:"q"`
	File    string          `flag:"file" type:"file:" default:".buildrc"`
}

func run() error {

	quiet := false
	for _, arg := range os.Args {
		if arg == "--quiet" || arg == "-q" {
			quiet = true
		}
	}

	ctx := context.Background()

	if !quiet {
		ctx = logging.NewVerboseLoggerContext(ctx)
	} else {
		ctx = logging.NewVerboseLoggerContextWithLevel(ctx, zerolog.Disabled)
	}

	cli := CLI{}

	k := kong.Parse(&cli, kong.Name("buildrc"))

	if k.Selected().Name == "version" {
		return k.Run(ctx)
	}

	execgit := git.NewGitGoGitProvider()

	var pr git.PullRequestProvider
	var release git.ReleaseProvider
	var repometa git.RemoteRepositoryMetadataProvider
	var pipe pipeline.Pipeline
	var fs afero.Fs

	if actions.IAmInAGithubAction(ctx) {
		actionpipe, err := actions.NewGithubActionPipeline(ctx)
		if err != nil {

			zerolog.Ctx(ctx).Error().Err(err).Msg("failed to create content provider")
			return err
		}

		pipe = actionpipe

		ghrestapi, err := restapi.NewGithubClient(ctx, execgit, actionpipe)
		if err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("failed to create github client")
			return err
		}

		pr = ghrestapi
		release = ghrestapi
		repometa = ghrestapi

		fs = afero.NewOsFs()
	} else {
		zerolog.Ctx(ctx).Warn().Msg("not running in github action, using local filesystem")
		pipe = pipeline.NewMemoryPipeline()
		fs = afero.NewMemMapFs()

		pr = git.NewMemoryPullRequestProvider([]*git.PullRequest{})
		release = git.NewMemoryReleaseProvider([]*git.Release{})
		repometa = git.NewMemoryRepoMetadataProvider(&git.RemoteRepositoryMetadata{
			Description: "test repo",
			Homepage:    "test.com",
			License:     "Nunya",
		})
	}

	res, err := pipeline.WrapGeneric(ctx, "main-buildrc", pipe, fs, nil, func(ctx context.Context, a any) (*buildrc.Buildrc, error) {
		return buildrc.Parse(ctx, cli.File)
	})

	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("failed to parse buildrc")
		return err
	}

	prov2 := common.NewProvider(execgit, release, pipe, pr, res, repometa, fs)

	// err = pipeline.EnsureCacheDB(ctx, pipe, fs)
	// if err != nil {
	// 	zerolog.Ctx(ctx).Error().Err(err).Msg("failed to ensure cache db")
	// 	return err
	// }

	err = pipeline.SetEnvFromCache(ctx, pipe, fs)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("failed to set env from cache")
		return err
	}

	k.BindTo(ctx, (*context.Context)(nil))
	k.BindTo(prov2, (*common.Provider)(nil))

	err = k.Run(ctx, prov2)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("pipeline failed")
		return err
	}

	return nil

}

func main() {
	if run() != nil {
		os.Exit(1)
	}
}
