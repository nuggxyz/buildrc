package buildrc

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/walteh/buildrc/internal/errd"
)

func (c *Buildrc) validate(ctx context.Context) (err error) {

	defer errd.DeferContext(ctx, &err, "buildrc.Validate", c)

	if c.Packages != nil {
		if len(c.Packages) != 1 {
			// TODO: support multiple packages
			return errors.New("buildrc: only one package is supported")
		}
		for _, pkg := range c.Packages {
			if err := pkg.validate(ctx); err != nil {
				return err
			}
		}
	}

	if len(c.Packages) != len(c.PackageByName()) {
		return errors.New("buildrc: duplicate package names")
	}

	return nil
}

func (pkg *Package) validate(ctx context.Context) (err error) {
	if pkg.Name == "" {
		return errors.New("buildrc: no package name")
	}

	if pkg.Dir == "" {
		pkg.Dir = "."
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	pkg.Dir = filepath.Join(cwd, pkg.Dir)

	if len(pkg.Os) > 0 && len(pkg.Arch) > 0 {
		pkg.Platforms = make([]Platform, 0, len(pkg.Os)*len(pkg.Arch))
		for _, os := range pkg.Os {
			for _, arch := range pkg.Arch {
				pkg.Platforms = append(pkg.Platforms, Platform(os+"/"+arch))
			}
		}
	}

	for _, platform := range pkg.Platforms {

		if err := platform.validate(); err != nil {
			return err
		}
		if platform.isDocker() {
			pkg.DockerPlatforms = append(pkg.DockerPlatforms, platform)
		}
	}

	return nil
}

// func (me *Golang) validate(ctx context.Context) (err error) {

// 	if me.Version.Major() < 1 {
// 		return errors.New("buildrc: golang version must be >= 1.x")
// 	}

// 	if me.Private == "" {
// 		return errors.New("buildrc: no golang private")
// 	}

// 	return nil
// }

func (me Platform) validate() error {
	oss := strings.Split(string(me), "/")
	if len(oss) != 2 {
		return errors.New("invalid platform: " + string(me))
	}
	return nil
}
