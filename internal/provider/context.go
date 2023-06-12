package provider

import (
	"context"

	"github.com/alecthomas/kong"
)

type KongContextKey struct{}

type ContentProviderKey struct{}

func BindToKongContext(ctx context.Context, kctx *kong.Context) context.Context {
	ctx = context.WithValue(ctx, KongContextKey{}, kctx)
	kctx.BindTo(ctx, (*context.Context)(nil))
	return ctx
}

func BindContentProvider(ctx context.Context, cp ContentProvider) context.Context {
	return context.WithValue(ctx, ContentProviderKey{}, cp)
}

func GetContentProvider(ctx context.Context) (ContentProvider, bool) {
	if v, ok := ctx.Value(ContentProviderKey{}).(ContentProvider); ok {
		return v, true
	} else {
		return nil, false
	}
}

func GetKongContext(ctx context.Context) *kong.Context {
	return ctx.Value(KongContextKey{}).(*kong.Context)
}
