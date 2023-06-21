package pipeline

import (
	"context"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/spf13/afero"
)

const (
	ExpressTag = "express"
)

type BuildrcEnvVar string

const (
	BuildrcCacheDir BuildrcEnvVar = "BUILDRC_CACHE_DIR"
	BuildrcTempDir  BuildrcEnvVar = "BUILDRC_TEMP_DIR"
)

func (me BuildrcEnvVar) Load(ctx context.Context, p Pipeline, fs afero.Fs) (string, error) {
	return p.GetFromEnv(ctx, string(me), fs)
}

func TempFileName(ctx context.Context, p Pipeline, fs afero.Fs, cmd string) (string, error) {
	r, err := BuildrcTempDir.Load(ctx, p, fs)
	if err != nil {
		return "", err
	}
	return filepath.Join(r, fmt.Sprintf("%s.provider-content.json", cmd)), nil
}

func CacheDir(ctx context.Context, p Pipeline, fs afero.Fs) (string, error) {
	r, err := BuildrcCacheDir.Load(ctx, p, fs)
	if err != nil {
		return "", err
	}
	return r, nil
}

func SetEnvFromCache(ctx context.Context, pipe Pipeline, fs afero.Fs) error {
	v, hit, err := loadCachedEnvVars(ctx, pipe, fs)
	if err != nil {
		return err
	}

	if !hit {

		for k, v := range v {
			err := pipe.AddToEnv(ctx, k, v, fs)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func AddContentToEnv(ctx context.Context, prov Pipeline, fsp afero.Fs, id string, cmd map[string]string) error {
	// save to tmp folder
	for k, v := range cmd {
		start := k
		if !strings.HasPrefix(k, "BUILDRC_") {
			start = fmt.Sprintf("BUILDRC_%s_%s", strings.ToUpper(id), strings.ToUpper(k))
		}
		err := cacheEnvVar(ctx, prov, fsp, start, v)
		if err != nil {
			return err
		}
		err = prov.AddToEnv(ctx, start, v, fsp)
		if err != nil {
			return err
		}
	}

	return nil
}

func AddContentToEnvButDontCache(ctx context.Context, prov Pipeline, fsp afero.Fs, id string, cmd map[string]string) error {
	// save to tmp folder
	for k, v := range cmd {
		start := k
		if !strings.HasPrefix(k, "BUILDRC_") {
			start = fmt.Sprintf("BUILDRC_%s_%s", strings.ToUpper(id), strings.ToUpper(k))
		}
		err := prov.AddToEnv(ctx, start, v, fsp)
		if err != nil {
			return err
		}
	}

	return nil
}

func Express(x any) map[string]string {

	ref := reflect.ValueOf(x)

	if ref.Kind() == reflect.Ptr {
		ref = ref.Elem()
	}

	if ref.Kind() != reflect.Struct {
		return nil
	}

	res := map[string]string{}

	for i := 0; i < ref.NumField(); i++ {
		field := ref.Type().Field(i)
		if field.Type.Kind() == reflect.String {
			if t, ok := field.Tag.Lookup(ExpressTag); ok && t != "-" && t != "" {
				res[t] = ref.Field(i).String()
			} else {
				res[field.Name] = ref.Field(i).String()
			}
		}
	}

	return res

}
