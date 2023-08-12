package file

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"io"
	"path"
	"path/filepath"

	"github.com/spf13/afero"
)

func Targz(ctx context.Context, fs afero.Fs, pth string) (afero.File, error) {
	name := filepath.Base(pth)
	var err error
	var writer *gzip.Writer
	var body []byte

	fle, err := fs.Open(pth)
	if err != nil {
		return nil, err
	}
	defer fle.Close()

	wrk, err := afero.TempFile(fs, "", filepath.Join(name, ".tar.gz"))
	if err != nil {
		return nil, err
	}
	defer wrk.Close()

	wrkstats, err := wrk.Stat()
	if err != nil {
		return nil, err
	}

	if writer, err = gzip.NewWriterLevel(wrk, gzip.BestCompression); err != nil {
		return nil, err
	}
	defer writer.Close()

	tw := tar.NewWriter(writer)
	defer tw.Close()

	if body, err = io.ReadAll(fle); err != nil {
		return nil, err
	}

	if body != nil {
		hdr := &tar.Header{
			Name: path.Base(wrkstats.Name()),
			Mode: int64(0644),
			Size: int64(len(body)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return nil, err
		}
		if _, err := tw.Write(body); err != nil {
			return nil, err
		}
	}

	return wrk, nil
}
