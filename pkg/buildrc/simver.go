package buildrc

import (
	"context"
	"strings"

	"github.com/walteh/simver/gitexec"
	"golang.org/x/mod/semver"
)

func GetVersionWithSimver(ctx context.Context, def string) (string, string, error) {
	gitp, err := gitexec.NewLocalReadOnlyGitProvider("git", ".")
	if err != nil {
		return "", "", err
	}

	tagp, err := gitexec.NewLocalReadOnlyTagProvider("git", ".")
	if err != nil {
		return "", "", err
	}

	head, err := gitp.Branch(ctx)
	if err != nil {
		return "", "", err
	}

	tags, err := tagp.TagsFromBranch(ctx, head)
	if err != nil {
		return "", "", err
	}

	tgstrs := tags.SemversMatching(func(s string) bool {
		return semver.IsValid(s) && !strings.Contains(s, "base") && !strings.Contains(s, "reserved")
	})

	semver.Sort(tgstrs)

	if len(tgstrs) == 0 {
		return def, head, nil
	}

	return tgstrs[len(tgstrs)-1], head, nil

}
