package main

import (
	"context"
	"fmt"
)

var (
	// Package is filled at linking time
	Package = "github.com/walteh/buildrc"

	// Version holds the complete version number. Filled in at linking time.
	Version = ""

	PR = ""

	ContentHash = ""

	Temporary = ""
)

type VersionHandler struct {
}

func (me *VersionHandler) Run(ctx context.Context) (err error) {

	if Version == "" {
		if Temporary != "" {
			_, err = fmt.Println("temporary")
		} else {
			_, err = fmt.Println("unknown")
		}
	} else {
		_, err = fmt.Println(Version)
	}

	return err
}
