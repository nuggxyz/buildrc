package file

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"io"
	"path"
	"strings"

	"github.com/spf13/afero"
)

// Targz compresses the content of the given file.
// The file's read/write position will be reset to the beginning
// of the file before Targz returns, so the caller can continue to read from
// the file if needed.
func Targz(ctx context.Context, fls afero.Fs, pth string) (afero.File, error) {

	// name := filepath.Base(pth)
	var err error
	var writer *gzip.Writer
	var body []byte

	fle, err := fls.Open(pth)
	if err != nil {
		return nil, err
	}
	defer fle.Close()

	wrk, err := fls.Create(fle.Name() + ".tar.gz")
	if err != nil {
		return nil, err
	}

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

func Untargz(ctx context.Context, fls afero.Fs, pth string) (afero.File, error) {

	fle, err := fls.Open(pth)
	if err != nil {
		return nil, err
	}
	defer fle.Close()

	gr, err := gzip.NewReader(fle)
	if err != nil {
		return nil, err
	}
	defer gr.Close()

	tr := tar.NewReader(gr)

	// Assuming you want to extract to the same directory with the original name
	destPath := strings.TrimSuffix(fle.Name(), ".tar.gz")
	destFile, err := fls.Create(destPath)
	if err != nil {
		return nil, err
	}

	_, err = tr.Next()
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(destFile, tr)
	if err != nil {
		return nil, err
	}

	return destFile, nil
}
