package pipeline

import (
	"context"

	"github.com/alecthomas/kong"
)

type KongContextKey struct{}

// type PipelineKey struct{}

func BindToKongContext(ctx context.Context, kctx *kong.Context) context.Context {
	ctx = context.WithValue(ctx, KongContextKey{}, kctx)
	kctx.BindTo(ctx, (*context.Context)(nil))
	return ctx
}

func GetKongContext(ctx context.Context) *kong.Context {
	return ctx.Value(KongContextKey{}).(*kong.Context)
}
