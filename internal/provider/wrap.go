package provider

import (
	"context"
	"encoding/json"
)

type RunnerFunc[I any] func(context.Context, ContentProvider) (*I, error)

func Wrap[A any](id string, i RunnerFunc[A]) RunnerFunc[A] {
	return func(ctx context.Context, r ContentProvider) (res *A, err error) {
		return wrap[A](ctx, id, i, r)
	}
}

func wrap[I any, C RunnerFunc[I]](ctx context.Context, id string, cmd C, cp ContentProvider) (res *I, err error) {

	wrk, err := cp.Load(ctx, id)
	if err != nil {
		return nil, err
	}

	res = new(I)

	if len(wrk) > 0 {
		err := json.Unmarshal(wrk, res)
		return res, err
	}

	res, err = cmd(ctx, cp)
	if err != nil {
		return nil, err
	}

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
			err := cp.Express(ctx, id, exp)
			if err != nil {
				return nil, err
			}
		}
	}

	err = cp.Save(ctx, id, wrk)
	if err != nil {
		return nil, err
	}

	return res, nil
}
