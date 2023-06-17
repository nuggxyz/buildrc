package provider

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

const (
	ExpressTag = "express"
)

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

func AddContentToEnv(ctx context.Context, prov ContentProvider, id string, cmd map[string]string) error {
	// save to tmp folder
	for k, v := range cmd {
		start := k
		if !strings.HasPrefix(k, "BUILDRC_") {
			start = fmt.Sprintf("BUILDRC_%s_%s", strings.ToUpper(id), strings.ToUpper(k))
		}
		err := prov.AddToEnv(ctx, start, v)
		if err != nil {
			return err
		}
	}

	return nil
}
