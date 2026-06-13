package storage

import "context"

type NoopStorage struct{}

func NewNoop() Storage {
	return &NoopStorage{}
}

func (n *NoopStorage) Put(context.Context, string, []byte) error {
	return nil
}

func (n *NoopStorage) Delete(context.Context, string) error {
	return nil
}
