package version

import (
	"context"
	"fmt"

	"github.com/nuggxyz/buildrc/internal/provider"
	"github.com/nuggxyz/buildrc/version"
)

const (
	CommandID = "version"
)

type Handler struct {
	Revision bool `flag:"revision" type:"revision:" default:"false"`
	Time     bool `flag:"time" type:"time:" default:"false"`
	Raw      bool `flag:"raw" type:"raw:" default:"false"`
}

func (me *Handler) Run(ctx context.Context, cp provider.ContentProvider) (err error) {
	check := version.Version

	if check == "" {
		if version.Temporary {
			_, err = fmt.Println("temporary")
		} else {
			_, err = fmt.Println("unknown")
		}
	} else {
		if me.Raw {
			_, err = fmt.Println(version.RawVersion)
		} else if me.Revision {
			_, err = fmt.Println(version.Revision)
		} else if me.Time {
			_, err = fmt.Println(version.Time)
		} else {
			_, err = fmt.Println(version.Version)
		}
	}

	return err
}
