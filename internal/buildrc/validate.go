package buildrc

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/nuggxyz/buildrc/internal/errd"
)

func (c *BuildRC) validate(ctx context.Context) (err error) {

	defer errd.DeferContext(ctx, &err, "buildrc.Validate", c)

	if c.Version.Patch() != 0 {
		return fmt.Errorf("buildrc: invalid version: '%s' - must be major.minor", c.Version)
	}

	if err := c.Golang.validate(ctx); err != nil {
		return err
	}

	for _, pkg := range c.Packages {
		if err := pkg.validate(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (pkg *Package) validate(ctx context.Context) (err error) {
	if pkg.Name == "" {
		return errors.New("buildrc: no package name")
	}

	err = pkg.Type.validate()
	if err != nil {
		return err
	}

	err = pkg.Language.validate()
	if err != nil {
		return err
	}

	if pkg.Entry == "" {
		return errors.New("buildrc: no package file")
	}

	if s, err := pkg.EntryInfo(); err != nil {
		return err
	} else if s.IsDir() {
		return errors.New("buildrc: package file is a directory")
	} else if s.Size() == 0 {
		return errors.New("buildrc: package file is empty")
	} else {
		pkg.Entry, err = pkg.RelativeEntry()
		if err != nil {
			return err
		}
	}

	if pkg.Dockerfile != "" {
		if s, err := pkg.DockerfileInfo(); err != nil {
			return err
		} else if s.IsDir() {
			return errors.New("buildrc: dockerfile is a directory")
		} else if s.Size() == 0 {
			return errors.New("buildrc: dockerfile is empty")
		} else {
			pkg.Dockerfile, err = pkg.RelativeDockerfile()
			if err != nil {
				return err
			}
		}
	}

	if pkg.PrebuildHook != "" {
		if s, err := pkg.PrebuildHookInfo(); err != nil {
			return err
		} else if s.IsDir() {
			return errors.New("buildrc: prebuild_hook is a directory")
		} else if s.Size() == 0 {
			return errors.New("buildrc: prebuild_hook is empty")
		} else {
			pkg.PrebuildHook, err = pkg.RelativePrebuildHook()
			if err != nil {
				return err
			}
		}
	}

	if len(pkg.Platforms) == 0 {
		return errors.New("buildrc: no platforms")
	}

	for _, platform := range pkg.Platforms {
		if err := platform.validate(); err != nil {
			return err
		}
		pkg.PlatformsCSV += platform.String()
		if platform != pkg.Platforms[len(pkg.Platforms)-1] {
			pkg.PlatformsCSV += ","
		}
	}

	return nil
}

func (me *Golang) validate(ctx context.Context) (err error) {

	if me.Version.Major() < 1 {
		return errors.New("buildrc: golang version must be >= 1.x")
	}

	if me.Private == "" {
		return errors.New("buildrc: no golang private")
	}

	return nil
}

func (me Platform) validate() error {
	oss := strings.Split(string(me), "/")
	if len(oss) != 2 {
		return errors.New("invalid platform: " + string(me))
	}
	return nil
}
