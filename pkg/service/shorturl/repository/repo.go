package repository

import (
	"context"

	"github.com/sappy5678/dcard/pkg/domain"
)

// ShortURLRepository is a repository for shorturl
type ShortURLRepository interface {
	Create(ctx context.Context, short *domain.ShortURL) (*domain.ShortURL, error)
	Get(ctx context.Context, shortCode string) (*domain.ShortURL, error)
}
