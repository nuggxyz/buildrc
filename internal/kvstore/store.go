package kvstore

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/peterbourgon/diskv"
	"github.com/spf13/afero"
)

type Store struct {
	dv *diskv.Diskv
}

func NewStore(ctx context.Context, fs afero.Fs, basePath string) *Store {
	dv := diskv.New(diskv.Options{
		BasePath:     basePath,
		CacheSizeMax: 1024 * 1024, // 1MB
	})
	return &Store{dv: dv}
}

var ErrNotFound = errors.New("key not found")

func (s *Store) Save(key string, data any) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return s.dv.Write(key, bytes)
}

func (s *Store) Load(key string, data any) error {
	bytes, err := s.dv.Read(key)
	if err != nil {
		return ErrNotFound
	}
	return json.Unmarshal(bytes, data)
}

func (s *Store) ListKeys() ([]string, error) {
	strs := []string{}
	for key := range s.dv.Keys(nil) {
		strs = append(strs, key)
	}
	return strs, nil
}

func (s *Store) LoadAll(cb func(string, any)) error {
	for key := range s.dv.Keys(nil) {
		var data any
		v, err := s.dv.Read(key)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(v, &data); err != nil {
			return err
		}
		cb(key, data)
	}
	return nil
}

func (s *Store) Delete(key string) error {
	return s.dv.Erase(key)
}

func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}
