package resolvers

import (
	"github.com/spf13/afero"
	"github.com/spf13/pflag"
	"github.com/walteh/snake"
)

var _ snake.Flagged = (*FSResolver)(nil)

type FSResolver struct {
}

func (me *FSResolver) Flags(_ *pflag.FlagSet) {
}

func (me *FSResolver) Run() (afero.Fs, error) {
	osf := afero.NewOsFs()

	return osf, nil
}
