package buildrc

import (
	"context"
	"path/filepath"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

type Buildrc struct {
	Major uint64 `json:"version" yaml:"version"`
}

func LoadBuildrc(ctx context.Context, fs afero.Fs, dir string) (*Buildrc, error) {
	fle, err := fs.Open(filepath.Join(dir, ".buildrc"))
	if err != nil {
		return nil, err
	}

	defer fle.Close()

	brc := &Buildrc{}

	var buf []byte

	_, err = fle.Read(buf)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(buf, brc)
	if err != nil {
		return nil, err
	}

	return brc, nil
}
