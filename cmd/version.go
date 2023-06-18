package main

import (
	"context"
	"fmt"
)

var (
	// Package is filled at linking time
	Package = "github.com/nuggxyz/buildrc"

	// Version holds the complete version number. Filled in at linking time.
	Version = ""

	// Revision is filled with the VCS (e.g. git) revision being used to build
	// the program at linking time.
	Revision = ""

	Time = ""

	RawVersion = ""

	Temporary = ""
)

type VersionHandler struct {
	Revision bool `flag:"revision" type:"revision:" default:"false"`
	Time     bool `flag:"time" type:"time:" default:"false"`
	Raw      bool `flag:"raw" type:"raw:" default:"false"`
}

func (me *VersionHandler) Run(ctx context.Context) (err error) {

	if Version == "" {
		if Temporary != "" {
			_, err = fmt.Println("temporary")
		} else {
			_, err = fmt.Println("unknown")
		}
	} else {
		if me.Raw {
			_, err = fmt.Println(RawVersion)
		} else if me.Revision {
			_, err = fmt.Println(Revision)
		} else if me.Time {
			_, err = fmt.Println(Time)
		} else {
			_, err = fmt.Println(Version)
		}
	}

	return err
}
