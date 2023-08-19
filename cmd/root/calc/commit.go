package calc

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/spf13/cobra"
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
	Major                 uint64     `json:"major"`
	PatchIndicator        string     `json:"patch-indicator"`
	PRNumber              int64      `json:"pr-number"`
	CommitMessageOverride string     `json:"commit-message-override"`
	LatestTagOverride     string     `json:"latest-tag-override"`
	Patch                 bool       `json:"patch"`
}

func (me *Handler) BuildCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Short: "calculate next pre-release tag",
	}

	cmd.Args = cobra.ExactArgs(0)

	cmd.Flags().StringVarP(&me.PatchIndicator, "patch-indicator", "i", "patch", "The ref to calculate the patch from")
	cmd.Flags().Uint64VarP(&me.Major, "major", "m", 0, "The major version to set")
	cmd.Flags().StringVarP((*string)(&me.Type), "type", "t", "local", "The type of commit to calculate")
	cmd.Flags().Int64VarP(&me.PRNumber, "pr-number", "n", 0, "The pr number to set")
	cmd.Flags().StringVarP(&me.CommitMessageOverride, "commit-message-override", "c", "", "The commit message to use")
	cmd.Flags().StringVarP(&me.LatestTagOverride, "latest-tag-override", "l", "", "The tag to use")

	cmd.Flags().BoolVarP(&me.Patch, "patch", "p", false, "shortcut for --type=release --patch-indicator=x --commit-message-override=x")

	return cmd
}

func (me *Handler) ParseArguments(ctx context.Context, cmd *cobra.Command, file []string) error {

	if me.Patch {
		me.Type = CommitTypeRelease
		me.PatchIndicator = "patch"
		me.CommitMessageOverride = "patch"
	}

	if me.Type == CommitTypePR {
		if me.PRNumber == 0 {
			return fmt.Errorf("'--pr-number=#' is required for type %s", me.Type)
		}
	}

	return nil

}

func (me *Handler) Run(ctx context.Context, cmd *cobra.Command, gitp git.GitProvider) error {

	switch me.Type {
	case CommitTypeRelease:
		{

			var latestHead *semver.Version
			var message string
			var err error

			if me.LatestTagOverride != "" {
				latestHead, err = semver.NewVersion(me.LatestTagOverride)
				if err != nil {
					return err
				}
			} else {
				latestHead, err = gitp.GetLatestSemverTagFromRef(ctx, "HEAD")
				if err != nil {
					return err
				}
			}

			if me.CommitMessageOverride != "" {
				message = me.CommitMessageOverride
			} else {
				message, err = gitp.GetCurrentCommitMessageFromRef(ctx, "HEAD")
				if err != nil {
					return err
				}
			}

			patch := strings.Contains(message, me.PatchIndicator)

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

		}
	case CommitTypeLocal:
		{
			work := *semver.New(0, 0, 0, "local", time.Now().Format("2006-01-02-15:04:05"))
			cmd.Printf("%s\n", work.String())
		}
	case CommitTypePR:
		{

			latestHead, err := gitp.GetLatestSemverTagFromRef(ctx, "HEAD")
			if err != nil {
				return err
			}

			revision, err := gitp.GetCurrentShortHashFromRef(ctx, "HEAD")
			if err != nil {
				return err
			}

			work := *latestHead

			work, err = work.SetPrerelease("pr." + strconv.FormatInt(me.PRNumber, 10))
			if err != nil {
				return err
			}

			work, err = work.SetMetadata(revision)
			if err != nil {
				return err
			}

			cmd.Printf("%s\n", work.String())
		}
	}

	return nil
}
