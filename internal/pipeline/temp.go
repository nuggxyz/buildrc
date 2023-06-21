package pipeline

import (
	"context"
	"path/filepath"

	"github.com/nuggxyz/buildrc/internal/kvstore"
	"github.com/rs/zerolog"
	"github.com/spf13/afero"
)

func tempFile(ctx context.Context, p Pipeline, fs afero.Fs) (string, error) {
	var dir string
	if envvar, err := BuildrcTempDir.Load(ctx, p, fs); err == nil && envvar != "" {
		dir = envvar
	} else {
		return "", err
	}
	return filepath.Join(dir, "temp.db"), nil
}

func HasCacheBeenHit(ctx context.Context, p Pipeline, fs afero.Fs, flag string) (bool, error) {
	dir, err := tempFile(ctx, p, fs)
	if err != nil {
		return false, err
	}

	zerolog.Ctx(ctx).Debug().Str("db", dir).Msg("checking if cache has been hit")
	res := false
	l, err := kvstore.Load(ctx, fs, dir, flag, &res)
	if err != nil {
		if kvstore.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}

	return l, nil
}

func RecordCacheHit(ctx context.Context, p Pipeline, fs afero.Fs, flag string) error {
	dir, err := tempFile(ctx, p, fs)
	if err != nil {
		return err
	}

	zerolog.Ctx(ctx).Debug().Str("db", dir).Msg("recording cache hit")

	dat := true

	return kvstore.Save(ctx, fs, dir, flag, &dat)
}
