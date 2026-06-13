package storage

import "context"

type Storage interface {
	Put(ctx context.Context, path string, contents []byte) error
	Delete(ctx context.Context, path string) error
}
