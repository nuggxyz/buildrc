package build

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/nuggxyz/buildrc/cmd/release/setup"
	"github.com/nuggxyz/buildrc/internal/buildrc"
	"github.com/nuggxyz/buildrc/internal/common"
	"github.com/nuggxyz/buildrc/internal/git"
	"github.com/nuggxyz/buildrc/internal/pipeline"

	"github.com/rs/zerolog"
)

const (
	CommandID = "build"
	BuildFile = "build.sh"
)

type Handler struct {
	Package string `arg:"name" help:"The name of the package to load."`
}

func (me *Handler) Run(ctx context.Context, cmp common.Provider) (err error) {
	_, err = me.CachedBuild(ctx, cmp)
	return err
}

func (me *Handler) CachedBuild(ctx context.Context, prov common.Provider) (out *any, err error) {
	return pipeline.Cache(ctx, "build", prov, me.build)
}

func (me *Handler) build(ctx context.Context, prov common.Provider) (out *any, err error) {

	sv, err := setup.NewHandler("", "").Invoke(ctx, prov)
	if err != nil {
		return nil, err
	}

	zerolog.Ctx(ctx).Info().Msg("checking if build is required")

	ok, tagg, err := git.ReleaseAlreadyExists(ctx, prov.Release(), prov.Git())
	if err != nil {
		return nil, err
	}

	if ok {
		zerolog.Ctx(ctx).Info().Bool("release_aleady_exists", ok).Str("tag", tagg).Msg("build not required")
		return nil, nil
	} else {
		zerolog.Ctx(ctx).Info().Msg("build required, continuing")
	}

	// make sure the prebuild hook exists and is executable
	if _, err := os.Stat(BuildFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("build hook %s does not exist", BuildFile)
	}

	if err := os.Chmod(BuildFile, 0755); err != nil {
		return nil, fmt.Errorf("error making build hook %s executable: %v", BuildFile, err)
	}

	sha, err := prov.Git().GetCurrentCommitFromRef(ctx, "HEAD")
	if err != nil {
		return nil, err
	}

	err = me.run(ctx, BuildFile, prov.Buildrc(), sv.Tag, sha, prov)
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
	return buildrc.RunAllPackagePlatforms(ctx, brc, 10*time.Minute, func(ctx context.Context, pkg *buildrc.Package, arc buildrc.Platform) error {

		dir, err := pipeline.NewTempDir(ctx, prov.Pipeline(), prov.FileSystem())
		if err != nil {
			return fmt.Errorf("error running script %s with [%s:%s]: %v", scriptPath, arc.OS(), arc.Arch(), err)
		}

		artifactName := pkg.Name + "-" + arc.OS() + "-" + arc.Arch()

		outputFile := filepath.Join(dir.String(), artifactName)

		custom, err := pkg.CustomJSON()
		if err != nil {
			return fmt.Errorf("error marshalling custom JSON: %v", err)
		}

		cmd := exec.Command("bash", "./"+scriptPath, outputFile, pkg.Name, custom)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = append(
			os.Environ(),
			fmt.Sprintf("GOOS=%s", arc.OS()),
			fmt.Sprintf("GOARCH=%s", arc.Arch()),
			fmt.Sprintf("GO_LDFLAGS=%s", ldflags),
			fmt.Sprintf("BUILDRC_VERSION=%s", tag),
			fmt.Sprintf("BUILDRC_COMMIT=%s", commit),
			fmt.Sprintf("BUILDRC_OS=%s", arc.OS()),
			fmt.Sprintf("BUILDRC_ARCH=%s", arc.Arch()),
			fmt.Sprintf("BUILDRC_OUTPUT=%s", outputFile),
			fmt.Sprintf("BUILDRC_CUSTOM=%s", custom),
			fmt.Sprintf("BUILDRC_NAME=%s", pkg.Name),
		)
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("error running script  %s with [%s:%s]: %v", scriptPath, arc.OS(), arc.Arch(), err)
		}

		if err = pipeline.UploadDirAsTar(ctx, prov.Pipeline(), prov.FileSystem(), dir.String(), artifactName, &pipeline.UploadDirAsTarOpts{
			RequireFiles:  true,
			ProduceSHA256: true,
		}); err != nil {
			return err
		}

		return nil
	})

}
