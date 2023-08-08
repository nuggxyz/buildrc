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

func UploadDirAsTar(ctx context.Context, pipe Pipeline, fs afero.Fs, dir string, label string) error {
	// check if dir is empty
	osFiles, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("error reading directory %s: %v", dir, err)
	}

	if len(osFiles) == 0 {
		return fmt.Errorf("no files were output by script %s%s", dir, label)
	}

	fname := fmt.Sprintf("%s.tar.gz", label)

	fid := xid.New().String()

	// Create .tar.gz archive at pkg.OutputFile(arc).tar.gz
	tarCmd := exec.Command("tar", "-czvf", fid, "-C", dir, ".")
	tarCmd.Stdout = os.Stdout
	tarCmd.Stderr = os.Stderr
	err = tarCmd.Run()
	if err != nil {
		return fmt.Errorf("error creating .tar.gz archive: %v", err)
	}

	res, err := fs.Open(fid)
	if err != nil {
		return fmt.Errorf("error opening archive file %s: %v", fid, err)
	}

	err = pipe.UploadArtifact(ctx, fs, fname, res)
	if err != nil {
		return fmt.Errorf("error uploading archive file %s: %v", fname, err)
	}

	zerolog.Ctx(ctx).Debug().Str("dest_file", fid).Str("source_dir", dir).Msgf("created archive: %s", fname)

	return nil

}
