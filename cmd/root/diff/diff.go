package diff

import (
	"context"
	"os"

	"github.com/rs/zerolog"
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
	glob    []string // glob pattern
}

func (me *Handler) BuildCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Short: "get current revision",
	}

	cmd.Args = cobra.ExactArgs(0)

	cmd.Flags().StringVarP(&me.current, "current", "c", ".", "current directory")
	cmd.Flags().StringVarP(&me.correct, "correct", "r", ".", "correct directory")
	cmd.Flags().StringSliceVar(&me.globs, "globs", []string{}, "glob pattern")
	cmd.Flags().StringArrayVar(&me.glob, "glob", []string{}, "glob pattern")

	return cmd
}

func (me *Handler) ParseArguments(ctx context.Context, cmd *cobra.Command, args []string) error {

	me.globs = append(me.globs, me.glob...)

	if len(me.globs) == 0 {
		me.globs = append(me.globs, "**/*")
	}

	return nil

}

func (me *Handler) Run(ctx context.Context, cmd *cobra.Command, gitp afero.Fs) error {

	zerolog.Ctx(ctx).Debug().Str("current", me.current).Str("correct", me.correct).Strs("globs", me.globs).Msg("diff")

	diffs, err := file.Diff(ctx, gitp, me.current, me.correct, me.globs)
	if err != nil {
		return err
	}

	notignored, err := file.FilterGitIgnored(ctx, gitp, diffs)
	if err != nil {
		return err
	}

	if len(notignored) > 0 {
		cmd.PrintErrln("============= buildrc diff ==============")
		cmd.PrintErrf(" %d DIFFERENCES FOUND\n", len(notignored))
		cmd.PrintErrln("=========================================")
		cmd.PrintErrf("dir:  %s\n", me.current)
		cmd.PrintErrln("=========================================")
		total := 0
		for _, diff := range notignored {
			cmd.PrintErrf("%s\n", diff)
			total++
			if total > 10 {
				cmd.PrintErrf("... and %d more\n", len(notignored)-total)
				break
			}
		}
		cmd.PrintErrln("=========================================")
		os.Exit(1)
	} else {
		os.Exit(0)
	}

	return nil
}
