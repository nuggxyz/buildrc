package pipeline

import (
	"context"
	"path/filepath"

	"github.com/nuggxyz/buildrc/internal/kvstore"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"github.com/spf13/afero"
)

type TempFile string

func (me TempFile) String() string {
	return string(me)
}

func tempDBFile(ctx context.Context, pipe Pipeline, fs afero.Fs, name string) TempFile {
	dir := name + ".temp.db"
	if envvar, err := BuildrcTempDir.Load(ctx, pipe, fs); err == nil && envvar != "" {
		dir = filepath.Join(envvar, dir)
	}

	return TempFile(dir)
}

func cacheHitFile(ctx context.Context, pipe Pipeline, fs afero.Fs) TempFile {
	return tempDBFile(ctx, pipe, fs, "cache-hit")
}

func NewTempDir(ctx context.Context, pipe Pipeline, fs afero.Fs) (TempFile, error) {
	dir := xid.New().String()

	if envvar, err := BuildrcTempDir.Load(ctx, pipe, fs); err == nil && envvar != "" {
		dir = filepath.Join(envvar, dir)
	}

	if err := fs.MkdirAll(dir, 0755); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("dir", dir).Msg("failed to create temp dir")
		return "", err
	}

	return TempFile(dir), nil
}

func NewNamedTempDir(ctx context.Context, pipe Pipeline, fs afero.Fs, name string) (TempFile, error) {

	if envvar, err := BuildrcTempDir.Load(ctx, pipe, fs); err == nil && envvar != "" {
		name = filepath.Join(envvar, name)
	}

	if err := fs.MkdirAll(name, 0755); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("dir", name).Msg("failed to create temp dir")
		return "", err
	}

	return TempFile(name), nil
}

func NewTempFile(ctx context.Context, pipe Pipeline, fs afero.Fs) (TempFile, error) {
	fle := xid.New().String()

	drr, err := NewTempDir(ctx, pipe, fs)
	if err != nil {
		return "", err
	}

	fle = filepath.Join(drr.String(), fle)

	res, err := fs.Create(fle)
	if err != nil {
		return "", err
	}

	if err := res.Close(); err != nil {
		return "", err
	}

	return TempFile(fle), nil
}

func HasCacheBeenHit(ctx context.Context, p Pipeline, fs afero.Fs, flag string) (bool, error) {
	dir := cacheHitFile(ctx, p, fs)

	zerolog.Ctx(ctx).Debug().Str("db", dir.String()).Msg("checking if cache has been hit")
	res := false
	l, err := kvstore.Load(ctx, fs, dir.String(), flag, &res)
	if err != nil {
		if kvstore.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}

	return l, nil
}

func RecordCacheHit(ctx context.Context, p Pipeline, fs afero.Fs, flag string) error {
	dir := cacheHitFile(ctx, p, fs)

	zerolog.Ctx(ctx).Debug().Str("db", dir.String()).Msg("recording cache hit")

	dat := true

	return kvstore.Save(ctx, fs, dir.String(), flag, &dat)
}
