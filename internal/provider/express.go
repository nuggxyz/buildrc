package provider

import "reflect"

const (
	ExpressTag = "express"
)

func Express(x interface{}) map[string]string {

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
