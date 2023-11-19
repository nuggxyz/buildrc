package resolvers

import (
	"github.com/spf13/afero"
	"github.com/spf13/pflag"
	"github.com/walteh/snake"
)

var _ snake.Flagged = (*FSResolver)(nil)

type FSResolver struct {
	GitDir string `json:"git-dir"`
}

func (me *FSResolver) Flags(flgs *pflag.FlagSet) {
	flgs.StringVar(&me.GitDir, "git-dir", "", "Git directory")
}

func (me *FSResolver) Run() (afero.Fs, error) {
	osf := afero.NewOsFs()

	if me.GitDir != "" {
		return afero.NewBasePathFs(osf, me.GitDir), nil
	}

	return osf, nil
}
