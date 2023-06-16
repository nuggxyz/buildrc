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

func SaveRelease(ctx context.Context, name string, r *github.RepositoryRelease) error {

	dir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	return kvstore.Save(ctx, filepath.Join(dir, CACHE_DIR, "cache.db"), "cache", name, r)
}

func LoadRelease(ctx context.Context, name string) (*github.RepositoryRelease, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	var r github.RepositoryRelease
	ok, err := kvstore.Load(ctx, filepath.Join(dir, CACHE_DIR, "cache.db"), "cache", name, &r)
	if err != nil {
		return nil, err
	}

	if !ok {
		zerolog.Ctx(ctx).Warn().Str("name", name).Msg("cache miss")
		return nil, nil
	}
	return &r, err
}
