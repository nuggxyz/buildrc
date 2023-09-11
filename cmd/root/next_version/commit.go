package next_version

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/walteh/buildrc/pkg/buildrc"
	"github.com/walteh/buildrc/pkg/git"
	"github.com/walteh/snake"
)

var _ snake.Snakeable = (*Handler)(nil)

type CommitType string

const (
	CommitTypePR      CommitType = "pr"
	CommitTypeLocal   CommitType = "local"
	CommitTypeRelease CommitType = "release"
)

type Handler struct {
	Type                  CommitType `json:"type"`
	PatchIndicator        string     `json:"patch-indicator"`
	PRNumber              uint64     `json:"pr-number"`
	CommitMessageOverride string     `json:"commit-message-override"`
	LatestTagOverride     string     `json:"latest-tag-override"`
	Patch                 bool       `json:"patch"`
	Auto                  bool       `json:"auto"`
	NoV                   bool       `json:"no-v"`
}

func (me *Handler) BuildCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Short: "calculate next pre-release tag",
	}

	cmd.Args = cobra.ExactArgs(0)

	cmd.Flags().StringVarP(&me.PatchIndicator, "patch-indicator", "i", "patch", "The ref to calculate the patch from")
	cmd.Flags().StringVarP((*string)(&me.Type), "type", "t", "local", "The type of commit to calculate")
	cmd.Flags().Uint64VarP(&me.PRNumber, "pr-number", "n", 0, "The pr number to set")
	cmd.Flags().StringVarP(&me.CommitMessageOverride, "commit-message-override", "c", "", "The commit message to use")
	cmd.Flags().StringVarP(&me.LatestTagOverride, "latest-tag-override", "l", "", "The tag to use")

	cmd.Flags().BoolVarP(&me.Patch, "patch", "p", false, "shortcut for --patch-indicator=x --commit-message-override=x")

	cmd.Flags().BoolVarP(&me.Auto, "auto", "a", false, "shortcut for if CI != 'true' then local else if '--pr-number' > 0 then pr")

	cmd.Flags().BoolVarP(&me.NoV, "no-v", "", false, "do not prefix with 'v'")

	return cmd
}

func (me *Handler) ParseArguments(ctx context.Context, cmd *cobra.Command, file []string) error {

	return nil

}

func (me *Handler) Run(ctx context.Context, cmd *cobra.Command, gitp git.GitProvider) error {

	brc, err := buildrc.LoadBuildrc(ctx, gitp)
	if err != nil {
		return err
	}

	vers, err := buildrc.GetVersion(ctx, gitp, brc, &buildrc.GetVersionOpts{
		Type:                  buildrc.CommitType(me.Type),
		PatchIndicator:        me.PatchIndicator,
		PRNumber:              me.PRNumber,
		CommitMessageOverride: me.CommitMessageOverride,
		LatestTagOverride:     me.LatestTagOverride,
		Patch:                 me.Patch,
		Auto:                  me.Auto,
		ExcludeV:              me.NoV,
	})

	if err != nil {
		return err
	}

	cmd.Printf("%s\n", vers)

	return nil

}
