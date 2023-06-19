package provider

import "context"

type RunnerFunc2[I any] func(context.Context, ContentProvider, ...any) (*I, error)

func Wrap2[A any](id string, i RunnerFunc2[A]) RunnerFunc2[A] {
	return func(ctx context.Context, cp ContentProvider, r ...any) (res *A, err error) {
		return wrap2[A](ctx, id, i, cp, r)
	}
}

func wrap2[I any, C RunnerFunc2[I]](ctx context.Context, id string, cmd C, cp ContentProvider, a ...any) (res *I, err error) {
	return nil, nil
}
