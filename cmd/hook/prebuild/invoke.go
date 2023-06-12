package prebuild

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/nuggxyz/buildrc/cmd/buildrc/load"
	"github.com/nuggxyz/buildrc/internal/buildrc"
	"github.com/nuggxyz/buildrc/internal/provider"
	"github.com/rs/zerolog"
)

type Handler struct {
	BuildrcFile string `arg:"buildrc_file" type:"file:" required:"true"`
	PackageName string `arg:"package_name" type:"string:" required:"true"`

	buildrcHandler *load.Handler
}

func (me *Handler) Init(ctx context.Context) (err error) {

	me.buildrcHandler, err = load.NewHandler(ctx, me.BuildrcFile)
	if err != nil {
		return err
	}
	return
}

func (me *Handler) Invoke(ctx context.Context, prv provider.ContentProvider) (out *output, err error) {

	brc, err := me.buildrcHandler.Helper().Run(ctx, prv)
	if err != nil {
		return nil, err
	}

	var pkg *buildrc.Package
	for _, p := range brc.Packages {
		if p.Name == me.PackageName {
			pkg = p
			break
		}
	}

	if pkg == nil {
		return nil, fmt.Errorf("package %s not found in buildrc", me.PackageName)
	}

	if pkg.PrebuildHook == "" {
		zerolog.Ctx(ctx).Debug().Msg("no prebuild hook defined, skipping")
		return &output{}, nil
	}

	// make sure the prebuild hook exists and is executable
	if _, err := os.Stat(pkg.PrebuildHook); os.IsNotExist(err) {
		return nil, fmt.Errorf("prebuild hook %s does not exist", pkg.PrebuildHook)
	}

	if err := os.Chmod(pkg.PrebuildHook, 0755); err != nil {
		return nil, fmt.Errorf("error making prebuild hook %s executable: %v", pkg.PrebuildHook, err)
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(pkg.Platforms))

	for _, arch := range pkg.Platforms {
		wg.Add(1)
		go runScript(pkg.PrebuildHook, pkg, arch, &wg, errChan)
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
		return &output{}, nil
	}
}

func runScript(scriptPath string, pkg *buildrc.Package, arc buildrc.Platform, wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()

	cmd := exec.Command("bash", "./"+scriptPath, pkg.Entry, arc.OS(), arc.Arch())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		errChan <- fmt.Errorf("error running script  %s with [%s:%s:%s]: %v", scriptPath, pkg.Entry, arc.OS(), arc.Arch(), err)
		return
	}
}
