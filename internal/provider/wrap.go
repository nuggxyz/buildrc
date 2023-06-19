package provider

import (
	"context"
	"encoding/json"
)

type ContentProviderProvider interface {
	Content() ContentProvider
}

type GenericRunnerFunc[I any, O any] func(context.Context, I) (*O, error)

func WrapGeneric[I any, O any](ctx context.Context, id string, cp ContentProvider, in I, i GenericRunnerFunc[I, O]) (*O, error) {
	return wrap(ctx, id, i, cp, in)
}

func Wrap[I ContentProviderProvider, O any](ctx context.Context, id string, in I, i GenericRunnerFunc[I, O]) (*O, error) {
	return WrapGeneric(ctx, id, in.Content(), in, i)
}

func wrap[I any, O any, R GenericRunnerFunc[I, O]](ctx context.Context, id string, cmd R, cp ContentProvider, in I) (res *O, err error) {

	wrk, err := Load(ctx, cp, id)
	if err != nil {
		return nil, err
	}

	if len(wrk) > 0 {
		err := json.Unmarshal(wrk, &res)
		return nil, err
	}

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
			err := AddContentToEnv(ctx, cp, id, exp)
			if err != nil {
				return nil, err
			}
		}
	}

	err = Save(ctx, cp, id, wrk)
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

// func Wrap[A any, C RunnerFunc[A]](id string, i C, r ContentProvider, a ...any) C {
// 	return func(ctx context.Context) (res *A, err error) {
// 		return wrap[A](ctx, id, i, r, a...)
// 	}
// }
