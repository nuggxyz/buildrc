package cache

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/go-github/v53/github"
	"github.com/nuggxyz/buildrc/internal/kvstore"
	"github.com/rs/zerolog"
)

const (
	CACHE_DIR_ENV_VAR = "BUILDRC_CACHE_DIR"
)

func EnsureCacheDB(ctx context.Context) error {
	dir, err := cacheFile(ctx)
	if err != nil {
		return err
	}

	zerolog.Ctx(ctx).Debug().Str("db", dir).Msg("ensuring cache db")

	return kvstore.EnsureStoreFile(ctx, dir)
}

func cacheFile(ctx context.Context) (string, error) {
	var dir string
	if envvar := os.Getenv(CACHE_DIR_ENV_VAR); envvar != "" {
		dir = envvar
	} else {
		return "", fmt.Errorf("cache dir not set, please set %s", CACHE_DIR_ENV_VAR)
	}
	return filepath.Join(dir, dir, "cache.db"), nil
}

func SaveRelease(ctx context.Context, name string, r *github.RepositoryRelease) error {
	dir, err := cacheFile(ctx)
	if err != nil {
		return err
	}

	zerolog.Ctx(ctx).Debug().Str("name", name).Str("db", dir).Msg("saving release to cache")

	return kvstore.Save(ctx, dir, "cache", name, r)
}

func LoadRelease(ctx context.Context, name string) (*github.RepositoryRelease, error) {
	dir, err := cacheFile(ctx)
	if err != nil {
		return nil, err
	}

	zerolog.Ctx(ctx).Debug().Str("name", name).Str("db", dir).Msg("loading release from cache")

	var r github.RepositoryRelease
	ok, err := kvstore.Load(ctx, dir, "cache", name, &r)
	if err != nil {
		return nil, err
	}

	if !ok {
		zerolog.Ctx(ctx).Warn().Str("name", name).Msg("cache miss")
		return nil, nil
	}
	return &r, err
}
