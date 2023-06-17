package build

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
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
	File      string `flag:"file" type:"file:" default:".buildrc"`
	JustBuild bool   `flag:"just-build" type:"bool" default:"false"`
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

	var wg sync.WaitGroup
	errChan := make(chan error, len(brc.Packages))

	for _, pkg := range brc.Packages {

		// make sure the prebuild hook exists and is executable
		if _, err := os.Stat(BuildFile); os.IsNotExist(err) {
			return nil, fmt.Errorf("build hook %s does not exist", BuildFile)
		}

		if err := os.Chmod(BuildFile, 0755); err != nil {
			return nil, fmt.Errorf("error making build hook %s executable: %v", BuildFile, err)
		}

		for _, arch := range pkg.Platforms {
			wg.Add(1)
			go me.runScript(ctx, BuildFile, ghclient, pkg, arch, &wg, errChan)
		}
	}

	ctx, cancel := context.WithTimeout(ctx, 600*time.Second)
	defer cancel()

	go func() {
		defer cancel()
		wg.Wait()
	}()

	errs := make([]error, 0)

HERE:
	for {
		select {
		case <-ctx.Done():
			if err := ctx.Err(); err != nil && err != context.Canceled {
				errs = append(errs, err)
			}
			break HERE
		case err := <-errChan:
			zerolog.Ctx(ctx).Error().Err(err).Msg("error running build script")
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		zerolog.Ctx(ctx).Error().Errs("errors", errs).Msg("completed with errors")
		return nil, fmt.Errorf("completed with %d error(s)", len(errs))
	} else {
		zerolog.Ctx(ctx).Info().Msg("completed successfully")
		return nil, nil
	}
}

func (me *Handler) runScript(ctx context.Context, scriptPath string, clnt *github.GithubClient, pkg *buildrc.Package, arc buildrc.Platform, wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()

	file, err := arc.OutputFile(pkg)
	if err != nil {
		errChan <- fmt.Errorf("error running script %s with [%s:%s]: %v", scriptPath, arc.OS(), arc.Arch(), err)
		return
	}

	cmd := exec.Command("bash", "./"+scriptPath, arc.OS(), arc.Arch(), file)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		errChan <- fmt.Errorf("error running script  %s with [%s:%s]: %v", scriptPath, arc.OS(), arc.Arch(), err)
		return
	}

	zerolog.Ctx(ctx).Debug().Msgf("ran script %s with [%s:%s]", scriptPath, arc.OS(), arc.Arch())

	if _, err := os.Stat(file); os.IsNotExist(err) {
		errChan <- fmt.Errorf("error running script %s with [%s:%s]: expected file %s to be created but it was not", scriptPath, arc.OS(), arc.Arch(), file)
		return
	}

	zerolog.Ctx(ctx).Debug().Msgf("script %s with [%s:%s] completed successfully", scriptPath, arc.OS(), arc.Arch())

	// Create .tar.gz archive at pkg.OutputFile(arc).tar.gz
	tarCmd := exec.Command("tar", "-czvf", file+".tar.gz", "-C", filepath.Dir(file), filepath.Base(file))
	tarCmd.Stdout = os.Stdout
	tarCmd.Stderr = os.Stderr
	err = tarCmd.Run()
	if err != nil {
		errChan <- fmt.Errorf("error creating .tar.gz archive: %v", err)
		return
	}

	zerolog.Ctx(ctx).Debug().Msgf("created archive %s.tar.gz", file)

	// Compute and write SHA-256 checksum to pkg.OutputFile(arc).sha256
	hashCmd := exec.Command("shasum", "-a", "256", file)
	hashOutput, err := hashCmd.Output()
	if err != nil {
		errChan <- fmt.Errorf("error computing SHA-256 checksum: %v", err)
		return
	}

	zerolog.Ctx(ctx).Debug().Msgf("computed SHA-256 checksum for %s", file)

	err = os.WriteFile(file+".sha256", hashOutput, 0644)
	if err != nil {
		errChan <- fmt.Errorf("error writing SHA-256 checksum to file: %v", err)
		return
	}

	if me.JustBuild {
		return
	}

	zerolog.Ctx(ctx).Debug().Msgf("wrote SHA-256 checksum to %s.sha256", file)

	err = clnt.Upload(ctx, file+".tar.gz")
	if err != nil {
		errChan <- fmt.Errorf("error uploading archive: %v", err)
		return
	}

	err = clnt.Upload(ctx, file+".sha256")
	if err != nil {
		errChan <- fmt.Errorf("error uploading checksum: %v", err)
		return
	}

	zerolog.Ctx(ctx).Debug().Msgf("uploaded checksum %s.sha256", file)

}
