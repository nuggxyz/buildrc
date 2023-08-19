package pr

import (
	"context"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/walteh/buildrc/pkg/git"
	"github.com/walteh/snake"
)

var _ snake.Snakeable = (*Handler)(nil)

type Handler struct {
	PR int64 `json:"pr"`
}

func (me *Handler) BuildCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pr",
		Short: "calculate next pre-release tag",
	}

	cmd.Args = cobra.ExactArgs(1)

	return cmd
}

func (me *Handler) ParseArguments(ctx context.Context, cmd *cobra.Command, file []string) error {

	inter, err := strconv.ParseInt(file[0], 10, 64)
	if err != nil {
		return err
	}

	me.PR = inter

	return nil

}

func (me *Handler) Run(ctx context.Context, cmd *cobra.Command, gitp git.GitProvider) error {

	latestHead, err := gitp.GetLatestSemverTagFromRef(ctx, "HEAD")
	if err != nil {
		return err
	}

	revision, err := gitp.GetCurrentShortHashFromRef(ctx, "HEAD")
	if err != nil {
		return err
	}

	work := *latestHead

	work, err = work.SetPrerelease("pr." + strconv.FormatInt(me.PR, 10))
	if err != nil {
		return err
	}

	work, err = work.SetMetadata(revision)
	if err != nil {
		return err
	}

	cmd.Printf("%s\n", work.String())

	return nil
}
