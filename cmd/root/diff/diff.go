package diff

import (
	"context"
	"os"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/walteh/buildrc/pkg/file"
	"github.com/walteh/snake"
)

var _ snake.Snakeable = (*Handler)(nil)

type Handler struct {
	current string   // real directory
	correct string   // real directory
	globs   []string // glob pattern
}

func (me *Handler) BuildCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Short: "get current revision",
	}

	cmd.Args = cobra.ExactArgs(0)

	cmd.Flags().StringVarP(&me.current, "current", "c", ".", "current directory")
	cmd.Flags().StringVarP(&me.correct, "correct", "r", ".", "correct directory")
	cmd.Flags().StringArrayVarP(&me.globs, "glob", "g", []string{"**/*"}, "glob pattern")

	return cmd
}

func (me *Handler) ParseArguments(ctx context.Context, cmd *cobra.Command, args []string) error {

	return nil

}

func (me *Handler) Run(ctx context.Context, cmd *cobra.Command, gitp afero.Fs) error {

	diffs, err := file.Diff(ctx, gitp, me.current, me.correct, me.globs)
	if err != nil {
		return err
	}

	if len(diffs) > 0 {
		cmd.PrintErrln("============= buildrc diff ==============")
		cmd.PrintErrf(" %d DIFFERENCES FOUND\n", len(diffs))
		cmd.PrintErrln("=========================================")
		cmd.PrintErrf("current:  %s\n", me.current)
		cmd.PrintErrf("correct:  %s\n", me.correct)
		cmd.PrintErrln("=========================================")
		for _, diff := range diffs {
			cmd.PrintErrf("%s\n", diff)
		}
		cmd.PrintErrln("=========================================")
		os.Exit(1)
	} else {
		os.Exit(0)
	}

	return nil
}
