package buildrc

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rs/zerolog"
	"github.com/walteh/buildrc/internal/errd"
	"gopkg.in/yaml.v3"
)

type Provider interface {
	Current() *Buildrc
}

type Buildrc struct {
	Version  int        `yaml:"version,flow" json:"version"`
	Packages []*Package `yaml:"packages,flow" json:"packages"`
	Aws      *Aws       `yaml:"aws,flow" json:"aws"`
	Github   *Github    `yaml:"github,flow" json:"github"`
}

func (me *Buildrc) Current() *Buildrc {
	return me
}

type Package struct {
	Language        PackageLanguage `yaml:"lang" json:"lang"`
	Name            string          `yaml:"name" json:"name"`
	Dir             string          `yaml:"dir" json:"dir"`
	Os              []string        `yaml:"os" json:"os"`
	Arch            []string        `yaml:"arch" json:"arch"`
	DockerPlatforms []Platform      `yaml:"docker_platforms" json:"docker_platforms"`
	Platforms       []Platform      `yaml:"platforms" json:"platforms"`
	Custom          any             `yaml:"custom" json:"custom"`
	Artifacts       []string        `yaml:"artifacts" json:"artifacts"`
	On              string          `yaml:"on" json:"on"`
	Uses            []string        `yaml:"uses" json:"uses"`
}

func (me *Buildrc) PackageByName() map[string]*Package {
	m := make(map[string]*Package)
	for _, pkg := range me.Packages {
		m[pkg.Name] = pkg
	}
	return m
}

type Platform string

func Parse(ctx context.Context, src string) (cfg *Buildrc, err error) {

	defer errd.DeferContext(ctx, &err, "buildrc.Parse", src)

	cfg = &Buildrc{}

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

func (me Platform) isDocker() bool {
	return strings.HasPrefix(string(me), "linux") || strings.Contains(string(me), "darwin")
}

func (me Platform) OutputFile(name *Package) (string, error) {
	main := fmt.Sprintf("%s-%s-%s", name.Name, me.OS(), me.Arch())

	if me.OS() == "windows" {
		main += ".exe"
	}
	return main, nil
}

type PackageType string

// const (
// 	PackageTypeLambda    PackageType = "lambda"
// 	PackageTypeImage     PackageType = "image"
// 	PackageTypeContainer PackageType = "container"
// 	PackageTypeCLI       PackageType = "cli"
// )

// func (me PackageType) validate() error {

// 	options := []string{
// 		string(PackageTypeLambda),
// 		string(PackageTypeImage),
// 		string(PackageTypeContainer),
// 		string(PackageTypeCLI),
// 	}

// 	for _, o := range options {
// 		if o == string(me) {
// 			return nil
// 		}
// 	}

// 	return fmt.Errorf("invalid package type '%s', valid options are { %s }", string(me), strings.Join(options, " "))
// }

type PackageLanguage string

const (
	PackageLanguageGo     PackageLanguage = "golang"
	PackageLanguageDocker PackageLanguage = "docker"
	PackageLanguageSwift  PackageLanguage = "swift"
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

func (me *Buildrc) GolangPackagesNamesArray() []string {
	strs := make([]string, len(me.Packages))
	for i, pkg := range me.Packages {
		if pkg.Language == PackageLanguageGo {
			strs[i] = pkg.Name
		}
	}
	return strs
}

func (me *Buildrc) PackagesNamesArray() []string {
	strs := make([]string, len(me.Packages))
	for i, pkg := range me.Packages {
		strs[i] = pkg.Name
	}
	return strs
}

func (me *Buildrc) PackagesNamesArrayJSON() string {
	return "[\"" + strings.Join(me.PackagesNamesArray(), "\",\"") + "\"]"
}

func (me *Buildrc) PackagesOnArray() []string {
	strs := make([]string, len(me.Packages))
	for i, pkg := range me.Packages {
		strs[i] = pkg.On
	}
	return strs
}

func (me *Buildrc) PackagesOnArrayJSON() string {
	return "[\"" + strings.Join(me.PackagesOnArray(), "\",\"") + "\"]"
}

func (me *Buildrc) PackagesArray() []*Package {
	return me.Packages
}

func (me *Buildrc) PackagesArrayJSON() (string, error) {
	data, err := json.Marshal(me.Packages)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (me *Buildrc) PackagesMap() map[string]*Package {
	m := make(map[string]*Package)
	for _, pkg := range me.Packages {
		m[pkg.Name] = pkg
	}
	return m
}

func (me *Buildrc) PackagesMapJSON() (string, error) {
	data, err := json.Marshal(me.PackagesMap())
	if err != nil {
		return "", err
	}
	return string(data), nil
}
