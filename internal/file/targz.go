package file

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"io"
	"path/filepath"
	"strings"

	"github.com/nuggxyz/buildrc/internal/logging"
	"github.com/spf13/afero"
)

func Targz(ctx context.Context, fs afero.Fs, pth string) (afero.File, error) {
	wrk, err := fs.Create(pth + ".tar.gz")
	if err != nil {
		return nil, logging.WrapError(ctx, err)
	}

	writer, err := gzip.NewWriterLevel(wrk, gzip.BestCompression)
	if err != nil {
		return nil, logging.WrapError(ctx, err)
	}
	defer writer.Close()

	tw := tar.NewWriter(writer)
	defer tw.Close()

	if err := addFilesToTar(ctx, fs, tw, pth, ""); err != nil {
		return nil, logging.WrapError(ctx, err)
	}

	return wrk, nil
}

func addFilesToTar(ctx context.Context, fs afero.Fs, tw *tar.Writer, pth string, prefix string) error {

	stats, err := fs.Stat(pth)
	if err != nil {
		return logging.WrapError(ctx, err)
	}

	if stats.IsDir() {
		infos, err := afero.ReadDir(fs, pth)
		if err != nil {
			return logging.WrapError(ctx, err)
		}
		for _, info := range infos {
			newPath := filepath.Join(pth, info.Name())
			newPrefix := filepath.Join(prefix, info.Name())
			if err := addFilesToTar(ctx, fs, tw, newPath, newPrefix); err != nil {
				return logging.WrapError(ctx, err)
			}
		}
		return nil
	}

	file, err := fs.Open(pth)
	if err != nil {
		return logging.WrapError(ctx, err)
	}
	defer file.Close()

	body, err := io.ReadAll(file)
	if err != nil {
		return logging.WrapError(ctx, err)
	}

	hdr := &tar.Header{
		Name: prefix,
		Mode: int64(0644),
		Size: int64(len(body)),
	}

	if err := tw.WriteHeader(hdr); err != nil {
		return logging.WrapError(ctx, err)
	}

	if _, err := tw.Write(body); err != nil {
		return logging.WrapError(ctx, err)
	}
	// }

	return nil
}

// Targz compresses the content of the given file.
// The file's read/write position will be reset to the beginning
// of the file before Targz returns, so the caller can continue to read from
// the file if needed.
// func Targz(ctx context.Context, fls afero.Fs, pth string) (afero.File, error) {

// 	// name := filepath.Base(pth)
// 	var err error
// 	var writer *gzip.Writer
// 	var body []byte

// 	stat, err := fls.Stat(pth)
// 	if err != nil {
// 		return nil, logging.WrapError(ctx, err,)
// 	}

// 	fle, err := fls.Open(pth)
// 	if err != nil {
// 		return nil, logging.WrapError(ctx, err,)
// 	}
// 	defer fle.Close()

// 	wrk, err := fls.Create(fle.Name() + ".tar.gz")
// 	if err != nil {
// 		return nil, logging.WrapError(ctx, err,)
// 	}

// 	if writer, err = gzip.NewWriterLevel(wrk, gzip.BestCompression); err != nil {
// 		return nil, logging.WrapError(ctx, err,)
// 	}
// 	defer writer.Close()

// 	tw := tar.NewWriter(writer)
// 	defer tw.Close()

// 	files := []afero.File{}

// 	if stat.IsDir() {
// 		fileinfos, err := afero.ReadDir(fls, pth)
// 		if err != nil {
// 			return nil, logging.WrapError(ctx, err,)
// 		}

// 		for _, fileinfo := range fileinfos {
// 			if fileinfo.IsDir() {
// 				continue
// 			}
// 			file, err := fls.Open(path.Join(pth, fileinfo.Name()))
// 			if err != nil {
// 				return nil, logging.WrapError(ctx, err,)
// 			}
// 			files = append(files, file)
// 		}
// 	} else {
// 		files = append(files, fle)
// 	}

// 	for _, fle := range files {

// 		if body, err = io.ReadAll(fle); err != nil {
// 			return nil, logging.WrapError(ctx, err,)
// 		}

// 		if body != nil {
// 			hdr := &tar.Header{
// 				Name: path.Base(fle.Name()),
// 				Mode: int64(0644),
// 				Size: int64(len(body)),
// 			}
// 			if err := tw.WriteHeader(hdr); err != nil {
// 				return nil, logging.WrapError(ctx, err,)
// 			}
// 			if _, err := tw.Write(body); err != nil {
// 				return nil, logging.WrapError(ctx, err,)
// 			}
// 		}
// 	}

// 	return wrk, nil
// }

func Untargz(ctx context.Context, fs afero.Fs, pth string) (afero.File, error) {
	fle, err := fs.Open(pth)
	if err != nil {
		return nil, logging.WrapError(ctx, err)
	}
	defer fle.Close()

	dest := strings.TrimSuffix(fle.Name(), ".tar.gz")

	gr, err := gzip.NewReader(fle)
	if err != nil {
		return nil, logging.WrapError(ctx, err)
	}
	defer gr.Close()

	tr := tar.NewReader(gr)

	first := true

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, logging.WrapError(ctx, err)
		}

		destPath := filepath.Join(dest, hdr.Name) // Update the destination directory as needed
		if hdr.Typeflag == tar.TypeDir {
			if err := fs.MkdirAll(destPath, 0755); err != nil {
				return nil, logging.WrapError(ctx, err)
			}
			continue
		} else if first {
			// if it is the first and only file, we want to extract it to the same directory with the original name
			destPath = dest
		}

		first = false

		destFile, err := fs.Create(destPath)
		if err != nil {
			return nil, logging.WrapError(ctx, err)
		}

		_, err = io.Copy(destFile, tr)
		if err != nil {
			return nil, logging.WrapError(ctx, err)
		}

		if err := destFile.Close(); err != nil {
			return nil, logging.WrapError(ctx, err)
		}
	}

	return fle, nil
}

// func Untargz(ctx context.Context, fls afero.Fs, pth string) (afero.File, error) {

// 	fle, err := fls.Open(pth)
// 	if err != nil {
// 		return nil, logging.WrapError(ctx, err,)
// 	}
// 	defer fle.Close()

// 	gr, err := gzip.NewReader(fle)
// 	if err != nil {
// 		return nil, logging.WrapError(ctx, err,)
// 	}
// 	defer gr.Close()

// 	tr := tar.NewReader(gr)

// 	// // Assuming you want to extract to the same directory with the original name
// 	// destPath := strings.TrimSuffix(fle.Name(), ".tar.gz")
// 	// destFile, err := fls.Create(destPath)
// 	// if err != nil {
// 	// 	return nil, logging.WrapError(ctx, err,)
// 	// }

// 	// Iterate through the files in the tar archive
// 	for {
// 		hdr, err := tr.Next()
// 		if err == io.EOF {
// 			break // Reached end of archive
// 		}
// 		if err != nil {
// 			return nil, logging.WrapError(ctx, err,)
// 		}

// 		// Create destination file based on header name
// 		destPath := filepath.Join("destination_directory", hdr.Name) // Update the destination directory as needed
// 		destFile, err := fls.Create(destPath)
// 		if err != nil {
// 			return nil, logging.WrapError(ctx, err,)
// 		}

// 		// Copy content to destination file
// 		_, err = io.Copy(destFile, tr)
// 		if err != nil {
// 			return nil, logging.WrapError(ctx, err,)
// 		}

// 		// Close destination file
// 		if err := destFile.Close(); err != nil {
// 			return nil, logging.WrapError(ctx, err,)
// 		}
// 	}

// 	return destFile, nil
// }
