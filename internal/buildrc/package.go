package buildrc

import (
	"os"
	"path/filepath"
	"strings"
)

func (me *Package) AbsolutePrebuildHook() (string, error) {
	return filepath.Abs(me.PrebuildHook)
}

func (me *Package) PrebuildHookInfo() (os.FileInfo, error) {
	return os.Stat(me.PrebuildHook)
}

func (me *Package) RelativePrebuildHook() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	abs, err := me.AbsolutePrebuildHook()
	if err != nil {
		return "", err
	}
	return filepath.Rel(cwd, abs)
}

func (me *Package) AbsoluteEntry() (string, error) {
	return filepath.Abs(me.Entry)
}

func (me *Package) EntryInfo() (os.FileInfo, error) {
	return os.Stat(me.Entry)
}

func (me *Package) RelativeEntry() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	abs, err := me.AbsoluteEntry()
	if err != nil {
		return "", err
	}
	return filepath.Rel(cwd, abs)
}

func (me *Package) AbsoluteDockerfile() (string, error) {
	return filepath.Abs(me.Dockerfile)
}

func (me *Package) RelativeDockerfile() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	abs, err := me.AbsoluteDockerfile()
	if err != nil {
		return "", err
	}
	return filepath.Rel(cwd, abs)
}

func (me *Package) DockerfileInfo() (os.FileInfo, error) {
	return os.Stat(me.Dockerfile)
}

func (me *Platform) String() string {
	return me.OS() + "/" + me.Arch()
}

func (me Platform) OS() string {
	oss := strings.Split(string(me), "/")
	return oss[0]
}

func (me Platform) Arch() string {
	oss := strings.Split(string(me), "/")
	return oss[1]
}
