package file

import (
	"context"
)

type FileAPI interface {
	Get(ctx context.Context, key string) (res []byte, err error)
	Put(ctx context.Context, key string, data []byte) error
	AppendString(ctx context.Context, key string, data string) error
	Delete(ctx context.Context, key string) error
}

type Client struct {
}
