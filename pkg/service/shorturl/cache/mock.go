package cache

import (
	"context"

	"github.com/sappy5678/dcard/pkg/domain"
)

type MockShortURLCacheRepository struct {
	CreateFunc func(ctx context.Context, short *domain.ShortURL) (*domain.ShortURL, error)
	GetFunc    func(ctx context.Context, shortCode string) (*domain.ShortURL, error)
}

func (m *MockShortURLCacheRepository) Create(ctx context.Context, short *domain.ShortURL) (*domain.ShortURL, error) {
	return m.CreateFunc(ctx, short)
}

func (m *MockShortURLCacheRepository) Get(ctx context.Context, shortCode string) (*domain.ShortURL, error) {
	return m.GetFunc(ctx, shortCode)
}
