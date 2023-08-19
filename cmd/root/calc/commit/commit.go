package commit

import (
	"context"
	"strconv"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/spf13/cobra"
	"github.com/walteh/buildrc/pkg/git"
	"github.com/walteh/snake"
)

var _ snake.Snakeable = (*Handler)(nil)

type Handler struct {
	Message string `json:"message"`
	Major   uint64 `json:"major"`

	PatchIndicator string `json:"patch-indicator"`
}

func (me *Handler) BuildCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "calc commit",
		Short: "calculate next pre-release tag",
	}

	cmd.Args = cobra.ExactArgs(1)

	cmd.Flags().StringVarP(&me.PatchIndicator, "patch-indicator", "i", "patch", "The ref to calculate the patch from")
	cmd.Flags().Uint64VarP(&me.Major, "major", "m", 0, "The major version to set")

	return cmd
}

func (me *Handler) ParseArguments(ctx context.Context, cmd *cobra.Command, file []string) error {

	me.Message = file[0]

	return nil

}

func (me *Handler) Run(ctx context.Context, cmd *cobra.Command, gitp git.GitProvider) error {

	patch := strings.Contains(me.Message, me.PatchIndicator)

	latestHead, err := gitp.GetLatestSemverTagFromRef(ctx, "HEAD")
	if err != nil {
		return err
	}

	if latestHead.Major() < me.Major {
		latestHead, err = semver.NewVersion(strconv.FormatUint(me.Major, 10) + ".0.0")
		if err != nil {
			return err
		}
		cmd.Printf("%s\n", latestHead.String())
		return nil
	}

	// we do not care about the prerelease or metadata and this safely removes it
	work := *semver.New(latestHead.Major(), latestHead.Minor(), latestHead.Patch(), "", "")

	if patch {
		work = work.IncPatch()
	} else {
		work = work.IncMinor()
	}

	cmd.Printf("%s\n", work.String())

	return nil
}
