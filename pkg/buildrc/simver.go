package buildrc

import (
	"context"
	"strings"

	"github.com/walteh/simver"
	"golang.org/x/mod/semver"
)

func GetVersionWithSimver(ctx context.Context, def string, exc simver.Execution, gp simver.GitProvider) (string, string, error) {

	c := exc.HeadCommitTags()

	var head string

	if len(c) > 0 {
		head = c[len(c)-1].Ref
	}

	tgstrs := c.SemversMatching(func(s string) bool {
		return semver.IsValid(s) && !strings.Contains(s, "base") && !strings.Contains(s, "reserved")
	})

	semver.Sort(tgstrs)

	if len(tgstrs) == 0 {
		gp, err := gp.GetHeadRef(ctx)
		if err != nil {
			return "", "", err
		}
		return def, gp, nil
	}

	return tgstrs[len(tgstrs)-1], head, nil

}
