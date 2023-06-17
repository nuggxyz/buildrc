package kvstore

import (
	"context"
	"encoding/json"
	"errors"

	"go.etcd.io/bbolt"
)

type Store struct {
	db *bbolt.DB
}

func EnsureStoreFile(ctx context.Context, dir string) error {
	_, closer, err := newStore(ctx, dir)
	if err != nil {
		return err
	}

	defer closer()

	return nil
}

func newStore(ctx context.Context, dir string) (*Store, func() error, error) {

	db, err := bbolt.Open(dir, 0600, nil)
	if err != nil {
		return nil, nil, err
	}

	return &Store{db: db}, func() error {
		return db.Close()
	}, nil
}

func (s *Store) save(bucket, key string, data any) error {

	return s.db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		bytes, err := json.Marshal(data)
		if err != nil {
			return err
		}
		return b.Put([]byte(key), bytes)
	})
}

var ErrNotFound = errors.New("key not found")

func (s *Store) load(bucket, key string, data any) error {

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return ErrNotFound
		}
		bytes := b.Get([]byte(key))
		if bytes == nil {
			return ErrNotFound
		}
		return json.Unmarshal(bytes, &data)
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) listKeys(bucket string) ([]string, error) {

	keys := make([]string, 0)

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return ErrNotFound
		}
		return b.ForEach(func(k, v []byte) error {
			keys = append(keys, string(k))
			return nil
		})
	})

	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (s *Store) loadAll(bucket string, cb func(string, any)) error {

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return ErrNotFound
		}
		return b.ForEach(func(k, v []byte) error {
			cb(string(k), v)
			return nil
		})
	})

	if err != nil {
		return err
	}

	return nil
}
