package buildrc

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
)

func (me *Package) DockerBuildArgs() (map[string]string, error) {

	cachedir, err := BuildrcCacheDir.Load()
	if err != nil {
		panic(err)
	}

	// args = append(args, fmt.Sprintf("DIR=%s", cachedir))

	// args = append(args, fmt.Sprintf("NAME=%s", me.Name))

	return map[string]string{
		"DIR":  cachedir,
		"NAME": me.Name,
	}, nil
}

func (me *Package) DockerBuildArgsArray() ([]string, error) {
	args, err := me.DockerBuildArgs()
	if err != nil {
		return nil, err
	}

	var strArgs []string
	for k, v := range args {
		strArgs = append(strArgs, fmt.Sprintf("%s=%s", k, v))
	}

	return strArgs, nil
}

func (me *Package) DockerBuildArgsCSV() (string, error) {
	args, err := me.DockerBuildArgsArray()
	if err != nil {
		return "", err
	}

	return strings.Join(args, ","), nil
}

func (me *Package) Dockerfile() string {
	return filepath.Join(me.Dir, "Dockerfile")
}

func (me *Package) DockerPlatformsCSV() string {
	return StringsToCSV(me.DockerPlatforms)
}

func (me *BuildRC) Images(pkg *Package, org string, repo string) []string {
	strs := make([]string, 0)
	if me.Aws != nil {
		strs = append(strs, fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com/%s/%s/%s", me.Aws.AccountID, me.Aws.Region, org, repo, pkg.Name))
	}

	strs = append(strs, fmt.Sprintf("ghcr.io/%s/%s/%s", org, repo, pkg.Name))

	return strs

}

func (me *BuildRC) ImagesCSV(pkg *Package, org string, repo string) string {
	return strings.Join(me.Images(pkg, org, repo), ",")
}

func (me *BuildRC) ImagesCSVJSON(pkg *Package, org string, repo string) (string, error) {
	data, err := json.Marshal(me.ImagesCSV(pkg, org, repo))
	if err != nil {
		return "", err
	}

	return string(data), nil

}
