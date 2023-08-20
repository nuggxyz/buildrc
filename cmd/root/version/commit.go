package version

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
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

	return cmd
}

func (me *Handler) ParseArguments(ctx context.Context, cmd *cobra.Command, file []string) error {

	if me.Patch {
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

func (me *Handler) Run(ctx context.Context, cmd *cobra.Command, gitp git.GitProvider, brc *buildrc.Buildrc) error {

	if me.Auto {
		me.Type = CommitTypeRelease
		if gitp.Dirty(ctx) {
			me.Type = CommitTypeLocal
		} else {
			svt, err := gitp.TryGetSemverTag(ctx)
			if err != nil {
				return err
			}

			if svt != nil {
				cmd.Printf("%s\n", svt.String())
				return nil
			}

			n, err := gitp.TryGetPRNumber(ctx)
			if err != nil {
				return err
			}

			me.PRNumber = n
			if me.PRNumber > 0 {
				me.Type = CommitTypePR
			}
		}
	}

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

			if latestHead.Major() < brc.Major {
				latestHead, err = semver.NewVersion(strconv.FormatUint(brc.Major, 10) + ".0.0")
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

			work, err = work.SetPrerelease("pr." + strconv.FormatUint(me.PRNumber, 10))
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
