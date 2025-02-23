package logservice_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/sappy5678/dcard/pkg/domain"
	"github.com/sappy5678/dcard/pkg/service/shorturl"
	sl "github.com/sappy5678/dcard/pkg/service/shorturl/logservice"
	"github.com/sappy5678/dcard/pkg/utl/zlog"
)

var (
	mockCreatedTimeString = "2025-01-01T00:00:00Z"
	mockExpireTimeString  = "2025-01-01T01:00:00Z"
	mockCreatedTime, _    = time.Parse(time.RFC3339, mockCreatedTimeString)
	mockExpireTime        = mockCreatedTime.Add(time.Hour)
	mockOriginalURL       = "test-original-url"
	mockShortCode         = "test-short-code"
	mockShortURL          = "test-short-url"
	mockShort             = &domain.ShortURL{
		ShortCode:   mockShortCode,
		OriginalURL: mockOriginalURL,
		ShortURL:    mockShortURL,
		ExpireTime:  uint64(mockExpireTime.Unix()),
		CreatedTime: uint64(mockCreatedTime.Unix()),
	}
	mockError = errors.New("error")
)

var mockShortURLService = &shorturl.MockShortURLService{
	CreateFunc: func(ctx context.Context, originalURL string, expireTime uint64) (*domain.ShortURL, error) {
		if originalURL != mockShort.OriginalURL || expireTime != mockShort.ExpireTime {
			return nil, mockError
		}
		return mockShort, nil
	},
	GetFunc: func(ctx context.Context, shortCode string) (*domain.ShortURL, error) {
		return mockShort, nil
	},
}

func TestCreate(t *testing.T) {
	log := zlog.New()
	svc := sl.New(mockShortURLService, log)
	r1, e1 := svc.Create(context.Background(), mockOriginalURL, uint64(mockExpireTime.Unix()))
	r2, e2 := mockShortURLService.Create(context.Background(), mockOriginalURL, uint64(mockExpireTime.Unix()))

	assert.Equal(t, r1, r2)
	assert.Equal(t, e1, e2)
}

func TestGet(t *testing.T) {
	log := zlog.New()
	svc := sl.New(mockShortURLService, log)
	r1, e1 := svc.Get(context.Background(), mockShortCode)
	r2, e2 := mockShortURLService.Get(context.Background(), mockShortCode)

	assert.Equal(t, r1, r2)
	assert.Equal(t, e1, e2)
}
