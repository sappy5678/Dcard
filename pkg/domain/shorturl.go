package domain

import (
	"context"
	"fmt"
)

var ErrShortURLNotFound = fmt.Errorf("short url not found")

type ShortURL struct {
	ShortCode   string `json:"shortCode" db:"short_code"`
	OriginalURL string `json:"originalUrl" db:"original_url"`
	ShortURL    string `json:"shortUrl" db:"-"`
	ExpireTime  int    `json:"expireTime" db:"expire_time"`
	CreatedTime int    `json:"createdTime" db:"created_time"`
}

type ShortURLService interface {
	Create(ctx context.Context, originalURL string, expireTime int) (*ShortURL, error)
	Get(ctx context.Context, shortCode string) (*ShortURL, error)
}
