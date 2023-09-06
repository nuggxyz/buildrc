package binary

import (
	"context"
	"fmt"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/walteh/buildrc/pkg/install"
	"github.com/walteh/snake"
)

var _ snake.Snakeable = (*Handler)(nil)

type Handler struct {
	Organization string
	Repository   string
	Version      string
	Token        string
	Provider     string
	OutFile      string
}

func (me *Handler) BuildCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Short: "install buildrc",
	}

	cmd.Args = cobra.ExactArgs(0)

	cmd.PersistentFlags().StringVar(&me.Provider, "provider", "github", "Provider to install from")
	cmd.PersistentFlags().StringVar(&me.Repository, "repository", "", "Repository to install from")
	cmd.PersistentFlags().StringVar(&me.Organization, "organization", "", "Organization to install from")
	cmd.PersistentFlags().StringVar(&me.Version, "version", "latest", "Version to install")
	cmd.PersistentFlags().StringVar(&me.OutFile, "outfile", "", "Output file")

	cmd.PersistentFlags().StringVar(&me.Token, "token", "", "Oauth2 token to use")

	return cmd
}

func (me *Handler) ParseArguments(ctx context.Context, cmd *cobra.Command, file []string) error {

	if me.Repository == "" || me.Organization == "" {
		return fmt.Errorf("Repository and organization must be specified")
	}

	return nil

}

func (me *Handler) Run(ctx context.Context, cmd *cobra.Command) error {
	var fle afero.File
	var err error

	switch me.Provider {
	case "github":
		{
			fle, err = install.DownloadGithubRelease(ctx, afero.NewOsFs(), me.Organization, me.Repository, me.Version, me.Token)
			if err != nil {
				return err
			}
		}
	default:
		{
			return fmt.Errorf("Unknown provider: %s", me.Provider)
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
