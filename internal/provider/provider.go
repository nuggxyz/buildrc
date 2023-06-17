package provider

import (
	"context"
	"os"

	"github.com/nuggxyz/buildrc/internal/file"
)

type ContentProvider interface {
	AddToEnv(context.Context, string, string) error
	FileSystem() file.FileAPI
}

type DefaultContentProvider struct {
	fs file.FileAPI
}

func NewDefaultContentProvider(fs file.FileAPI) ContentProvider {
	return &DefaultContentProvider{fs}
}

func (me *DefaultContentProvider) FileSystem() file.FileAPI {
	return me.fs
}

func (me *DefaultContentProvider) AddToEnv(ctx context.Context, key, value string) error {
	return os.Setenv(key, value)
}
