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
	"github.com/nuggxyz/buildrc/internal/github"
	"github.com/nuggxyz/buildrc/internal/pipeline"

	"github.com/nuggxyz/buildrc/internal/file"
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

	execgit := git.NewExecGitProvider()

	var pr git.PullRequestProvider
	var release git.ReleaseProvider
	var repometa git.RepositoryMetadataProvider
	var prov pipeline.Pipeline

	if pipeline.IAmInAGithubAction(ctx) {
		provd, err := pipeline.NewGithubActionPipeline(ctx, file.NewOSFile())
		if err != nil {

			zerolog.Ctx(ctx).Error().Err(err).Msg("failed to create content provider")
			return err
		}

		prov = provd

		ghp, err := github.NewGithubClient(ctx, "", "")
		if err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("failed to create github client")
			return err
		}

		pr = ghp
		release = ghp
		repometa = ghp

	}

	k := kong.Parse(&cli,
		kong.BindTo(ctx, (*context.Context)(nil)),
		kong.Name("buildrc"),
		kong.IgnoreFields("Command"),
	)

	if k.Selected().Name == "version" {
		return k.Run(ctx)
	}

	res, err := pipeline.WrapGeneric(ctx, "main-buildrc", prov, nil, func(ctx context.Context, a any) (*buildrc.Buildrc, error) {
		return buildrc.Parse(ctx, cli.File)
	})
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("failed to parse buildrc")
		return err
	}

	cmp := common.NewProvider(execgit, release, prov, pr, res, repometa)

	err = pipeline.SetEnvFromCache(ctx, prov)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("failed to set env from cache")
		return err
	}

	err = k.Run(ctx, kong.BindTo(cmp, (common.Provider)(nil)))
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
