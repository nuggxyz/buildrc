package cache

import (
	"context"
	"path/filepath"

	"github.com/nuggxyz/buildrc/internal/buildrc"
	"github.com/nuggxyz/buildrc/internal/kvstore"
	"github.com/rs/zerolog"
)

func tempFile(ctx context.Context) (string, error) {
	var dir string
	if envvar, err := buildrc.BuildrcTempDir.Load(); err == nil && envvar != "" {
		dir = envvar
	} else {
		return "", err
	}
	return filepath.Join(dir, "temp.db"), nil
}

func HasCacheBeenHit(ctx context.Context, flag string) (bool, error) {
	dir, err := tempFile(ctx)
	if err != nil {
		return false, err
	}

	zerolog.Ctx(ctx).Debug().Str("db", dir).Msg("checking if cache has been hit")
	var res bool
	l, err := kvstore.Load(ctx, dir, "cache", flag, &res)
	if err != nil {
		return false, err
	}

	return l, nil
}

func RecordCacheHit(ctx context.Context, flag string) error {
	dir, err := tempFile(ctx)
	if err != nil {
		return err
	}

	zerolog.Ctx(ctx).Debug().Str("db", dir).Msg("recording cache hit")

	dat := true

	return kvstore.Save(ctx, dir, "cache", flag, &dat)
}
