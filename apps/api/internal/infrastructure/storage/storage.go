package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Provider is the storage abstraction used by use cases.
type Provider interface {
	Upload(ctx context.Context, key string, r io.Reader, size int64, mimeType string) (string, error)
	Delete(ctx context.Context, key string) error
	URL(key string) string
}

// LocalProvider stores files on the local filesystem — for development/MVP.
type LocalProvider struct {
	basePath string
	baseURL  string
}

func NewLocalProvider(basePath, baseURL string) (*LocalProvider, error) {
	if err := os.MkdirAll(basePath, 0o755); err != nil {
		return nil, fmt.Errorf("storage: mkdir: %w", err)
	}
	return &LocalProvider{basePath: basePath, baseURL: baseURL}, nil
}

func (p *LocalProvider) Upload(_ context.Context, key string, r io.Reader, _ int64, _ string) (string, error) {
	dest := filepath.Join(p.basePath, key)
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return "", fmt.Errorf("storage: mkdir: %w", err)
	}
	f, err := os.Create(dest)
	if err != nil {
		return "", fmt.Errorf("storage: create file: %w", err)
	}
	defer f.Close()
	if _, err := io.Copy(f, r); err != nil {
		return "", fmt.Errorf("storage: write: %w", err)
	}
	return key, nil
}

func (p *LocalProvider) Delete(_ context.Context, key string) error {
	return os.Remove(filepath.Join(p.basePath, key))
}

func (p *LocalProvider) URL(key string) string {
	return p.baseURL + "/" + key
}
