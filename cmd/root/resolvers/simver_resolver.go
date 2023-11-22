package resolvers

import (
	"context"

	"github.com/spf13/pflag"
	"github.com/walteh/buildrc/pkg/git"
	"github.com/walteh/simver"
	"github.com/walteh/simver/gitexec"
	"github.com/walteh/snake"
)

var _ snake.Flagged = (*SimverResolver)(nil)

type SimverResolver struct {
	githubActions bool
	cacheDir      string
}

func (me *SimverResolver) Flags(gg *pflag.FlagSet) {
	gg.BoolVar(&me.githubActions, "github-actions", false, "use github actions providers")
	gg.StringVar(&me.cacheDir, "cache-dir", "", "cache directory")
}

func (me *SimverResolver) Run(ctx context.Context, prov git.GitProvider) (simver.Execution, simver.GitProvider, simver.TagReader, error) {

	if me.githubActions {

		g, tr, _, _, prr, err := gitexec.BuildGitHubActionsProviders()
		if err != nil {
			return nil, nil, nil, err
		}

		ex, _, err := simver.LoadExecutionFromPR(ctx, tr, prr)
		if err != nil {
			return nil, nil, nil, err
		}

		return ex, g, tr, nil
	} else {
		g, tr, _, _, err := gitexec.BuildLocalProviders(prov.Fs())
		if err != nil {
			return nil, nil, nil, err
		}

		ex, err := simver.NewLocalProjectState(ctx, g, tr)
		if err != nil {
			return nil, nil, nil, err
		}

		return ex, g, tr, nil
	}

}
