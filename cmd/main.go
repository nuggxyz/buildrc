package main

import (
	"context"
	"log"
	"os"

	"github.com/nuggxyz/buildrc/cmd/load"
	"github.com/nuggxyz/buildrc/cmd/simple/close"
	"github.com/nuggxyz/buildrc/cmd/simple/docker"
	"github.com/nuggxyz/buildrc/cmd/simple/open"
	"github.com/nuggxyz/buildrc/cmd/simple/upload"
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
	Load   *load.Handler `cmd:""`
	Simple struct {
		Docker *docker.Handler `cmd:"" help:"build a docker image"`
		Open   *open.Handler   `cmd:"" help:"open a url"`
		Close  *close.Handler  `cmd:"" help:"close a url"`
		Upload *upload.Handler `cmd:"" help:"upload a file"`
	} `cmd:"" help:"build a tag"`
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
		ctx = logging.NewVerboseLoggerContextWithLevel(ctx, zerolog.TraceLevel)
	} else {
		ctx = logging.NewVerboseLoggerContextWithLevel(ctx, zerolog.Disabled)
	}

	cli := CLI{}

	k := kong.Parse(&cli, kong.Name("buildrc"))

	if k.Selected().Name == "version" {
		k.BindTo(ctx, (*context.Context)(nil))
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

		pr = git.NewMemoryPullRequestProvider([]*git.PullRequest{
			{
				Number: 1,
				Open:   true,
			},
		})
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

	zerolog.Ctx(ctx).Debug().Interface("buildrc", res).Msg("parsed buildrc")

	prov2 := common.NewProvider(execgit, release, pipe, pr, res, repometa, fs)

	err = pipeline.SetEnvFromCache(ctx, pipe, fs)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("failed to set env from cache")
		return err
	}

	k.BindTo(pipe, (*pipeline.Pipeline)(nil))
	k.BindTo(fs, (*afero.Fs)(nil))
	k.BindTo(res, (*buildrc.Provider)(nil))
	k.BindTo(execgit, (*git.GitProvider)(nil))
	k.BindTo(pr, (*git.PullRequestProvider)(nil))
	k.BindTo(release, (*git.ReleaseProvider)(nil))
	k.BindTo(repometa, (*git.RemoteRepositoryMetadataProvider)(nil))
	k.BindTo(prov2, (*common.Provider)(nil))
	k.BindTo(ctx, (*context.Context)(nil))

	err = k.Run(ctx)
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
