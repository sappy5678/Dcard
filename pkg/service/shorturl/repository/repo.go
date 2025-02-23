package repository

import (
	"context"

	"github.com/sappy5678/dcard/pkg/domain"
)

// Repository is a repository for shorturl
type Repository interface {
	Create(ctx context.Context, short *domain.ShortURL) (*domain.ShortURL, error)
	Get(ctx context.Context, shortCode string) (*domain.ShortURL, error)
}
