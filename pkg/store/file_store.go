package store

import (
	"context"
	"errors"
	"fmt"
	"os"
)

type FileStore struct {
	filename string
}

func NewFileStore(filename string) *FileStore {
	return &FileStore{
		filename: filename,
	}
}

func (t *FileStore) Store(ctx context.Context, data []byte) error {
	f, err := os.OpenFile(t.filename, os.O_WRONLY, os.ModePerm)
	if errors.Is(err, os.ErrNotExist) {
		f, err = os.Create(t.filename)
	}
	if err != nil {
		return fmt.Errorf("when os.Open: %w", err)
	}
	_, err = f.Write(data)
	if err != nil {
		return fmt.Errorf("when f.Write: %w", err)
	}
	return nil
}

func (t *FileStore) Get(ctx context.Context) ([]byte, error) {
	return os.ReadFile(t.filename)
}
