package logservice

import (
	"context"
	"time"

	"github.com/sappy5678/dcard/pkg/domain"
)

func New(svc domain.ShortURLService, logger domain.Logger) *LogService {
	return &LogService{
		ShortURLService: svc,
		logger:          logger,
	}
}

// LogService represents wallet logging service
type LogService struct {
	domain.ShortURLService
	logger domain.Logger
}

const name = "shorturl"

// Create(ctx context.Context, originalURL string, expireTime uint64) (*ShortURL, error)
//
//	Get(ctx context.Context, shortCode string) (*ShortURL, error)
func (ls *LogService) Create(ctx context.Context, originalURL string, expireTime uint64) (short *domain.ShortURL, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			ctx,
			name, "Create shorturl request", err,
			map[string]interface{}{
				"originalURL": originalURL,
				"expireTime":  expireTime,
				"took":        time.Since(begin),
			},
		)
	}(time.Now())

	return ls.ShortURLService.Create(ctx, originalURL, expireTime)
}

func (ls *LogService) Get(ctx context.Context, shortCode string) (short *domain.ShortURL, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			ctx,
			name, "Create shorturl request", err,
			map[string]interface{}{
				"shortCode": shortCode,
				"took":      time.Since(begin),
			},
		)
	}(time.Now())

	return ls.ShortURLService.Get(ctx, shortCode)
}
