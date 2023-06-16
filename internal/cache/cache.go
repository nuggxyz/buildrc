package cache

import (
	"context"
	"os"
	"path/filepath"

	"github.com/google/go-github/v53/github"
	"github.com/nuggxyz/buildrc/internal/kvstore"
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
	_, err = kvstore.Load(ctx, filepath.Join(dir, CACHE_DIR, "cache.db"), "cache", name, &r)
	return &r, err
}
