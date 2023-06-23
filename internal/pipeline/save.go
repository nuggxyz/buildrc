package pipeline

import (
	"bytes"
	"context"
	"errors"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/afero"
)

func Load(ctx context.Context, prov Pipeline, cmd string, fs afero.Fs) ([]byte, error) {
	tmp, err := TempFileName(ctx, prov, fs, cmd)
	if err != nil {
		return nil, err
	}

	zerolog.Ctx(ctx).Debug().Str("tmp", tmp).Any("fs", fs).Msg("Load")
	// try to load from tmp folder
	f, err := fs.Open(tmp)
	if err != nil {

		// if not found do nothing
		if errors.Is(err, os.ErrNotExist) {
			return []byte{}, nil
		}

		return nil, err
	}

	z, err := afero.ReadAll(f)
	if err != nil {
		return nil, err
	}

	zerolog.Ctx(ctx).Debug().RawJSON("data", z).Msgf("loaded result from %s", tmp)

	return z, nil
}

func Save(ctx context.Context, prov Pipeline, cmd string, result []byte, fs afero.Fs) error {

	tmp, err := TempFileName(ctx, prov, fs, cmd)
	if err != nil {
		return err
	}

	return afero.WriteReader(fs, tmp, bytes.NewReader(result))
}
