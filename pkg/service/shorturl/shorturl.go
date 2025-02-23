package shorturl

import (
	"context"
	"fmt"

	"github.com/sappy5678/dcard/pkg/domain"
)

func (im *shorturlService) getShortURL(shortCode string) string {
	return fmt.Sprintf("%s/%s", im.host, shortCode)
}

func (im *shorturlService) Create(ctx context.Context, originalURL string, expireTime uint64) (*domain.ShortURL, error) {
	shortCode := im.shortcodeGenerator.NextID()
	shortURL := &domain.ShortURL{
		ShortCode:   shortCode,
		OriginalURL: originalURL,
		ShortURL:    im.getShortURL(shortCode),
		ExpireTime:  expireTime,
		CreatedTime: im.now(),
	}
	if !shortURL.IsValid(im.now()) {
		return nil, domain.ErrShortURLInvalid
	}
	shortURL, err := im.repo.Create(ctx, shortURL)
	if err != nil {
		return nil, err
	}
	return shortURL, nil
}

func (im *shorturlService) Get(ctx context.Context, shortCode string) (*domain.ShortURL, error) {
	shortURL, err := im.repo.Get(ctx, shortCode)
	if err != nil {
		return nil, err
	}
	shortURL.ShortURL = im.getShortURL(shortCode)
	if shortURL == nil || !shortURL.IsValid(im.now()) {
		return nil, domain.ErrShortURLNotFound
	}
	return shortURL, nil
}
