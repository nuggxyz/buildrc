package binary

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/walteh/buildrc/pkg/buildrc"
	"github.com/walteh/buildrc/pkg/install"
	"github.com/walteh/snake"
)

var _ snake.Flagged = (*Handler)(nil)

type Handler struct {
	Organization string
	Repository   string
	Version      string
	Token        string
	Provider     string
	OutFile      string
	Platform     string
}

func (me *Handler) Flags(flgs *pflag.FlagSet) {
	flgs.StringVar(&me.Organization, "organization", "", "Organization")
	flgs.StringVar(&me.Repository, "repository", "", "Repository")
	flgs.StringVar(&me.Version, "binary-version", "", "Version")
	flgs.StringVar(&me.Token, "token", "", "Token")
	flgs.StringVar(&me.Provider, "provider", "github", "Provider")
	flgs.StringVar(&me.OutFile, "out-file", "", "OutFile")
	flgs.StringVar(&me.Platform, "platform", "runtime.GOOS/runtime.GOARCH", "Platform")
}

func (me *Handler) Cobra() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download",
		Short: "downloads a binary",
	}

	return cmd
}

func (me *Handler) Run(ctx context.Context) error {

	if me.Repository == "" || me.Organization == "" {
		return errors.Errorf("Repository and organization must be specified")
	}

	var fle afero.File
	var err error

	switch me.Provider {
	case "github":
		{

			var plat *buildrc.Platform
			if me.Platform == "runtime.GOOS/runtime.GOARCH" {
				plat = buildrc.GetGoPlatform(ctx)
			} else {

				plat, err = buildrc.NewPlatformFromFullString(me.Platform)
				if err != nil {
					return err
				}
			}

			fle, err = install.DownloadGithubReleaseWithOptions(ctx, afero.NewOsFs(), &install.DownloadGithubReleaseOptions{
				Org:      me.Organization,
				Name:     me.Repository,
				Version:  me.Version,
				Token:    me.Token,
				Platform: plat,
			})
			if err != nil {
				return err
			}
		}
	default:
		{
			return errors.Errorf("Unknown provider: %s", me.Provider)
		}
	}

	defer fle.Close()

	fls := afero.NewOsFs()

	err = afero.WriteReader(fls, me.OutFile, fle)
	if err != nil {
		return err
	}

	return fls.Chmod(me.OutFile, 0755)

}
