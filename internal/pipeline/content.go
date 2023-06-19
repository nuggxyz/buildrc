package pipeline

import (
	"context"
	"os"

	"github.com/nuggxyz/buildrc/internal/file"
)

type Pipeline interface {
	AddToEnv(context.Context, string, string) error
	FileSystem() file.FileAPI
}

type DefaultPipeline struct {
	fs file.FileAPI
}

func NewDefaultPipeline(fs file.FileAPI) Pipeline {
	return &DefaultPipeline{fs}
}

func (me *DefaultPipeline) FileSystem() file.FileAPI {
	return me.fs
}

func (me *DefaultPipeline) AddToEnv(ctx context.Context, key, value string) error {
	return os.Setenv(key, value)
}
