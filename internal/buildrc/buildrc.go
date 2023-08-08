package buildrc

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nuggxyz/buildrc/internal/errd"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
)

type Buildrc struct {
	Version  int        `yaml:"version,flow" json:"version"`
	Packages []*Package `yaml:"packages,flow" json:"packages"`
	Aws      *Aws       `yaml:"aws,flow" json:"aws"`
	Github   *Github    `yaml:"github,flow" json:"github"`
	On       string     `yaml:"on" json:"on"`
	Uses     []string   `yaml:"uses" json:"uses"`
}

type Package struct {
	Type     PackageType     `yaml:"type" json:"type"`
	Language PackageLanguage `yaml:"lang" json:"lang"`
	Name     string          `yaml:"name" json:"name"`
	// Dockerfile      string          `yaml:"dockerfile" json:"dockerfile"`
	Dir             string     `yaml:"dir" json:"dir"`
	Os              []string   `yaml:"os" json:"os"`
	Arch            []string   `yaml:"arch" json:"arch"`
	DockerPlatforms []Platform `yaml:"docker_platforms" json:"docker_platforms"`
	Platforms       []Platform `yaml:"platforms" json:"platforms"`
	Custom          any        `yaml:"custom" json:"custom"`
	Artifacts       []string   `yaml:"artifacts" json:"artifacts"`
}

func (me *Buildrc) PackageByName() map[string]*Package {
	m := make(map[string]*Package)
	for _, pkg := range me.Packages {
		m[pkg.Name] = pkg
	}
	return m
}

func (me *Buildrc) UsesMap() map[string]string {
	m := make(map[string]string)
	for _, use := range me.Uses {
		m["uses_"+use] = "1"
	}
	return m
}

func StringsToCSV[I ~string](ss []I) string {
	strs := make([]string, len(ss))
	for i, s := range ss {
		strs[i] = string(s)
	}
	return strings.Join(strs, ",")
}

func (me *Package) ArtifactFileNames() ([]string, error) {
	names := make([]string, 0)
	for _, s := range me.Platforms {
		tmp, err := s.OutputFile(me)
		if err != nil {
			return nil, err
		}
		names = append(names, tmp+".tar.gz", tmp+".sha256")
	}
	return names, nil
}

func (me *Package) ToArtifactCSV(ss []Platform) (string, error) {
	names, err := me.ArtifactFileNames()
	if err != nil {
		return "", err
	}

	return strings.Join(names, ","), nil
}

func (me *Package) CustomJSON() (string, error) {
	if me.Custom == nil {
		return "{}", nil
	}
	cust, err := json.Marshal(me.Custom)
	if err != nil {
		return "", err
	}

	return string(cust), nil
}

func (me *Package) TestArchiveFileName() string {
	return fmt.Sprintf("%s-test-output.tar.gz", me.Name)
}

func (me *Package) VerifyArchiveFileName() string {
	return fmt.Sprintf("%s-test-output.tar.gz", me.Name)
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
	return strings.HasPrefix(string(me), "linux")
}

func (me Platform) OutputFile(name *Package) (string, error) {
	main := fmt.Sprintf("%s-%s-%s", name.Name, me.OS(), me.Arch())

	if me.OS() == "windows" {
		main += ".exe"
	}
	return main, nil
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
