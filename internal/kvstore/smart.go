package kvstore

import (
	"context"
	"errors"

	"github.com/rs/zerolog"
	"github.com/spf13/afero"
)

func Load[T any](ctx context.Context, fs afero.Fs, database string, name string, data *T) (bool, error) {

	store := NewStore(ctx, fs, database)

	err := store.Load(name, data)
	if err != nil {
		if err == ErrNotFound {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func Save[T any](ctx context.Context, fs afero.Fs, database string, name string, data *T) error {
	store := NewStore(ctx, fs, database)

	if data == nil {
		zerolog.Ctx(ctx).Error().Str("name", name).Msg("nil data")
		return errors.New("nil token")
	}

	return store.Save(name, data)
}

func LoadAll[T any](ctx context.Context, fs afero.Fs, database string, data map[string]T) error {
	store := NewStore(ctx, fs, database)

	if data == nil {
		zerolog.Ctx(ctx).Error().Msg("nil data")
		return errors.New("nil token")
	}

	return store.LoadAll(func(s string, a any) {
		data[s] = a.(T)
	})
}
