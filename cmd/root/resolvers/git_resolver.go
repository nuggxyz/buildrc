package resolvers

import (
	"github.com/spf13/afero"
	"github.com/spf13/pflag"
	"github.com/walteh/buildrc/pkg/git"
	"github.com/walteh/snake"
)

var _ snake.Flagged = (*GitResolver)(nil)

type GitResolver struct {
	GitDir string `json:"git-dir"`
}

func (me *GitResolver) Flags(flgs *pflag.FlagSet) {
	flgs.StringVar(&me.GitDir, "git-dir", "", "Git directory")
}

func (me *GitResolver) Run(fls afero.Fs) (git.GitProvider, error) {

	gitp, err := git.NewGitGoGitProvider(fls, me.GitDir)
	if err != nil {
		return nil, err
	}

	return gitp, nil
}
