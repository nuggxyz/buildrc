package revision

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/walteh/buildrc/pkg/buildrc"
	"github.com/walteh/buildrc/pkg/git"
	"github.com/walteh/snake"
)

var _ snake.Snakeable = (*Handler)(nil)

type Handler struct {
}

func (me *Handler) BuildCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Short: "get current revision",
	}

	cmd.Args = cobra.ExactArgs(0)

	return cmd
}

func (me *Handler) ParseArguments(ctx context.Context, cmd *cobra.Command, file []string) error {

	return nil

}

func (me *Handler) Run(ctx context.Context, cmd *cobra.Command, gitp git.GitProvider) error {

	revision, err := buildrc.GetRevision(ctx, gitp)
	if err != nil {
		return err
	}

	cmd.Printf("%s\n", revision)

	return nil
}
