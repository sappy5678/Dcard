package repository

import (
	"context"

	"github.com/sappy5678/dcard/pkg/domain"
)

type MockShortURLRepository struct {
	CreateFunc func(ctx context.Context, short *domain.ShortURL) (*domain.ShortURL, error)
	GetFunc    func(ctx context.Context, shortCode string) (*domain.ShortURL, error)
}

func (m *MockShortURLRepository) Create(ctx context.Context, short *domain.ShortURL) (*domain.ShortURL, error) {
	return m.CreateFunc(ctx, short)
}

func (m *MockShortURLRepository) Get(ctx context.Context, shortCode string) (*domain.ShortURL, error) {
	return m.GetFunc(ctx, shortCode)
}
