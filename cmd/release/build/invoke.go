package build

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/nuggxyz/buildrc/cmd/buildrc/load"
	"github.com/nuggxyz/buildrc/internal/buildrc"
	"github.com/nuggxyz/buildrc/internal/github"
	"github.com/nuggxyz/buildrc/internal/provider"

	"github.com/rs/zerolog"
)

const (
	CommandID = "build"
	BuildFile = "build.sh"
)

type Handler struct {
	File string `flag:"file" type:"file:" default:".buildrc"`
}

func (me *Handler) Run(ctx context.Context, cp provider.ContentProvider) (err error) {
	_, err = me.Build(ctx, cp)
	return err
}

func (me *Handler) Build(ctx context.Context, cp provider.ContentProvider) (out *any, err error) {

	return provider.Wrap(CommandID, me.build)(ctx, cp)
}

func (me *Handler) build(ctx context.Context, prv provider.ContentProvider) (out *any, err error) {

	brc, err := load.NewHandler(me.File).Load(ctx, prv)
	if err != nil {
		return nil, err
	}

	ghclient, err := github.NewGithubClient(ctx, "", "")
	if err != nil {
		return nil, err
	}

	zerolog.Ctx(ctx).Info().Msg("checking if build is required")

	ok, reason, err := ghclient.ShouldBuild(ctx)
	if err != nil {
		return nil, err
	}

	if !ok {
		zerolog.Ctx(ctx).Info().Str("reason", reason).Msg("build not required")
		return nil, nil
	} else {
		zerolog.Ctx(ctx).Info().Str("reason", reason).Msg("build required, continuing")
	}

	// make sure the prebuild hook exists and is executable
	if _, err := os.Stat(BuildFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("build hook %s does not exist", BuildFile)
	}

	if err := os.Chmod(BuildFile, 0755); err != nil {
		return nil, fmt.Errorf("error making build hook %s executable: %v", BuildFile, err)
	}

	err = me.run(ctx, BuildFile, ghclient, brc)
	if err != nil {
		return nil, err
	}

	return nil, nil

}

func (me *Handler) run(ctx context.Context, scriptPath string, clnt *github.GithubClient, brc *buildrc.BuildRC) error {
	return buildrc.RunAllPackages(ctx, brc, 10*time.Minute, func(ctx context.Context, pkg *buildrc.Package, arc buildrc.Platform) error {
		file, err := arc.OutputFile(pkg)
		if err != nil {
			return fmt.Errorf("error running script %s with [%s:%s]: %v", scriptPath, arc.OS(), arc.Arch(), err)
		}

		cmd := exec.Command("bash", "./"+scriptPath, file)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = append(os.Environ(), fmt.Sprintf("GOOS=%s", arc.OS()), fmt.Sprintf("GOARCH=%s", arc.Arch()))

		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("error running script  %s with [%s:%s]: %v", scriptPath, arc.OS(), arc.Arch(), err)
		}

		zerolog.Ctx(ctx).Debug().Msgf("ran script %s with [%s:%s]", scriptPath, arc.OS(), arc.Arch())

		if _, err := os.Stat(file); os.IsNotExist(err) {
			return fmt.Errorf("error running script %s with [%s:%s]: expected file %s to be created but it was not", scriptPath, arc.OS(), arc.Arch(), file)
		}

		zerolog.Ctx(ctx).Debug().Msgf("script %s with [%s:%s] completed successfully", scriptPath, arc.OS(), arc.Arch())

		// Create .tar.gz archive at pkg.OutputFile(arc).tar.gz
		tarCmd := exec.Command("tar", "-czvf", file+".tar.gz", "-C", filepath.Dir(file), filepath.Base(file))
		tarCmd.Stdout = os.Stdout
		tarCmd.Stderr = os.Stderr
		err = tarCmd.Run()
		if err != nil {
			return fmt.Errorf("error creating .tar.gz archive: %v", err)
		}

		zerolog.Ctx(ctx).Debug().Msgf("created archive %s.tar.gz", file)

		// Compute and write SHA-256 checksum to pkg.OutputFile(arc).sha256
		hashCmd := exec.Command("shasum", "-a", "256", file)
		hashOutput, err := hashCmd.Output()
		if err != nil {
			return fmt.Errorf("error computing SHA-256 checksum: %v", err)
		}

		zerolog.Ctx(ctx).Debug().Msgf("computed SHA-256 checksum for %s", file)

		err = os.WriteFile(file+".sha256", hashOutput, 0644)
		if err != nil {
			return fmt.Errorf("error writing SHA-256 checksum to file: %v", err)
		}

		zerolog.Ctx(ctx).Debug().Msgf("wrote SHA-256 checksum to %s.sha256", file)

		return nil
	})

	// if me.JustBuild {
	// 	return
	// }

	// zerolog.Ctx(ctx).Debug().Msgf("wrote SHA-256 checksum to %s.sha256", file)

	// err = clnt.Upload(ctx, file+".tar.gz")
	// if err != nil {
	// 	return fmt.Errorf("error uploading archive: %v", err)
	// }

	// err = clnt.Upload(ctx, file+".sha256")
	// if err != nil {
	// 	return fmt.Errorf("error uploading checksum: %v", err)
	// }

	// zerolog.Ctx(ctx).Debug().Msgf("uploaded checksum %s.sha256", file)

}
