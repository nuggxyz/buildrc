package pipeline

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"github.com/spf13/afero"
)

type UploadDirAsTarOpts struct {
	RequireFiles  bool
	ProduceSHA256 bool
}

func UploadDirAsTar(ctx context.Context, pipe Pipeline, fs afero.Fs, dir string, label string, opts *UploadDirAsTarOpts) error {

	if opts == nil {
		opts = &UploadDirAsTarOpts{
			RequireFiles:  false,
			ProduceSHA256: false,
		}
	}

	// check if dir is empty
	osFiles, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("error reading directory %s: %v", dir, err)
	}

	if opts.RequireFiles && len(osFiles) == 0 {
		return fmt.Errorf("no files were output by script %s%s", dir, label)
	}

	fname := fmt.Sprintf("%s.tar.gz", label)

	fid := xid.New().String()

	zerolog.Ctx(ctx).Debug().Msgf("creating archive %s", fid)

	// Create .tar.gz archive at pkg.OutputFile(arc).tar.gz
	tarCmd := exec.Command("tar", "-czvf", fid, "-C", dir, ".")
	tarCmd.Stdout = os.Stdout
	tarCmd.Stderr = os.Stderr
	err = tarCmd.Run()
	if err != nil {
		return fmt.Errorf("error creating .tar.gz archive: %v", err)
	}

	zerolog.Ctx(ctx).Debug().Msgf("created archive %s", fid)

	res, err := fs.Open(fid)
	if err != nil {
		return fmt.Errorf("error opening archive file %s: %v", fid, err)
	}

	err = pipe.UploadArtifact(ctx, fs, fname, res)
	if err != nil {
		zerolog.Ctx(ctx).Error().Msgf("error uploading archive file %s: %v", fname, err)
		return fmt.Errorf("error uploading archive file %s: %v", fname, err)
	}

	zerolog.Ctx(ctx).Debug().Msgf("uploaded archive %s", fname)

	if opts.ProduceSHA256 {
		hname := fmt.Sprintf("%s.sha256", label)

		// Compute and write SHA-256 checksum to pkg.OutputFile(arc).sha256
		hashCmd := exec.Command("shasum", "-a", "256", fid)
		hashOutput, err := hashCmd.Output()
		if err != nil {
			return fmt.Errorf("error computing SHA-256 checksum: %v", err)
		}

		zerolog.Ctx(ctx).Debug().Msgf("computed SHA-256 checksum for %s", hname)

		fle, err := fs.Create(hname)
		if err != nil {
			return fmt.Errorf("error creating SHA-256 checksum file %s: %v", hname, err)
		}

		_, err = fle.Write(hashOutput)
		if err != nil {
			defer fle.Close()
			return fmt.Errorf("error writing SHA-256 checksum file %s: %v", hname, err)
		}

		err = fle.Close()
		if err != nil {
			return fmt.Errorf("error closing SHA-256 checksum file %s: %v", hname, err)
		}

		err = pipe.UploadArtifact(ctx, fs, hname, fle)
		if err != nil {
			return fmt.Errorf("error uploading archive file %s: %v", hname, err)
		}

		zerolog.Ctx(ctx).Debug().Msgf("wrote SHA-256 checksum to %s", hname)

	}

	zerolog.Ctx(ctx).Debug().Str("dest_file", fid).Str("source_dir", dir).Msgf("created archive: %s", fname)

	return nil

}
