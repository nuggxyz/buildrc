package next

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/walteh/buildrc/internal/buildrc"
	"github.com/walteh/buildrc/internal/git"
	"github.com/walteh/snake"
)

var _ snake.Snakeable = (*Handler)(nil)

type Handler struct {
}

func (me *Handler) BuildCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "next",
		Short: "calculate next pre-release tag",
	}

	return cmd
}

func (me *Handler) ParseArguments(ctx context.Context, cmd *cobra.Command, file []string) error {

	return nil

}

func (me *Handler) Run(ctx context.Context, cmd *cobra.Command, gitp git.GitProvider, brc *buildrc.Buildrc, prp git.PullRequestProvider) error {

	targetSemver, err := git.CalculateNextPreReleaseTag(ctx, brc, gitp, prp)
	if err != nil {
		return err
	}

	cmd.Printf("%s\n", targetSemver.String())

	return nil
}
