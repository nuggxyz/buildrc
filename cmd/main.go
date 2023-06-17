package main

import (
	"context"
	"log"
	"os"

	"github.com/nuggxyz/buildrc/cmd/buildrc/load"
	packagecmd "github.com/nuggxyz/buildrc/cmd/buildrc/package"
	"github.com/nuggxyz/buildrc/cmd/gen/github"
	"github.com/nuggxyz/buildrc/cmd/release/build"
	"github.com/nuggxyz/buildrc/cmd/release/container"
	"github.com/nuggxyz/buildrc/cmd/release/finalize"
	"github.com/nuggxyz/buildrc/cmd/release/setup"
	"github.com/nuggxyz/buildrc/cmd/release/upload"

	"github.com/nuggxyz/buildrc/cmd/tag/list"

	"github.com/nuggxyz/buildrc/internal/file"
	"github.com/nuggxyz/buildrc/internal/logging"
	"github.com/nuggxyz/buildrc/internal/provider"
	"github.com/nuggxyz/buildrc/internal/runner"

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
	Tag struct {
		List *list.Handler `cmd:""`
	} `cmd:"" help:"tag related commands"`
	Hook struct {
	} `cmd:"" help:"hook related commands"`
	Gen struct {
		Github *github.Handler `cmd:"" help:"generate actions"`
	} `cmd:"" help:"generate actions"`
}

func (me *CLI) AfterApply(ctx context.Context, kctx *kong.Context) error {

	return nil
}

func run() error {

	ctx := context.Background()

	ctx = logging.NewVerboseLoggerContext(ctx)

	cli := CLI{}

	var prov provider.ContentProvider

	prov, err := runner.NewGHActionContentProvider(ctx, file.NewOSFile())
	if err != nil {

		// if os.Getenv("BYPASS_CI") == "1" {

		zerolog.Ctx(ctx).Debug().Msg("using mock content provider")

		prov = provider.NewDefaultContentProvider(file.NewOSFile())

		dir := os.TempDir()

		a, err := os.MkdirTemp(dir, "temp")
		if err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("failed to create temp dir")
			return err
		}

		b, err := os.MkdirTemp(dir, "cache")
		if err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("failed to create cache dir")
			return err
		}

		err = os.Setenv("BUILDRC_TEMP_DIR", a)
		if err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("failed to set temp dir")
			return err
		}

		err = os.Setenv("BUILDRC_CACHE_DIR", b)
		if err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("failed to set cache dir")
			return err
		}

		defer func() {
			err = os.RemoveAll(a)
			if err != nil {
				zerolog.Ctx(ctx).Error().Err(err).Msg("failed to remove temp dir")
			}
			err = os.RemoveAll(b)
			if err != nil {
				zerolog.Ctx(ctx).Error().Err(err).Msg("failed to remove cache dir")
			}

		}()

		// } else {
		// 	zerolog.Ctx(ctx).Error().Err(err).Msg("failed to create runner content provider")

		// 	return err
		// }
	}

	err = provider.SetEnvFromCache(ctx, prov)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("failed to set env from cache")
		return err
	}

	k := kong.Parse(&cli,
		kong.BindTo(ctx, (*context.Context)(nil)),
		kong.Name("buildrc"),
		kong.IgnoreFields("Command"),
		kong.BindTo(prov, (*provider.ContentProvider)(nil)),
	)

	err = k.Run(ctx)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("failed to run kong")
		return err
	}

	return nil
}

func main() {

	err := run()
	if err != nil {
		log.Fatal(err)
	}

}
