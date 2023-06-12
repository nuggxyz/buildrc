package provider

import (
	"context"

	"github.com/rs/zerolog"
)

type Identifiable interface {
	ID() string
}

type Expressable interface {
	Express() map[string]string
}

type ContentProvider interface {
	Load(context.Context, Identifiable) ([]byte, error)
	Express(context.Context, Identifiable, Expressable) error
	Save(context.Context, Identifiable, []byte) error
}

var _ ContentProvider = (*NoopContentProvider)(nil)

type NoopContentProvider struct {
	LoadFunc func(context.Context, Identifiable) ([]byte, error)
	SaveFunc func(context.Context, Identifiable, []byte) error

	SaveCalled bool
	LoadCalled bool

	SaveCalledWithId Identifiable
	LoadCalledWithId Identifiable

	SaveCalledWithBytes []byte

	LoadBytes []byte
	LoadError error

	SaveError error
}

func (me *NoopContentProvider) HasRun(ider Identifiable) bool {
	return me.LoadCalledWithId == ider
}

func (me *NoopContentProvider) Load(ctx context.Context, ider Identifiable) ([]byte, error) {
	return me.LoadFunc(ctx, ider)
}

func (me *NoopContentProvider) Save(ctx context.Context, ider Identifiable, b []byte) error {
	zerolog.Ctx(ctx).Debug().Str("id", ider.ID()).Str("json", string(b)).Msg("save")
	return me.SaveFunc(ctx, ider, b)
}

func NewNoopContentProvider(output []byte) *NoopContentProvider {
	mcp := &NoopContentProvider{}
	mcp.LoadBytes = output
	mcp.LoadFunc = func(_ context.Context, ider Identifiable) ([]byte, error) {
		mcp.LoadCalled = true
		mcp.LoadCalledWithId = ider
		return mcp.LoadBytes, mcp.LoadError
	}
	mcp.SaveFunc = func(_ context.Context, ider Identifiable, b []byte) error {
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

func (me *NoopContentProvider) Express(context.Context, Identifiable, Expressable) error {
	return nil
}
