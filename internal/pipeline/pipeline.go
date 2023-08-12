package pipeline

import (
	"context"

	"github.com/spf13/afero"
)

type Pipeline interface {
	AddToEnv(context.Context, string, string, afero.Fs) error
	GetFromEnv(context.Context, string, afero.Fs) (string, error)

	// RunId(ctx context.Context) (int64, error)
	UploadArtifact(ctx context.Context, fls afero.Fs, name string, fle afero.File) error

	DownloadArtifact(context.Context, afero.Fs, string) (afero.File, error)
	SupportsDocker() bool
}

var _ Pipeline = &MemoryPipeline{}

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

func (me *MemoryPipeline) UploadArtifact(ctx context.Context, _ afero.Fs, name string, _ afero.File) error {
	return nil
}

func (me *MemoryPipeline) DownloadArtifact(ctx context.Context, _ afero.Fs, name string) (afero.File, error) {
	return nil, nil
}

func (me *MemoryPipeline) SupportsDocker() bool {
	return true
}
