package build

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/nuggxyz/buildrc/cmd/buildrc/load"
	"github.com/nuggxyz/buildrc/internal/buildrc"
	"github.com/nuggxyz/buildrc/internal/provider"
	"github.com/rs/zerolog"
)

type Handler struct {
	File string `flag:"file" type:"file:" default:".buildrc"`
}

func (me *Handler) Invoke(ctx context.Context, prv provider.ContentProvider) (out any, err error) {

	brc, err := load.NewHandler(me.File).Load(ctx, prv)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	errChan := make(chan error)

	for _, pkg := range brc.Packages {

		if pkg.PrebuildHook == "" {
			zerolog.Ctx(ctx).Debug().Msg("no prebuild hook defined, skipping")
			return nil, nil
		}

		// make sure the prebuild hook exists and is executable
		if _, err := os.Stat(pkg.PrebuildHook); os.IsNotExist(err) {
			return nil, fmt.Errorf("prebuild hook %s does not exist", pkg.PrebuildHook)
		}

		if err := os.Chmod(pkg.PrebuildHook, 0755); err != nil {
			return nil, fmt.Errorf("error making prebuild hook %s executable: %v", pkg.PrebuildHook, err)
		}

		for _, arch := range pkg.Platforms {
			wg.Add(1)
			go runScript(pkg.PrebuildHook, pkg, arch, &wg, errChan)
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

func runScript(scriptPath string, pkg *buildrc.Package, arc buildrc.Platform, wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()

	file := arc.OutputFile(pkg)

	cmd := exec.Command("bash", "./"+scriptPath, pkg.Entry, arc.OS(), arc.Arch(), file)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		errChan <- fmt.Errorf("error running script  %s with [%s:%s:%s]: %v", scriptPath, pkg.Entry, arc.OS(), arc.Arch(), err)
		return
	}

	zerolog.Ctx(context.Background()).Debug().Msgf("ran script %s with [%s:%s:%s]", scriptPath, pkg.Entry, arc.OS(), arc.Arch())

	if _, err := os.Stat(file); os.IsNotExist(err) {
		errChan <- fmt.Errorf("error running script %s with [%s:%s:%s]: expected file %s to be created but it was not", scriptPath, pkg.Entry, arc.OS(), arc.Arch(), file)
		return
	}

	zerolog.Ctx(context.Background()).Debug().Msgf("script %s with [%s:%s:%s] completed successfully", scriptPath, pkg.Entry, arc.OS(), arc.Arch())

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
	err = ioutil.WriteFile(file+".sha256", hashOutput, 0644)
	if err != nil {
		errChan <- fmt.Errorf("error writing SHA-256 checksum to file: %v", err)
		return
	}

}
