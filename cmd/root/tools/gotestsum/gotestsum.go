package gotestsum

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/walteh/buildrc/pkg/git"
	"github.com/walteh/snake"

	gotestsumcmd "gotest.tools/gotestsum/cmd"
)

var _ snake.Snakeable = (*Handler)(nil)

type Handler struct {
	args []string
}

func (me *Handler) BuildCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Short: "get current revision",
	}

	return cmd
}

func (me *Handler) ParseArguments(ctx context.Context, cmd *cobra.Command, file []string) error {
	me.args = file
	return nil

}

func (me *Handler) Run(ctx context.Context, cmd *cobra.Command, gitp git.GitProvider) error {
	return gotestsumcmd.Run("gotestsum", me.args)
}
