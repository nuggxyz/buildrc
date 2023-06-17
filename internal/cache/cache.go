package cache

import (
	"context"
	"path/filepath"

	"github.com/google/go-github/v53/github"
	"github.com/nuggxyz/buildrc/internal/buildrc"
	"github.com/nuggxyz/buildrc/internal/kvstore"
	"github.com/rs/zerolog"
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
	if envvar, err := buildrc.BuildrcCacheDir.Load(); err == nil && envvar != "" {
		dir = envvar
	} else {
		return "", err
	}
	return filepath.Join(dir, "cache.db"), nil
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

func SaveEnvVar(ctx context.Context, name string, value string) error {
	dir, err := cacheFile(ctx)
	if err != nil {
		return err
	}

	zerolog.Ctx(ctx).Debug().Str("name", name).Str("db", dir).Msg("saving env var to cache")

	return kvstore.Save(ctx, dir, "env", name, &value)
}

func LoadAllEnvVars(ctx context.Context) (map[string]string, bool, error) {

	ok, err := HasCacheBeenHit(ctx, "load-all-env-vars")
	if err != nil {
		return nil, false, err
	}

	dir, err := cacheFile(ctx)
	if err != nil {
		return nil, false, err
	}

	zerolog.Ctx(ctx).Debug().Str("db", dir).Msg("loading all env vars from cache")

	vars := map[string]string{}
	err = kvstore.LoadAll(ctx, dir, "env", vars)
	if err != nil {
		if kvstore.IsNotFound(err) {
			return nil, false, nil
		}
		return nil, false, err
	}

	err = RecordCacheHit(ctx, "load-all-env-vars")
	if err != nil {
		return nil, false, err
	}

	return vars, ok, err
}
