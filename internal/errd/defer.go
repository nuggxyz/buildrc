package errd

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
)

// Wrap wraps err with fmt.Errorf if err is non nil.
// Intended for use with defer and a named error return.
// Inspired by https://github.com/golang/go/issues/32676.
func Wrap(err *error, f string, v ...interface{}) {
	if *err != nil {
		*err = fmt.Errorf(f+": %w", append(v, *err)...)
		// *err = New(f).WithKV("data", v).WithRoot(*err)
	}
}

// Wrap wraps err with fmt.Errorf if err is non nil.
// Intended for use with defer and a named error return.
// Inspired by https://github.com/golang/go/issues/32676.
func DeferContext(ctx context.Context, err *error, f string, v ...interface{}) {
	if *err != nil {
		WrapContext(ctx, *err, f, v...)
	}
}

func WrapContext(ctx context.Context, err error, f string, v ...interface{}) error {
	if err != nil {
		ev := zerolog.Ctx(ctx).Error().Err(err).CallerSkipFrame(2).Caller()
		for i := range v {
			ev = ev.Interface(fmt.Sprintf("data[%d]", i), v[i])
		}
		ev.Msg(f)

		return err
	}
	return nil
}
