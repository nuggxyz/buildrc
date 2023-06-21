package fs

import (
	"io/fs"

	"github.com/cockroachdb/pebble/vfs"
	"github.com/spf13/afero"
)

type FS interface {
	Afero() afero.Fs
	OS() fs.FS
	Pebble() vfs.FS
}

func NewMemoryFS() FS {
	return &baseFS{
		afero:  afero.NewMemMapFs(),
		pebble: vfs.NewMem(),
	}
}

func NewOSFS() FS {
	return &baseFS{
		afero:  afero.NewOsFs(),
		pebble: vfs.Default,
	}
}

type baseFS struct {
	afero  afero.Fs
	pebble vfs.FS
}

func (m *baseFS) Afero() afero.Fs {
	return m.afero
}

func (m *baseFS) Pebble() vfs.FS {
	return m.pebble
}

func (m *baseFS) OS() fs.FS {
	return afero.NewIOFS(m.afero)
}
