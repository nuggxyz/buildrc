package build

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

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

	var wg sync.WaitGroup
	errChan := make(chan error)

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
			go runScript(ctx, BuildFile, ghclient, pkg, arch, &wg, errChan)
		}
	}

	wg.Wait()
	close(errChan)

	errors := 0
	for err := range errChan {
		errors++
		zerolog.Ctx(ctx).Error().Err(err).Msg("error running prebuild hook")
	}

	if errors > 0 {
		zerolog.Ctx(ctx).Error().Int("errors", errors).Msg("completed with errors")
		return nil, fmt.Errorf("completed with %d error(s)", errors)
	} else {
		fmt.Println("All architectures completed successfully")
		return nil, nil
	}
}

func runScript(ctx context.Context, scriptPath string, clnt *github.GithubClient, pkg *buildrc.Package, arc buildrc.Platform, wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()

	file := arc.OutputFile(pkg)

	cmd := exec.Command("bash", "./"+scriptPath, arc.OS(), arc.Arch(), file)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
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

	// Compute and write SHA-256 checksum to pkg.OutputFile(arc).sha256
	hashCmd := exec.Command("shasum", "-a", "256", file)
	hashOutput, err := hashCmd.Output()
	if err != nil {
		errChan <- fmt.Errorf("error computing SHA-256 checksum: %v", err)
		return
	}

	err = os.WriteFile(file+".sha256", hashOutput, 0644)
	if err != nil {
		errChan <- fmt.Errorf("error writing SHA-256 checksum to file: %v", err)
		return
	}

	tar, err := os.Open(file + ".tar.gz")
	if err != nil {
		errChan <- fmt.Errorf("error opening archive file: %v", err)
		return
	}

	_, _, err = clnt.UploadWorkflowAsset(ctx, file+".tar.gz", tar)
	if err != nil {
		errChan <- fmt.Errorf("error uploading archive: %v", err)
		return
	}

	zerolog.Ctx(ctx).Debug().Msgf("uploaded archive %s.tar.gz", file)

	sha, err := os.Open(file + ".sha256")
	if err != nil {
		errChan <- fmt.Errorf("error opening checksum file: %v", err)
		return
	}

	_, _, err = clnt.UploadWorkflowAsset(ctx, file+".sha256", sha)
	if err != nil {
		errChan <- fmt.Errorf("error uploading checksum: %v", err)
		return
	}

	zerolog.Ctx(ctx).Debug().Msgf("uploaded checksum %s.sha256", file)

}
