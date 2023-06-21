package pipeline

import (
	"context"
	"path/filepath"

	"github.com/nuggxyz/buildrc/internal/kvstore"
	"github.com/rs/zerolog"
	"github.com/spf13/afero"
)

type TempFile string

func (me TempFile) String() string {
	return string(me)
}

func tempFile(ctx context.Context, pipe Pipeline, fs afero.Fs, name string) TempFile {
	dir := name + ".temp.db"
	if envvar, err := BuildrcTempDir.Load(ctx, pipe, fs); err == nil && envvar != "" {
		dir = filepath.Join(envvar, dir)
	}

	return TempFile(dir)
}

func cacheHitFile(ctx context.Context, pipe Pipeline, fs afero.Fs) TempFile {
	return tempFile(ctx, pipe, fs, "cache-hit")
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
