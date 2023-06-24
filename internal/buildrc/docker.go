package buildrc

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/nuggxyz/buildrc/internal/pipeline"
	"github.com/spf13/afero"
)

type DockerBuildArgs map[string]string

func (me *Package) DockerBuildArgs(ctx context.Context, p pipeline.Pipeline, fs afero.Fs) (DockerBuildArgs, error) {

	cachedir, err := pipeline.BuildrcCacheDir.Load(ctx, p, fs)
	if err != nil {
		panic(err)
	}

	return map[string]string{
		"DIR":  cachedir,
		"NAME": me.Name,
	}, nil
}

func (me DockerBuildArgs) Array() []string {

	var strArgs []string
	for k, v := range me {
		strArgs = append(strArgs, fmt.Sprintf("%s=%s", k, v))
	}

	return strArgs
}

func (me DockerBuildArgs) CSV() (string, error) {

	return strings.Join(me.Array(), ","), nil
}

func (me DockerBuildArgs) JSONString() (string, error) {
	args := me.Array()

	joiner := strings.Join(args, "\n")

	res, err := json.Marshal(joiner)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func (me *Package) Dockerfile() string {
	return filepath.Join(me.Dir, "Dockerfile")
}

func (me *Package) DockerPlatformsCSV() string {
	return StringsToCSV(me.DockerPlatforms)
}

func (me *Buildrc) Images(pkg *Package, org string, repo string) []string {

	var last string

	if repo == pkg.Name {
		last = pkg.Name
	} else {
		last = fmt.Sprintf("%s/%s", repo, pkg.Name)
	}
	strs := make([]string, 0)
	if me.Aws != nil {
		strs = append(strs, me.Aws.Repository(pkg, org, repo))
	}

	strs = append(strs, fmt.Sprintf("ghcr.io/%s/%s", org, last))

	return strs

}

func (me *Buildrc) ImagesCSV(pkg *Package, org string, repo string) string {
	return strings.Join(me.Images(pkg, org, repo), ",")
}

func (me *Buildrc) ImagesCSVJSON(pkg *Package, org string, repo string) (string, error) {
	data, err := json.Marshal(me.ImagesCSV(pkg, org, repo))
	if err != nil {
		return "", err
	}

	return string(data), nil

}
