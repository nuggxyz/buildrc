package buildrc

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/spf13/afero"
	"github.com/walteh/buildrc/pkg/git"
	"gopkg.in/yaml.v3"
)

type Buildrc struct {
	MajorRaw int `yaml:"major,flow" json:"major"`
}

func (me *Buildrc) Major() uint64 {
	return uint64(me.MajorRaw)
}

func LoadBuildrc(ctx context.Context, gitp git.GitProvider) (*Buildrc, error) {

	brc := &Buildrc{}

	buf, err := afero.ReadFile(gitp.Fs(), ".buildrc")
	if err != nil {
		return nil, err
	}

	zerolog.Ctx(ctx).Debug().Str("file", ".buildrc").Str("data", string(buf)).Msg("loaded buildrc")

	err = yaml.Unmarshal(buf, brc)
	if err != nil {
		return nil, err
	}

	return brc, nil
}
