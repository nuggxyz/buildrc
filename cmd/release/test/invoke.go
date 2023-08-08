package test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/nuggxyz/buildrc/cmd/release/setup"
	"github.com/nuggxyz/buildrc/internal/buildrc"
	"github.com/nuggxyz/buildrc/internal/common"
	"github.com/nuggxyz/buildrc/internal/pipeline"

	"github.com/rs/zerolog"
)

const (
	CommandID = "test"
	TestFile  = "test.sh"
)

type Handler struct {
}

func (me *Handler) Run(ctx context.Context, cmp common.Provider) (err error) {
	_, err = me.CachedTest(ctx, cmp)
	return err
}

func (me *Handler) CachedTest(ctx context.Context, prov common.Provider) (out *any, err error) {
	return pipeline.Cache(ctx, CommandID, prov, me.test)
}

func (me *Handler) test(ctx context.Context, prov common.Provider) (out *any, err error) {

	sv, err := setup.NewHandler("", "").Invoke(ctx, prov)
	if err != nil {
		return nil, err
	}

	// make sure the prebuild hook exists and is executable
	if _, err := os.Stat(TestFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("build hook %s does not exist", TestFile)
	}

	if err := os.Chmod(TestFile, 0755); err != nil {
		return nil, fmt.Errorf("error making build hook %s executable: %v", TestFile, err)
	}

	sha, err := prov.Git().GetCurrentCommitFromRef(ctx, "HEAD")
	if err != nil {
		return nil, err
	}

	err = me.run(ctx, TestFile, prov.Buildrc(), sv.Tag, sha, prov)
	if err != nil {
		return nil, err
	}

	return nil, nil

}

func (me *Handler) run(ctx context.Context, scriptPath string, brc *buildrc.Buildrc, tag string, commit string, prov common.Provider) error {
	ldflags, err := buildrc.GenerateGoLdflags(tag, commit)
	if err != nil {
		return err
	}
	return buildrc.RunAllPackages(ctx, brc, 10*time.Minute, func(ctx context.Context, pkg *buildrc.Package) error {

		dir, err := pipeline.NewTempDir(ctx, prov.Pipeline(), prov.FileSystem())
		if err != nil {
			return err
		}

		custom, err := pkg.CustomJSON()
		if err != nil {
			return fmt.Errorf("error marshalling custom JSON: %v", err)
		}

		cmd := exec.Command("bash", "./"+scriptPath, dir.String(), pkg.Name, custom)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = append(
			os.Environ(),
			fmt.Sprintf("GO_LDFLAGS=%s", ldflags),
			fmt.Sprintf("BUILDRC_VERSION=%s", tag),
			fmt.Sprintf("BUILDRC_COMMIT=%s", commit),
			fmt.Sprintf("BUILDRC_OUTPUT=%s", dir.String()),
			fmt.Sprintf("BUILDRC_CUSTOM=%s", custom),
			fmt.Sprintf("BUILDRC_NAME=%s", pkg.Name),
		)
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("error running test script %s for package %s: %v", scriptPath, pkg.Name, err)
		}

		zerolog.Ctx(ctx).Debug().Msgf("ran script %s for package %s", scriptPath, pkg.Name)

		// check if dir is empty
		osFiles, err := os.ReadDir(dir.String())
		if err != nil {
			return fmt.Errorf("error reading directory %s: %v", dir.String(), err)
		}

		if len(osFiles) == 0 {
			return fmt.Errorf("no files were output by script %s for package %s", scriptPath, pkg.Name)
		}

		fle := pipeline.GetNamedCacheFile(ctx, prov.Pipeline(), prov.FileSystem(), pkg.TestArchiveFileName())

		// Create .tar.gz archive at pkg.OutputFile(arc).tar.gz
		tarCmd := exec.Command("tar", "-czvf", pkg.TestArchiveFileName(), dir.String())
		tarCmd.Stdout = os.Stdout
		tarCmd.Stderr = os.Stderr
		err = tarCmd.Run()
		if err != nil {
			return fmt.Errorf("error creating .tar.gz archive: %v", err)
		}

		// move archive to cache
		err = os.Rename(pkg.TestArchiveFileName(), fle.String())
		if err != nil {
			return fmt.Errorf("error moving .tar.gz archive to cache: %v", err)
		}

		zerolog.Ctx(ctx).Debug().Str("dest_file", fle.String()).Str("source_dir", dir.String()).Msgf("created archive for package %s", pkg.Name)

		return nil
	})

}
