package cache

import (
	"context"

	"github.com/google/go-github/v53/github"
	"github.com/nuggxyz/buildrc/internal/kvstore"
)

const (
	CACHE_DIR = ".buildrc-cache"
)

func SaveRelease(ctx context.Context, name string, r *github.RepositoryRelease) error {
	return kvstore.Save(ctx, CACHE_DIR+"/cache.db", "cache", name, r)
}

func LoadRelease(ctx context.Context, name string) (*github.RepositoryRelease, error) {
	var r github.RepositoryRelease
	_, err := kvstore.Load(ctx, CACHE_DIR+"/cache.db", "cache", name, &r)
	return &r, err
}
