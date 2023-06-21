package kvstore

import (
	"context"
	"errors"

	"github.com/rs/zerolog"
	"github.com/spf13/afero"
)

func Load[T any](ctx context.Context, fs afero.Fs, database string, name string, data *T) (bool, error) {

	store, closer, err := newStore(ctx, database, fs)
	if err != nil {
		return false, err
	}

	defer closer()

	err = store.load(name, data)
	if err != nil {
		if err == ErrNotFound {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func Save[T any](ctx context.Context, fs afero.Fs, database string, name string, data *T) error {
	store, closer, err := newStore(ctx, database, fs)
	if err != nil {
		return err
	}
	defer closer()

	if data == nil {
		zerolog.Ctx(ctx).Error().Str("name", name).Msg("nil data")
		return errors.New("nil token")
	}

	return store.save(name, data)
}

func LoadAll[T any](ctx context.Context, fs afero.Fs, database string, data map[string]T) error {
	store, closer, err := newStore(ctx, database, fs)
	if err != nil {
		return err
	}
	defer closer()

	if data == nil {
		zerolog.Ctx(ctx).Error().Msg("nil data")
		return errors.New("nil token")
	}

	return store.loadAll(func(s string, a any) {
		data[s] = a.(T)
	})
}
