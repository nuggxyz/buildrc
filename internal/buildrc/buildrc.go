package buildrc

import (
	"context"
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/nuggxyz/buildrc/internal/errd"
	"github.com/nuggxyz/buildrc/internal/provider"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
)

type BuildRC struct {
	Version  *semver.Version `yaml:"version,flow" json:"version"`
	Golang   *Golang         `yaml:"golang,flow" json:"golang"`
	Packages []*Package      `yaml:"packages,flow" json:"packages"`
}

type Golang struct {
	Version   *semver.Version `yaml:"version" json:"version"`
	Private   string          `yaml:"private" json:"private"`
	CacheMods bool            `yaml:"cache_mods" json:"cache"`
}

type Package struct {
	Type         PackageType     `yaml:"type" json:"type"`
	Language     PackageLanguage `yaml:"lang" json:"lang"`
	Name         string          `yaml:"name" json:"name"`
	Entry        string          `yaml:"entry" json:"entry"`
	Dockerfile   string          `yaml:"dockerfile" json:"dockerfile"`
	Platforms    []Platform      `yaml:"platforms" json:"platforms"`
	PrebuildHook string          `yaml:"prebuild_hook" json:"prebuild_hook"`

	PlatformsCSV string `yaml:"platforms_csv" json:"platforms_csv"`
}

var _ provider.Expressable = (*BuildRC)(nil)

func (me *BuildRC) Express() map[string]string {
	return map[string]string{
		"version": me.Version.String(),
	}
}

type Platform string

func Parse(ctx context.Context, src string) (cfg *BuildRC, err error) {

	defer errd.DeferContext(ctx, &err, "buildrc.Parse", src)

	cfg = &BuildRC{}

	data, err := load(ctx, src)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	zerolog.Ctx(ctx).Debug().Any("config", cfg).Msgf("buildrc: loaded %s", src)

	err = cfg.validate(ctx)

	return
}

type PackageType string

const (
	PackageTypeLambda    PackageType = "lambda"
	PackageTypeImage     PackageType = "image"
	PackageTypeContainer PackageType = "container"
	PackageTypeCLI       PackageType = "cli"
)

func (me PackageType) validate() error {

	options := []string{
		string(PackageTypeLambda),
		string(PackageTypeImage),
		string(PackageTypeContainer),
		string(PackageTypeCLI),
	}

	for _, o := range options {
		if o == string(me) {
			return nil
		}
	}

	return fmt.Errorf("invalid package type '%s', valid options are { %s }", string(me), strings.Join(options, " "))
}

type PackageLanguage string

const (
	PackageLanguageGo     PackageLanguage = "golang"
	PackageLanguageDocker PackageLanguage = "docker"
)

func (me PackageLanguage) validate() error {

	options := []string{
		string(PackageLanguageGo),
		string(PackageLanguageDocker),
	}

	for _, o := range options {
		if o == string(me) {
			return nil
		}
	}

	return fmt.Errorf("invalid package language '%s', valid options are { %s }", string(me), strings.Join(options, " "))
}
