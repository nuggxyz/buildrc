package kvstore

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/spf13/afero"

	"github.com/cockroachdb/pebble"
	"github.com/cockroachdb/pebble/vfs"
)

type Store struct {
	vfs vfs.FS
	db  *pebble.DB
}

func NewFSFromAfero(afs afero.Fs) vfs.FS {
	if _, ok := afs.(*afero.MemMapFs); ok {
		return vfs.NewMem()
	} else {
		return vfs.Default
	}
}

func EnsureStoreFile(ctx context.Context, dir string, f afero.Fs) error {
	_, closer, err := newStore(ctx, dir, f)
	if err != nil {
		return err
	}

	defer closer()

	return nil
}

func newStore(ctx context.Context, dir string, f afero.Fs) (*Store, func() error, error) {

	v := NewFSFromAfero(f)

	r, err := pebble.Open(dir, &pebble.Options{FS: v})
	if err != nil {
		return nil, nil, err
	}

	return &Store{db: r, vfs: v}, r.Close, nil
}

var ErrNotFound = errors.New("key not found")

func IsNotFound(err error) bool {
	return errors.Is(err, pebble.ErrNotFound) || errors.Is(err, ErrNotFound)
}

func (s *Store) save(key string, data any) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return s.db.Set([]byte(key), bytes, pebble.Sync)
}

func (s *Store) load(key string, data any) error {
	bytes, closer, err := s.db.Get([]byte(key))
	if err != nil {
		if errors.Is(err, pebble.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}
	defer closer.Close()
	return json.Unmarshal(bytes, data)
}

func (s *Store) listKeys() ([]string, error) {
	var keys []string
	iter := s.db.NewIter(&pebble.IterOptions{})
	defer iter.Close()
	for iter.First(); iter.Valid(); iter.Next() {
		key := string(iter.Key())
		keys = append(keys, key)
	}
	if err := iter.Error(); err != nil {
		return nil, err
	}
	return keys, nil
}

func (s *Store) loadAll(cb func(string, any)) error {
	iter := s.db.NewIter(&pebble.IterOptions{})
	defer iter.Close()
	for iter.First(); iter.Valid(); iter.Next() {
		key := string(iter.Key())
		var data any
		v, err := iter.ValueAndErr()
		if err != nil {
			return err
		}
		if err := json.Unmarshal(v, &data); err != nil {
			return err
		}
		cb(key, data)
	}
	if err := iter.Error(); err != nil {
		return err
	}
	return nil
}
