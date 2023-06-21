package pipeline

import (
	"context"

	"github.com/spf13/afero"
)

type Pipeline interface {
	AddToEnv(context.Context, string, string, afero.Fs) error
	GetFromEnv(context.Context, string, afero.Fs) (string, error)
}

type MemoryPipeline struct {
	env map[string]string
}

func NewMemoryPipeline() Pipeline {
	return &MemoryPipeline{
		env: map[string]string{},
	}
}

func (me *MemoryPipeline) AddToEnv(ctx context.Context, key, value string, _ afero.Fs) error {
	me.env[key] = value
	return nil
}

func (me *MemoryPipeline) GetFromEnv(ctx context.Context, key string, _ afero.Fs) (string, error) {
	return me.env[key], nil
}
