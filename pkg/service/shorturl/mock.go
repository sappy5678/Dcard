package shorturl

import (
	"context"

	"github.com/sappy5678/dcard/pkg/domain"
)

type MockShortURLService struct {
	CreateFunc func(ctx context.Context, originalURL string, expireTime uint64) (*domain.ShortURL, error)
	GetFunc    func(ctx context.Context, shortCode string) (*domain.ShortURL, error)
}

func (m *MockShortURLService) Create(ctx context.Context, originalURL string, expireTime uint64) (*domain.ShortURL, error) {
	return m.CreateFunc(ctx, originalURL, expireTime)
}

func (m *MockShortURLService) Get(ctx context.Context, shortCode string) (*domain.ShortURL, error) {
	return m.GetFunc(ctx, shortCode)
}
