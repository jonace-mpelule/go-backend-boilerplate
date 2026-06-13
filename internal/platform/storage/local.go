package storage

import (
	"context"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	root string
}

func NewLocal(root string) Storage {
	return &LocalStorage{root: root}
}

func (l *LocalStorage) Put(ctx context.Context, path string, contents []byte) error {
	target := filepath.Join(l.root, path)
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}
	return os.WriteFile(target, contents, 0o644)
}

func (l *LocalStorage) Delete(ctx context.Context, path string) error {
	return os.Remove(filepath.Join(l.root, path))
}
