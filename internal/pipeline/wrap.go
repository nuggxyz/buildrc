package pipeline

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog"
	"github.com/spf13/afero"
)

type PipelineProvider interface {
	Pipeline() Pipeline
	FileSystem() afero.Fs
}

type GenericRunnerFunc[I any, O any] func(context.Context, I) (*O, error)

func WrapGeneric[I any, O any](ctx context.Context, id string, cp Pipeline, fs afero.Fs, in I, i GenericRunnerFunc[I, O]) (*O, error) {

	return wrap(ctx, id, i, cp, in, fs)
}

func Cache[I PipelineProvider, O any](ctx context.Context, id string, in I, i GenericRunnerFunc[I, O]) (out *O, err error) {
	zerolog.Ctx(ctx).Debug().Str("id", id).Msg("Cache")
	defer func() {
		if err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Str("id", id).Msg("Cache Done with error")
		} else {
			zerolog.Ctx(ctx).Debug().Str("id", id).Msg("Cache Done")
		}
	}()

	out, err = WrapGeneric(ctx, id, in.Pipeline(), in.FileSystem(), in, i)
	return
}

func wrap[I any, O any, R GenericRunnerFunc[I, O]](ctx context.Context, id string, cmd R, cp Pipeline, in I, fs afero.Fs) (res *O, err error) {

	wrk, err := Load(ctx, cp, id, fs)
	if err != nil {
		return nil, err
	}

	zerolog.Ctx(ctx).Debug().Str("id", id).RawJSON("wrk", wrk).Msg("wrap")

	if len(wrk) > 0 {
		err := json.Unmarshal(wrk, &res)
		if err != nil {
			return nil, err
		}

		zerolog.Ctx(ctx).Debug().Str("id", id).Str("wrk", string(wrk)).Msg("wrap loaded from cache")

		return res, nil
	}

	zerolog.Ctx(ctx).Debug().Str("id", id).Msg("wrap not loaded from cache")

	res2, err := cmd(ctx, in)
	if err != nil {
		return nil, err
	}

	res = res2

	inter := (interface{})(res)

	switch z := inter.(type) {
	case nil:
		wrk = []byte{}
	case *string:
		wrk = []byte(*z)
	default:
		wrk, err = json.Marshal(res)
		if err != nil {
			return nil, err
		}
		exp := Express(z)
		if len(exp) > 0 {
			err := AddContentToEnv(ctx, cp, fs, id, exp)
			if err != nil {
				return nil, err
			}
		}
	}

	err = Save(ctx, cp, id, wrk, fs)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// type InnerRunnerFunc[I any] func(context.Context, ...any) (*I, error)

// type Inner[I any] struct {
// 	Func InnerRunnerFunc[I]
// 	Args []any
// }

// func (me *Inner[I]) Run(ctx context.Context) (res *I, err error) {
// 	return me.Func(ctx, me.Cp, me.Args...)
// }

// func Wrap[A any, C RunnerFunc[A]](id string, i C, r Pipeline, a ...any) C {
// 	return func(ctx context.Context) (res *A, err error) {
// 		return wrap[A](ctx, id, i, r, a...)
// 	}
// }
