package buildrc

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/go-faster/errors"
)

type Platform struct {
	OS      string
	Arch    string
	Varient string
}

var (
	ErrCouldNotParsePlatform = errors.New("buildrc.ErrCouldNotParsePlatform")
)

func NewPlatformFromFullString(platform string) (*Platform, error) {
	parts := strings.Split(platform, "/")
	if len(parts) == 1 {
		parts = strings.Split(platform, "_")
	}
	if len(parts) == 1 {
		parts = strings.Split(platform, "-")
	}
	switch len(parts) {
	case 2:
		return &Platform{OS: parts[0], Arch: parts[1]}, nil
	case 3:
		return &Platform{OS: parts[0], Arch: parts[1], Varient: parts[2]}, nil
	default:
		return nil, errors.Wrap(ErrCouldNotParsePlatform, fmt.Sprintf("%q", platform))
	}
}

func (me *Platform) String() string {
	if me.Varient != "" {
		return me.OS + "/" + me.Arch + "/" + me.Varient
	}
	return me.OS + "/" + me.Arch
}

func (me *Platform) UnderscoreString() string {
	return strings.ReplaceAll(me.String(), "/", "_")
}

func (me *Platform) DashString() string {
	return strings.ReplaceAll(me.String(), "/", "-")
}

func GetGoPlatform(ctx context.Context) *Platform {
	osv := runtime.GOOS
	arch := runtime.GOARCH
	arm := os.Getenv("GOARM")

	plat := &Platform{
		OS:      osv,
		Arch:    arch,
		Varient: arm,
	}

	return plat
}

func GetTargetPlatform(ctx context.Context) (*Platform, error) {
	res := os.Getenv("TARGETPLATFORM")
	return NewPlatformFromFullString(res)
}

func GetBuildPlatform(ctx context.Context) (*Platform, error) {
	res := os.Getenv("BUILDPLATFORM")
	return NewPlatformFromFullString(res)
}

func (me *Platform) Aliases() []string {
	strs := []string{me.String(), me.UnderscoreString(), me.DashString()}

	if me.Varient == "v8" && me.Arch == "arm64" {
		similar := &Platform{OS: me.OS, Arch: "arm64", Varient: ""}
		strs = append(strs, similar.Aliases()...)
	}

	return strs
}
