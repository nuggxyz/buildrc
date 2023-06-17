package provider

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nuggxyz/buildrc/internal/buildrc"
	"github.com/rs/zerolog"
)

func ProviderTempFileName(cmd string) (string, error) {
	r, err := buildrc.BuildrcCacheDir.Load()
	if err != nil {
		return "", err
	}
	return filepath.Join(r, fmt.Sprintf("%s.provider-content.json", cmd)), nil
}

func Load(ctx context.Context, prov ContentProvider, cmd string) ([]byte, error) {
	tmp, err := ProviderTempFileName(cmd)
	if err != nil {
		return nil, err
	}
	// try to load from tmp folder
	f, err := prov.FileSystem().Get(ctx, tmp)
	if err != nil {

		// if not found do nothing
		if errors.Is(err, os.ErrNotExist) {
			return []byte{}, nil
		}

		return nil, err
	}

	zerolog.Ctx(ctx).Debug().Str("data", string(f)).Msgf("loaded result from %s", tmp)

	return f, nil
}

func Save(ctx context.Context, prov ContentProvider, cmd string, result []byte) error {

	tmp, err := ProviderTempFileName(cmd)
	if err != nil {
		return err
	}

	// save to tmp folder
	err = prov.FileSystem().Put(ctx, tmp, result)
	if err != nil {
		return err
	}

	return nil
}
