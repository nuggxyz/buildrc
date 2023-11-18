package revision

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/walteh/buildrc/pkg/buildrc"
	"github.com/walteh/buildrc/pkg/git"
)

type Handler struct {
}

func (me *Handler) Cobra() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revision",
		Short: "get current revision",
	}

	cmd.Args = cobra.ExactArgs(0)

	return cmd
}

func (me *Handler) Run(ctx context.Context, cmd *cobra.Command, gitp git.GitProvider) error {

	revision, err := buildrc.GetRevision(ctx, gitp)
	if err != nil {
		return err
	}

	cmd.Printf("%s\n", revision)

	return nil
}
