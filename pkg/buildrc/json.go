package buildrc

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/rs/zerolog"
	"github.com/spf13/afero"
	"github.com/walteh/buildrc/pkg/git"
)

var (
	ErrCouldNotParseRemoteURL = fmt.Errorf("could not parse remote url")
)

type BuildrcJSON struct {
	Version    string `json:"version"`
	Revision   string `json:"revision"`
	Executable string `json:"executable"`
	Org        string `json:"org"`
	Artifact   string `json:"artifact"`
	GoPkg      string `json:"go-pkg"`
	Name       string `json:"name"`
	Image      string `json:"image"`
	Platform   string `json:"platform"`
}

type BuildrcPackageName string

type BuildrcVersion string

func GetArtifactName(ctx context.Context, name string, version string, plat string) string {
	return name + "-" + version + "-" + strings.ReplaceAll(plat, "/", "-")
}

func GetPlatform(ctx context.Context) (string, error) {
	res := os.Getenv("TARGETPLATFORM")
	if res == "" {
		osv := os.Getenv("GOOS")
		arch := os.Getenv("GOARCH")
		arm := os.Getenv("GOARM")
		if arm != "" {
			arch = arch + "/" + arm
		}
		return osv + "/" + arch, nil
	}
	return res, nil
}

func GetRevision(ctx context.Context, gitp git.GitProvider) (string, error) {

	revision, err := gitp.GetCurrentCommitFromRef(ctx, "HEAD")
	if err != nil {
		return "", err
	}

	return revision, nil
}

func GetExecutable(ctx context.Context, name string) string {
	if runtime.GOOS == "windows" {
		return name + ".exe"
	}
	return name
}

func GetRepo(ctx context.Context, gitp git.GitProvider) (string, string, error) {

	url, err := gitp.GetRemoteURL(ctx)
	if err != nil {
		return "", "", err
	}

	parts := strings.Split(url, "/")

	if len(parts) < 2 {
		return "", "", ErrCouldNotParseRemoteURL
	}

	org := parts[len(parts)-2]

	if strings.Contains(org, ":") && !strings.HasSuffix(org, ":") {
		org = strings.Split(org, ":")[1]
	}

	trimmed := strings.TrimSuffix(parts[len(parts)-1], ".git")

	return org, trimmed, nil
}

func GetGoPkg(ctx context.Context, gitp git.GitProvider) (string, error) {

	fle, err := afero.ReadFile(gitp.Fs(), "go.mod")
	if err != nil {
		return "", err
	}

	// find the line with module on it
	lines := strings.Split(string(fle), "\n")

	var modine string

	for _, line := range lines {
		if strings.HasPrefix(line, "module") {
			modine = line
			break
		}
	}

	if modine == "" {
		return "", fmt.Errorf("could not find module line in go.mod")
	}

	// split on space
	parts := strings.Split(modine, " ")

	if len(parts) != 2 {
		return "", fmt.Errorf("could not parse module line in go.mod")
	}

	return parts[1], nil

}

func GetBuildrcJSON(ctx context.Context, gitp git.GitProvider, opts *GetVersionOpts) (*BuildrcJSON, error) {

	brc, err := LoadBuildrc(ctx, gitp)
	if err != nil {
		zerolog.Ctx(ctx).Debug().Err(err).Msg("could not load buildrc")
		return nil, err
	}

	plat, err := GetPlatform(ctx)
	if err != nil {
		zerolog.Ctx(ctx).Debug().Err(err).Msg("could not get platform")
		return nil, err
	}

	version, err := GetVersion(ctx, gitp, brc, opts)
	if err != nil {
		zerolog.Ctx(ctx).Debug().Err(err).Msg("could not get version")
		return nil, err
	}

	org, name, err := GetRepo(ctx, gitp)
	if err != nil {
		zerolog.Ctx(ctx).Debug().Err(err).Msg("could not get repo")
		return nil, err
	}

	revision, err := GetRevision(ctx, gitp)
	if err != nil {
		zerolog.Ctx(ctx).Debug().Err(err).Msg("could not get revision")
		return nil, err
	}

	goPkg, err := GetGoPkg(ctx, gitp)
	if err != nil {
		zerolog.Ctx(ctx).Debug().Err(err).Msg("could not get go pkg")
		return nil, err
	}

	exec := GetExecutable(ctx, name)

	artif := GetArtifactName(ctx, name, version, plat)

	return &BuildrcJSON{
		Version:    version,
		Revision:   revision,
		Executable: exec,
		Image:      org + "/" + name,
		Artifact:   artif,
		GoPkg:      goPkg,
		Name:       name,
		Org:        org,
		Platform:   plat,
	}, nil
}
