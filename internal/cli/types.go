package cli

import (
	"context"
	"encoding/json"
)

type WrappedHandlerFunc[I any] func(context.Context) (*I, error)

type Command[I any] interface {
	Invoke(context.Context, ContentProvider) (*I, error)
	// Wrapped(cp ContentProvider) WrappedHandlerFunc[I]
	Init(ctx context.Context) error
	Helper() CommandHelper[I]
	ID() string
}

type CommandRunner interface {
	AnyHelper() AnyHelper
	ID() string
}

func LoadFunc[I any](cmd Command[I], cp ContentProvider) WrappedHandlerFunc[I] {

	return func(ctx context.Context) (res *I, err error) {

		err = cmd.Init(ctx)
		if err != nil {
			return nil, err
		}

		wrk, err := cp.Load(ctx, cmd)
		if err != nil {
			return nil, err
		}

		res = new(I)

		if len(wrk) > 0 {
			err := json.Unmarshal(wrk, res)
			return res, err
		}

		res, err = cmd.Invoke(ctx, cp)
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
		}

		err = cp.Save(ctx, cmd, wrk)
		if err != nil {
			return nil, err
		}

		return res, nil
	}
}

var _ CommandHelper[string] = (*Helper[string])(nil)

var _ AnyHelper = (*Helper[string])(nil)

type Helper[I any] struct {
	Command Command[I]
}

func NewHelper[I any](cmd Command[I]) *Helper[I] {
	return &Helper[I]{Command: cmd}
}

type AnyHelper interface {
	Start(ctx context.Context, cp ContentProvider) (any, error)
}

func (me *Helper[I]) Start(ctx context.Context, cp ContentProvider) (any, error) {
	return me.Wrapped(cp)(ctx)
}

type CommandHelper[I any] interface {
	Wrapped(cp ContentProvider) WrappedHandlerFunc[I]
	Run(ctx context.Context, cp ContentProvider) (*I, error)
}

func (me *Helper[I]) Wrapped(cp ContentProvider) WrappedHandlerFunc[I] {
	return LoadFunc(me.Command, cp)
}

func (me *Helper[I]) Run(ctx context.Context, cp ContentProvider) (*I, error) {
	return me.Wrapped(cp)(ctx)
}
