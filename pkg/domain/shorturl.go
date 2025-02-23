package domain

import (
	"context"
	"fmt"
	"net/url"
)

var (
	ErrShortURLNotFound = fmt.Errorf("short url not found")
	ErrShortURLInvalid  = fmt.Errorf("short url invalid")
)

type ShortURL struct {
	ShortCode   string `json:"shortCode" db:"short_code"`
	OriginalURL string `json:"originalUrl" db:"original_url"`
	ShortURL    string `json:"shortUrl" db:"-"`
	ExpireTime  uint64 `json:"expireTime" db:"expire_time"`
	CreatedTime uint64 `json:"createdTime" db:"created_time"`
}

func (s *ShortURL) IsValid(nowUnix uint64) bool {
	if s.ExpireTime == 0 || nowUnix > s.ExpireTime || s.ExpireTime < s.CreatedTime {
		return false
	}
	if s.CreatedTime == 0 {
		return false
	}
	if s.ShortCode == "" {
		return false
	}
	if s.OriginalURL == "" {
		return false
	}
	_, err := url.ParseRequestURI(s.OriginalURL)
	if err != nil {
		return false
	}
	return true
}

type ShortURLService interface {
	Create(ctx context.Context, originalURL string, expireTime uint64) (*ShortURL, error)
	Get(ctx context.Context, shortCode string) (*ShortURL, error)
}
