package pipeline

type UploadDirAsTarOpts struct {
	ProduceSHA256 bool
	ProduceRaw    bool
	ProduceTarGz  bool
}

// func UploadDirAsTar(ctx context.Context, pipe Pipeline, fs afero.Fs, dir string, label string, opts *UploadDirAsTarOpts) error {

// 	if opts == nil {
// 		opts = &UploadDirAsTarOpts{
// 			RequireFiles:  false,
// 			ProduceSHA256: false,
// 		}
// 	}

// 	// check if dir is empty
// 	osFiles, err := fs.Open(dir)
// 	if err != nil {
// 		return fmt.Errorf("error reading directory %s: %v", dir, err)
// 	}

// 	osFilesInfo, err := osFiles.Readdir(0)
// 	if err != nil {
// 		return fmt.Errorf("error reading directory %s: %v", dir, err)
// 	}

// 	if opts.RequireFiles && len(osFilesInfo) == 0 {
// 		return fmt.Errorf("no files were output by script %s%s", dir, label)
// 	}

// 	if opts.ProduceRaw {
// 		res, err := fs.Open(
// 		if err != nil {
// 			return fmt.Errorf("error opening archive file %s: %v", fid, err)
// 		}
// 		err = pipe.UploadArtifact(ctx, fs, fname, res)
// 		if err != nil {
// 			zerolog.Ctx(ctx).Error().Msgf("error uploading archive file %s: %v", fname, err)
// 			return fmt.Errorf("error uploading archive file %s: %v", fname, err)
// 		}
// 	}

// 	if opts.ProduceTarGz {

// 		fname := fmt.Sprintf("%s.tar.gz", label)

// 		fid := xid.New().String()

// 		zerolog.Ctx(ctx).Debug().Msgf("creating archive %s", fid)

// 		// Create .tar.gz archive at pkg.OutputFile(arc).tar.gz
// 		tarCmd := exec.Command("tar", "-czvf", fid, "-C", dir, ".")
// 		tarCmd.Stdout = os.Stdout
// 		tarCmd.Stderr = os.Stderr
// 		err = tarCmd.Run()
// 		if err != nil {
// 			return fmt.Errorf("error creating .tar.gz archive: %v", err)
// 		}

// 		zerolog.Ctx(ctx).Debug().Msgf("created archive %s", fid)

// 		res, err := fs.Open(fid)
// 		if err != nil {
// 			return fmt.Errorf("error opening archive file %s: %v", fid, err)
// 		}

// 		err = pipe.UploadArtifact(ctx, fs, fname, res)
// 		if err != nil {
// 			zerolog.Ctx(ctx).Error().Msgf("error uploading archive file %s: %v", fname, err)
// 			return fmt.Errorf("error uploading archive file %s: %v", fname, err)
// 		}

// 		zerolog.Ctx(ctx).Debug().Msgf("uploaded archive %s", fname)
// 	}

// 	if opts.ProduceSHA256 {
// 		hname := fmt.Sprintf("%s.sha256", label)

// 		// Compute and write SHA-256 checksum to pkg.OutputFile(arc).sha256
// 		hashCmd := exec.Command("shasum", "-a", "256", fid)
// 		hashOutput, err := hashCmd.Output()
// 		if err != nil {
// 			return fmt.Errorf("error computing SHA-256 checksum: %v", err)
// 		}

// 		zerolog.Ctx(ctx).Debug().Msgf("computed SHA-256 checksum for %s", hname)

// 		fle, err := fs.Create(hname)
// 		if err != nil {
// 			return fmt.Errorf("error creating SHA-256 checksum file %s: %v", hname, err)
// 		}

// 		defer fle.Close()

// 		_, err = fle.Write(hashOutput)
// 		if err != nil {
// 			return fmt.Errorf("error writing SHA-256 checksum file %s: %v", hname, err)
// 		}

// 		err = pipe.UploadArtifact(ctx, fs, hname, fle)
// 		if err != nil {
// 			return fmt.Errorf("error uploading archive file %s: %v", hname, err)
// 		}

// 		zerolog.Ctx(ctx).Debug().Msgf("wrote SHA-256 checksum to %s", hname)

// 	}

// 	zerolog.Ctx(ctx).Debug().Str("dest_file", fid).Str("source_dir", dir).Msgf("created archive: %s", fname)

// 	return nil

// }
