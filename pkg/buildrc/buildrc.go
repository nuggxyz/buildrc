package buildrc

import (
	"context"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

type Buildrc struct {
	MajorRaw int `yaml:"major,flow" json:"major"`
}

func (me *Buildrc) Major() uint64 {
	return uint64(me.MajorRaw)
}

func LoadBuildrc(ctx context.Context, fs afero.Fs, dir string) (*Buildrc, error) {
	fle, err := fs.Open(filepath.Join(dir, ".buildrc"))
	if err != nil {
		return nil, err
	}

	defer fle.Close()

	brc := &Buildrc{}

	buf, err := afero.ReadFile(fs, filepath.Join(dir, ".buildrc"))
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
