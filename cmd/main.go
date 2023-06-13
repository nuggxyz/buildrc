package main

import (
	"context"
	"log"
	"os"

	"github.com/nuggxyz/buildrc/cmd/buildrc/build"
	"github.com/nuggxyz/buildrc/cmd/buildrc/load"
	"github.com/nuggxyz/buildrc/cmd/buildrc/packages"
	"github.com/nuggxyz/buildrc/cmd/gen/github"
	"github.com/nuggxyz/buildrc/cmd/tag/list"
	"github.com/nuggxyz/buildrc/cmd/tag/next"

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
	Load     *load.Handler     `cmd:""`
	Packages *packages.Handler `cmd:""`
	Build    *build.Handler    `cmd:""`
	Tag      struct {
		Next *next.Handler `cmd:""`
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

	CLI := CLI{}

	var prov provider.ContentProvider

	prov, err := runner.NewGHActionContentProvider(ctx, file.NewOSFile())

	if err != nil {

		if os.Getenv("BYPASS_CI") == "1" {

			zerolog.Ctx(ctx).Debug().Msg("using mock content provider")

			prov = provider.NewNoopContentProvider(nil)

		} else {
			zerolog.Ctx(ctx).Error().Err(err).Msg("failed to create runner content provider")

			return err
		}
	}

	k := kong.Parse(&CLI,
		kong.BindTo(ctx, (*context.Context)(nil)),
		kong.Name("ci"),
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
