package buildrc

import (
	"errors"
	"os"
)

type BuildrcEnvVar string

const (
	BuildrcCacheDir BuildrcEnvVar = "BUILDRC_CACHE_DIR"
	BuildrcTempDir  BuildrcEnvVar = "BUILDRC_TEMP_DIR"
)

func (me BuildrcEnvVar) Load() (string, error) {
	res := os.Getenv(string(me))
	if res == "" {
		return "", errors.New("env var not set or empty: " + string(me))
	}
	return res, nil
}
