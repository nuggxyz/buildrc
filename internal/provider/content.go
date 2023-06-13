package provider

import (
	"context"

	"github.com/rs/zerolog"
)

type Identifiable interface {
	ID() string
}

type ContentProvider interface {
	Load(context.Context, string) ([]byte, error)
	Express(context.Context, string, map[string]string) error
	Save(context.Context, string, []byte) error
}

var _ ContentProvider = (*NoopContentProvider)(nil)

type NoopContentProvider struct {
	LoadFunc func(context.Context, string) ([]byte, error)
	SaveFunc func(context.Context, string, []byte) error

	SaveCalled bool
	LoadCalled bool

	SaveCalledWithId string
	LoadCalledWithId string

	SaveCalledWithBytes []byte

	LoadBytes []byte
	LoadError error

	SaveError error
}

func (me *NoopContentProvider) HasRun(ider string) bool {
	return me.LoadCalledWithId == ider
}

func (me *NoopContentProvider) Load(ctx context.Context, ider string) ([]byte, error) {
	return me.LoadFunc(ctx, ider)
}

func (me *NoopContentProvider) Save(ctx context.Context, ider string, b []byte) error {
	zerolog.Ctx(ctx).Debug().Str("id", ider).Str("json", string(b)).Msg("save")
	return me.SaveFunc(ctx, ider, b)
}

func NewNoopContentProvider(output []byte) *NoopContentProvider {
	mcp := &NoopContentProvider{}
	mcp.LoadBytes = output
	mcp.LoadFunc = func(_ context.Context, ider string) ([]byte, error) {
		mcp.LoadCalled = true
		mcp.LoadCalledWithId = ider
		return mcp.LoadBytes, mcp.LoadError
	}
	mcp.SaveFunc = func(_ context.Context, ider string, b []byte) error {
		mcp.SaveCalled = true
		mcp.SaveCalledWithId = ider
		mcp.SaveCalledWithBytes = b
		return mcp.SaveError
	}
	return mcp
}

func (me *NoopContentProvider) LoadBytesReturn(b []byte) {
	me.LoadBytes = b
}

func (me *NoopContentProvider) Express(context.Context, string, map[string]string) error {
	return nil
}
