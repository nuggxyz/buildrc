package buildrc

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

func RunAllPackagePlatforms(ctx context.Context, brc *Buildrc, to time.Duration, f func(ctx context.Context, pkg *Package, arc Platform) error) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(brc.Packages))

	for _, pkg := range brc.Packages {

		if len(pkg.Platforms) == 0 {
			zerolog.Ctx(ctx).Warn().Str("package", pkg.Name).Msg("no platforms specified for package")
			wg.Add(1)
			go func(pkg *Package) {
				defer wg.Done()
				errChan <- f(ctx, pkg, Platform("unknown/unknown"))
			}(pkg)
		}

		for _, arch := range pkg.Platforms {
			wg.Add(1)
			go func(pkg *Package, arch Platform) {
				defer wg.Done()
				errChan <- f(ctx, pkg, arch)
			}(pkg, arch)
		}
	}

	ctx, cancel := context.WithTimeout(ctx, to)
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
			if err != nil {
				zerolog.Ctx(ctx).Error().Err(err).Msg("error running build script")
				errs = append(errs, err)
			}
		}
	}

	if len(errs) > 0 {
		zerolog.Ctx(ctx).Error().Errs("errors", errs).Msg("completed with errors")
		return fmt.Errorf("completed with %d error(s)", len(errs))
	} else {
		zerolog.Ctx(ctx).Info().Msg("completed successfully")
		return nil
	}
}

func RunAllPackages(ctx context.Context, brc *Buildrc, to time.Duration, f func(ctx context.Context, pkg *Package) error) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(brc.Packages))

	for _, pkg := range brc.Packages {

		wg.Add(1)
		go func(pkg *Package) {
			defer wg.Done()
			errChan <- f(ctx, pkg)
		}(pkg)
	}

	ctx, cancel := context.WithTimeout(ctx, to)
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
			if err != nil {
				zerolog.Ctx(ctx).Error().Err(err).Msg("error running build script")
				errs = append(errs, err)
			}
		}
	}

	if len(errs) > 0 {
		zerolog.Ctx(ctx).Error().Errs("errors", errs).Msg("completed with errors")
		return fmt.Errorf("completed with %d error(s)", len(errs))
	} else {
		zerolog.Ctx(ctx).Info().Msg("completed successfully")
		return nil
	}
}
