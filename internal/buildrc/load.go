package buildrc

import (
	"context"
	"errors"
	"os"

	"github.com/walteh/buildrc/internal/errd"
)

func load(ctx context.Context, file string) (res []byte, err error) {

	defer errd.DeferContext(ctx, &err, "buildrc.Load", file)

	// verify the .buildrc file exists
	stat, err := os.Stat(file)
	if err != nil {
		return nil, err
	}

	if stat.IsDir() {
		err = errors.New("buildrc: file is a directory")
		return
	}

	if stat.Size() == 0 {
		err = errors.New("buildrc: file is empty")
		return
	}

	res, err = os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return
}
