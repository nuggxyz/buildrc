package cache

import (
	"context"
	"os"
	"path/filepath"

	"github.com/google/go-github/v53/github"
	"github.com/nuggxyz/buildrc/internal/kvstore"
	"github.com/rs/zerolog"
)

const (
	CACHE_DIR = ".buildrc-cache"
)

func cacheFile(ctx context.Context) (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, CACHE_DIR, "cache.db"), nil
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
