package kvstore

import (
	"context"
	"errors"
)

func Load[T any](ctx context.Context, database string, bucket string, name string, data *T) (bool, error) {

	store, closer, err := newStore(ctx, database)
	if err != nil {
		return false, err
	}

	defer closer()

	err = store.load(bucket, name, data)
	if err != nil {
		if err == ErrNotFound {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func Save[T any](ctx context.Context, database string, bucket string, name string, data *T) error {
	store, closer, err := newStore(ctx, database)
	if err != nil {
		return err
	}
	defer closer()

	if data == nil {
		return errors.New("nil token")
	}

	return store.save(bucket, name, data)
}

func LoadAll[T any](ctx context.Context, database string, bucket string, data map[string]*T) error {
	store, closer, err := newStore(ctx, database)
	if err != nil {
		return err
	}
	defer closer()

	if data == nil {
		return errors.New("nil token")
	}

	return store.loadAll(bucket, func(s string, a any) {
		data[s] = a.(*T)
	})
}
